package models

type Address struct {
	Country  Country  `json:"country" db:",prefix=co."`
	City     City     `json:"city" db:",prefix=c."`
	Town     Town     `json:"town" db:",prefix=t."`
	District District `json:"district" db:",prefix=d."`
	Quarter  Quarter  `json:"quarter" db:",prefix=q"`
}

type Countries []Country
type Cities []City
type Towns []Town
type Districts []District
type Quarters []Quarter

type Country struct {
	ID           int    `json:"id" gorm:"column:id;primary_key;type:integer;not null" default:"1"`
	Name         string `json:"name" gorm:"column:name;type:varchar(64);not null" default:"Türkiye"`
	Abbreviation string `json:"abbreviation" gorm:"column:abbreviation;type:varchar(2);not null" default:"TR"`
	Language     string `json:"language" gorm:"column:language;type:varchar(2);not null" default:"tr"`
	DisplayOrder int    `json:"display_order" gorm:"column:display_order;type:integer;not null;default:0" default:"1"`
	SortOrder    int    `json:"sort_order" gorm:"column:sort_order;type:integer;not null;default:0" default:"1"`
	PhoneCode    string `json:"phone_code" gorm:"column:phone_code;type:varchar(5);not null" default:"+90"`
	Alpha2Code   string `json:"alpha2_code" gorm:"column:alpha2_code;type:varchar(2);not null" default:"TR"`
	Alpha3Code   string `json:"alpha3_code" gorm:"column:alpha3_code;type:varchar(3);not null" default:"TUR"`
}

type City struct {
	ID           int     `json:"id" gorm:"column:id;primary_key;type:integer;not null" example:"43"`
	CountryID    int     `json:"country_id" gorm:"column:country_id;type:integer;not null" example:"1"`
	Name         string  `json:"name" gorm:"column:name;type:varchar(64);not null" example:"Kütahya"`
	Tag          string  `json:"tag" gorm:"column:tag;type:varchar(64);not null" example:"kutahya"`
	DisplayOrder int     `json:"display_order" gorm:"column:display_order;type:integer;not null;default:0" example:"43"`
	SortOrder    int     `json:"sort_order" gorm:"column:sort_order;type:integer;not null;default:0" example:"43"`
	Country      Country `json:"country" gorm:"foreignKey:CountryID;references:ID"`
}

type Town struct {
	ID           int    `json:"id" gorm:"column:id;primary_key;type:integer;not null" example:"587"`
	CityID       int    `json:"city_id" gorm:"column:city_id;type:integer;not null" example:"43"`
	Name         string `json:"name" gorm:"column:name;type:varchar(64);not null" example:"Merkez"`
	DisplayOrder int    `json:"display_order" gorm:"column:display_order;type:integer;not null;default:0" default:"0"`
	SortOrder    int    `json:"sort_order" gorm:"column:sort_order;type:integer;not null;default:0" default:"0"`
	City         City   `json:"city" gorm:"foreignKey:CityID;references:ID"`
}

type District struct {
	ID           int    `json:"id" gorm:"column:id;primary_key;type:integer;not null" default:"2800"`
	TownID       int    `json:"town_id" gorm:"column:town_id;type:integer;not null" default:"587"`
	Name         string `json:"name" gorm:"column:name;type:varchar(64);not null" default:"Alipaşa"`
	DisplayOrder int    `json:"display_order" gorm:"column:display_order;type:integer;not null;default:0" default:"0"`
	SortOrder    int    `json:"sort_order" gorm:"column:sort_order;type:integer;not null;default:0" default:"2"`
	Town         Town   `json:"town" gorm:"foreignKey:TownID;references:ID"`
}

type Quarter struct {
	ID           int      `json:"id" gorm:"column:id;primary_key;type:integer;not null" default:"29981"`
	DistrictID   int      `json:"district_id" gorm:"column:district_id;type:integer;not null" default:"2800"`
	Name         string   `json:"name" gorm:"column:name;type:varchar(64);not null" default:"Alipaşa Mh."`
	DisplayOrder int      `json:"display_order" gorm:"column:display_order;type:integer;not null;default:0" default:"0"`
	SortOrder    int      `json:"sort_order" gorm:"column:sort_order;type:integer;not null;default:0" default:"0"`
	Detail       string   `json:"detail" gorm:"column:detail;type:varchar(255)" default:"{\"location_id\": 70910, \"lat\": 39.418975, \"lon\": 29.983876, \"zoom\": 9, \"quarterId\": 29981}"`
	District     District `json:"district" gorm:"foreignKey:DistrictID;references:ID"`
}
