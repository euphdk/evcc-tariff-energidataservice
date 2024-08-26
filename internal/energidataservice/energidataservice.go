package energidataservice

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/tidwall/gjson"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

const (
	ElspotpricesURI        = "https://api.energidataservice.dk/dataset/Elspotprices?offset=0&start=%s&end=%s&filter={\"PriceArea\":[\"%s\"]}&timezone=dk&limit=48"
	DatahubPricelistURI    = "https://api.energidataservice.dk/dataset/DatahubPricelist?offset=0&filter=%s&sort=ValidFrom%%20desc&limit=10"
	DatahubPricelistFilter = "{\"ChargeTypeCode\":%s,\"GLN_Number\":[\"%s\"],\"ChargeType\":%s}"
	TimeFormat             = "2006-01-02T15:04" // RFC3339 short
	TimeFormatSecond       = "2006-01-02T15:04:05"
)

type EvccAPIRate struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Price float64   `json:"price"`
}

type Gridcharge struct {
	Start  time.Time
	End    time.Time
	Prices [24]float64 // Hour of Day
}

func GetEvccAPIRates(gridCompany, region string, tax, vat float64) ([]EvccAPIRate, error) {

	datahubPricelist, err := getDatahubPricelist(*ChargeOwners[gridCompany])
	if err != nil {
		return []EvccAPIRate{}, err
	}

	elspotprices, err := getElspotprices(region)
	if err != nil {
		return []EvccAPIRate{}, err
	}

	data := make([]EvccAPIRate, 0)

	for unixTimestamp, price := range elspotprices {
		date := time.Unix(unixTimestamp, 0)
		var gridcharge float64

		// Loop through datahubPricelist and sum all charges (or rebates) if time matches
		for _, entry := range datahubPricelist {
			if date.After(entry.Start) && date.Before(entry.End) {
				gridcharge = gridcharge + entry.Prices[date.Hour()]
			}
		}

		r := EvccAPIRate{
			Start: date.Local(),
			End:   date.Add(time.Hour).Local(),
			Price: (price/1e3 + gridcharge + tax) * vat,
		}
		data = append(data, r)
	}

	return data, nil
}

func getElspotprices(region string) (map[int64]float64, error) {

	ts := time.Now().Truncate(time.Hour)
	uri := fmt.Sprintf(ElspotpricesURI,
		ts.Format(TimeFormat),
		ts.Add(48*time.Hour).Format(TimeFormat),
		region)

	slog.Debug("Elspotprices", "uri", uri)

	r, err := httpClient.Get(uri)
	if err != nil {
		slog.Error("Failed GET Elspotprices", "error", err.Error())
		return map[int64]float64{}, err
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed reading body", "error", err.Error())
		return map[int64]float64{}, err
	}

	records := gjson.GetBytes(body, "records")

	prices := make(map[int64]float64, 0)

	for _, record := range records.Array() {
		date, _ := time.Parse(TimeFormatSecond, record.Get("HourUTC").Str)
		price := record.Get("SpotPriceDKK").Float()
		prices[date.Unix()] = price
	}
	return prices, nil
}

func getDatahubPricelist(chargeOwner ChargeOwner) ([]Gridcharge, error) {

	// Build URI
	jsonChargeTypeCode, _ := json.Marshal(chargeOwner.ChargeTypeCode)
	jsonChargeType, _ := json.Marshal(chargeOwner.ChargeType)
	filter := fmt.Sprintf(DatahubPricelistFilter, jsonChargeTypeCode, chargeOwner.GLN, jsonChargeType)
	uri := fmt.Sprintf(DatahubPricelistURI, url.QueryEscape(filter))
	slog.Debug("Constructed URI for DatahubPricelist: " + uri)

	r, err := httpClient.Get(uri)
	if err != nil {
		slog.Error("Failed GET DatahubPricelist", "error", err.Error())
		return []Gridcharge{}, err
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed reading body", "error", err.Error())
		return []Gridcharge{}, err
	}

	records := gjson.GetBytes(body, "records")

	var gridcharges []Gridcharge

	for _, record := range records.Array() {
		gridcharges = append(gridcharges, jsonresultToGridcharge(record))
	}

	return gridcharges, nil
}

func jsonresultToGridcharge(record gjson.Result) Gridcharge {
	validFrom, err := time.Parse(TimeFormatSecond, record.Get("ValidFrom").Str)
	if err != nil {
		slog.Error("Invalid date", "error", err.Error())
		return Gridcharge{}
	}

	validTo := time.Now().Add(72 * time.Hour)
	if record.Get("ValidTo").Str != "" {
		validTo, err = time.Parse(TimeFormatSecond, record.Get("ValidTo").Str)
		if err != nil {
			slog.Error("Invalid date", "error", err.Error())
			return Gridcharge{}
		}
	}

	var prices [24]float64
	basePrice := record.Get("Price1").Float()

	for i := 0; i < 24; i++ {
		currentPrice := fmt.Sprintf("Price%d", i+1)
		if record.Get(currentPrice).Raw != "" {
			price := record.Get(currentPrice).Float()
			prices[i] = price
		} else {
			prices[i] = basePrice
		}
	}

	gridcharge := Gridcharge{
		Start:  validFrom,
		End:    validTo,
		Prices: prices,
	}

	return gridcharge
}
