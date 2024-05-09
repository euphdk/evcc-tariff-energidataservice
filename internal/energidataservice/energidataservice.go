package energidataservice

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps"
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

func GetEvccAPIRates(gridCompany, region string) []EvccAPIRate {

	// FIXME: returned unix-timestamps doesn't align
	datahubPricelist := getDatahubPricelist(*ChargeOwners[gridCompany])
	elspotprices := getElspotprices(region)

	data := make([]EvccAPIRate, 0)

	for unixTimestamp, price := range elspotprices {
		date := time.Unix(unixTimestamp, 0)
		r := EvccAPIRate{
			Start: date.Local(),
			End: date.Add(time.Hour).Local(),
			Price: price / 1e3 + datahubPricelist[unixTimestamp],
		}
		data = append(data, r)
	}

	return data
}

func getElspotprices(region string) map[int64]float64 {

	ts := time.Now().Truncate(time.Hour)
	uri := fmt.Sprintf(ElspotpricesURI,
		ts.Format(TimeFormat),
		ts.Add(24*time.Hour).Format(TimeFormat),
		region)

	r, err := httpClient.Get(uri)
	if err != nil {
		slog.Error("Failed GET Elspotprices", "error", err.Error())
		return map[int64]float64{}
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed reading body", "error", err.Error())
		return map[int64]float64{}
	}

	records := gjson.GetBytes(body, "records")

	prices := make(map[int64]float64,0)

	for _, record := range records.Array() {
		date, _ := time.Parse(TimeFormatSecond, record.Get("HourUTC").Str)
		price := record.Get("SpotPriceDKK").Float()
		prices[date.Unix()] = price
	}
	return prices
}

func getDatahubPricelist(chargeOwner ChargeOwner) map[int64]float64 {

	// Build URI
	jsonChargeTypeCode, _ := json.Marshal(chargeOwner.ChargeTypeCode)
	jsonChargeType, _ := json.Marshal(chargeOwner.ChargeType)
	filter := fmt.Sprintf(DatahubPricelistFilter, jsonChargeTypeCode, chargeOwner.GLN, jsonChargeType)
	uri := fmt.Sprintf(DatahubPricelistURI, url.QueryEscape(filter))
	slog.Info("Constructed URI for DatahubPricelist: " + uri)

	r, err := httpClient.Get(uri)
	if err != nil {
		slog.Error("Failed GET DatahubPricelist", "error", err.Error())
		return map[int64]float64{}
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed reading body", "error", err.Error())
		return map[int64]float64{}
	}

	records := gjson.GetBytes(body, "records")

	gridCharge := parseDatahubPricelistRecord(records, time.Now())
	gridChargeTomorrow := parseDatahubPricelistRecord(records, time.Now().Add(24*time.Hour))
	maps.Copy(gridCharge, gridChargeTomorrow)

	// slog.Info(fmt.Sprintf("%#v", gridCharge))

	return gridCharge
}

func parseDatahubPricelistRecord(records gjson.Result, date time.Time) map[int64]float64 {
	gridCharge := make(map[int64]float64)

	for _, record := range records.Array() {
		validFrom, err := time.Parse(TimeFormatSecond, record.Get("ValidFrom").Str)
		if err != nil {
			slog.Error("Invalid date", "error", err.Error())
		}

		// validTo might be blank - prepare for that and just override if not...
		validTo := time.Now().Add(24 * time.Hour)
		if record.Get("ValidTo").Str != "" {
			validTo, err = time.Parse(TimeFormatSecond, record.Get("ValidTo").Str)
			if err != nil {
				slog.Error("Invalid date", "error", err.Error())
				return map[int64]float64{}
			}
		}

		if dateInRange(validFrom, validTo, date) {
			baseTime := time.Date(
				date.Year(),
				date.Month(),
				date.Day(),
				0, 0, 0, 0, date.Location(),
			)

			basePrice := record.Get("Price1").Float()

			for i := 1; i <= 24; i++ {
				currentPrice := fmt.Sprintf("Price%d", i)
				price := basePrice
				if record.Get(currentPrice).Raw != "" {
					price = record.Get(currentPrice).Float()
				}
				currentHour := baseTime.Add(time.Duration(1-i) * time.Hour).Unix()
				gridCharge[currentHour] = price + gridCharge[currentHour]
			}
		}
	}

	return gridCharge

}

func dateInRange(from, to, date time.Time) bool {
	if from.Before(date) && to.After(date) {
		return true
	}

	return false
}
