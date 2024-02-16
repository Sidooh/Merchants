package entities

type Location struct {
	ModelID

	CountyId    int    `json:"county_id" gorm:"size:8"`
	County      string `json:"county" gorm:"varchar; size:32"`
	SubCountyId int    `json:"sub_county_id" gorm:"size:32"`
	SubCounty   string `json:"sub_county" gorm:"varchar; size:32"`
	WardId      int    `json:"ward_id" gorm:"size:32"`
	Ward        string `json:"ward" gorm:"varchar; size:32"`
	LandmarkId  int    `json:"landmark_id" gorm:"uniqueIndex; size:64"`
	Landmark    string `json:"landmark" gorm:"varchar; size:64"`
}
