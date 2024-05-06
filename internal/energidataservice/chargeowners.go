package energidataservice

type ChargeOwner struct {
	GLN            string
	Company        string
	ChargeTypeCode []string
	ChargeType     []string
}

// Values collected from:
// https://github.com/MTrab/energidataservice/blob/master/custom_components/energidataservice/tariffs/energidataservice/chargeowners.py

var ChargeOwners = map[string]*ChargeOwner{
	"N1": {GLN: "5790001089030", Company: "N1 A/S - 131", ChargeTypeCode: []string{"CD", "CD R"}, ChargeType: []string{"D03"}},
}
