package energidataservice

const (
	ElspotpricesURI        = "https://api.energidataservice.dk/dataset/Elspotprices?offset=0&start=%s&end=%s&filter={\"PriceArea\":[\"%s\"]}&timezone=dk&limit=48"
	DatahubPricelistURI    = "https://api.energidataservice.dk/dataset/DatahubPricelist?offset=0&filter=%s&sort=ValidFrom%%20desc&limit=10"
	DatahubPricelistFilter = "{\"ChargeTypeCode\":%s,\"GLN_Number\":[\"%s\"],\"ChargeType\":%s}"
	TimeFormat             = "2006-01-02T15:04" // RFC3339 short
	TimeFormatSecond       = "2006-01-02T15:04:05"
)

type ChargesRecords struct {
	Records []ChargeRecord `json:"records"`
}

type ChargeRecord struct {
	ChargeOwner          string
	GLN_Number           string
	ChargeType           string
	ChargeTypeCode       string
	Note                 string
	Description          string
	ValidFrom            string
	ValidTo              string
	VATClass             string
	Price1               float64
	Price2               float64
	Price3               float64
	Price4               float64
	Price5               float64
	Price6               float64
	Price7               float64
	Price8               float64
	Price9               float64
	Price10              float64
	Price11              float64
	Price12              float64
	Price13              float64
	Price14              float64
	Price15              float64
	Price16              float64
	Price17              float64
	Price18              float64
	Price19              float64
	Price20              float64
	Price21              float64
	Price22              float64
	Price23              float64
	Price24              float64
	TransparentInvoicing int
	TaxIndicator         int
	ResolutionDuration   string
}

