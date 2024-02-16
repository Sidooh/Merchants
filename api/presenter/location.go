package presenter

type County struct {
	CountyId uint   `json:"id"`
	County   string `json:"county"`
}

type SubCounty struct {
	SubCountyId uint   `json:"id"`
	SubCounty   string `json:"sub_county"`
}

type Ward struct {
	WardId uint   `json:"id"`
	Ward   string `json:"ward"`
}

type Landmark struct {
	LandmarkId uint   `json:"id"`
	Landmark   string `json:"landmark"`
}
