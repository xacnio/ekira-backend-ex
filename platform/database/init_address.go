package database

import (
	"ekira-backend/app/models"
	"os"
)

var CitiesData = models.Cities{
	{ID: 1, CountryID: 1, Name: "Adana", Tag: "adana", DisplayOrder: 10, SortOrder: 1},
	{ID: 2, CountryID: 1, Name: "Adıyaman", Tag: "adiyaman", DisplayOrder: 10, SortOrder: 2},
	{ID: 3, CountryID: 1, Name: "Afyonkarahisar", Tag: "afyonkarahisar", DisplayOrder: 10, SortOrder: 3},
	{ID: 4, CountryID: 1, Name: "Ağrı", Tag: "agri", DisplayOrder: 10, SortOrder: 4},
	{ID: 5, CountryID: 1, Name: "Amasya", Tag: "amasya", DisplayOrder: 10, SortOrder: 5},
	{ID: 6, CountryID: 1, Name: "Ankara", Tag: "ankara", DisplayOrder: 10, SortOrder: 6},
	{ID: 7, CountryID: 1, Name: "Antalya", Tag: "antalya", DisplayOrder: 10, SortOrder: 7},
	{ID: 8, CountryID: 1, Name: "Artvin", Tag: "artvin", DisplayOrder: 10, SortOrder: 8},
	{ID: 9, CountryID: 1, Name: "Aydın", Tag: "aydin", DisplayOrder: 10, SortOrder: 9},
	{ID: 10, CountryID: 1, Name: "Balıkesir", Tag: "balikesir", DisplayOrder: 10, SortOrder: 10},
	{ID: 11, CountryID: 1, Name: "Bilecik", Tag: "bilecik", DisplayOrder: 10, SortOrder: 11},
	{ID: 12, CountryID: 1, Name: "Bingöl", Tag: "bingol", DisplayOrder: 10, SortOrder: 12},
	{ID: 13, CountryID: 1, Name: "Bitlis", Tag: "bitlis", DisplayOrder: 10, SortOrder: 13},
	{ID: 14, CountryID: 1, Name: "Bolu", Tag: "bolu", DisplayOrder: 10, SortOrder: 14},
	{ID: 15, CountryID: 1, Name: "Burdur", Tag: "burdur", DisplayOrder: 10, SortOrder: 15},
	{ID: 16, CountryID: 1, Name: "Bursa", Tag: "bursa", DisplayOrder: 10, SortOrder: 16},
	{ID: 17, CountryID: 1, Name: "Çanakkale", Tag: "canakkale", DisplayOrder: 10, SortOrder: 17},
	{ID: 18, CountryID: 1, Name: "Çankırı", Tag: "cankiri", DisplayOrder: 10, SortOrder: 18},
	{ID: 19, CountryID: 1, Name: "Çorum", Tag: "corum", DisplayOrder: 10, SortOrder: 19},
	{ID: 20, CountryID: 1, Name: "Denizli", Tag: "denizli", DisplayOrder: 10, SortOrder: 20},
	{ID: 21, CountryID: 1, Name: "Diyarbakır", Tag: "diyarbakir", DisplayOrder: 10, SortOrder: 21},
	{ID: 22, CountryID: 1, Name: "Edirne", Tag: "edirne", DisplayOrder: 10, SortOrder: 22},
	{ID: 23, CountryID: 1, Name: "Elazığ", Tag: "elazig", DisplayOrder: 10, SortOrder: 23},
	{ID: 24, CountryID: 1, Name: "Erzincan", Tag: "erzincan", DisplayOrder: 10, SortOrder: 24},
	{ID: 25, CountryID: 1, Name: "Erzurum", Tag: "erzurum", DisplayOrder: 10, SortOrder: 25},
	{ID: 26, CountryID: 1, Name: "Eskişehir", Tag: "eskisehir", DisplayOrder: 10, SortOrder: 26},
	{ID: 27, CountryID: 1, Name: "Gaziantep", Tag: "gaziantep", DisplayOrder: 10, SortOrder: 27},
	{ID: 28, CountryID: 1, Name: "Giresun", Tag: "giresun", DisplayOrder: 10, SortOrder: 28},
	{ID: 29, CountryID: 1, Name: "Gümüşhane", Tag: "gumushane", DisplayOrder: 10, SortOrder: 29},
	{ID: 30, CountryID: 1, Name: "Hakkari", Tag: "hakkari", DisplayOrder: 10, SortOrder: 30},
	{ID: 31, CountryID: 1, Name: "Hatay", Tag: "hatay", DisplayOrder: 10, SortOrder: 31},
	{ID: 32, CountryID: 1, Name: "Isparta", Tag: "isparta", DisplayOrder: 10, SortOrder: 32},
	{ID: 33, CountryID: 1, Name: "Mersin", Tag: "mersin", DisplayOrder: 10, SortOrder: 33},
	{ID: 34, CountryID: 1, Name: "İstanbul", Tag: "istanbul", DisplayOrder: 10, SortOrder: 34},
	{ID: 35, CountryID: 1, Name: "İzmir", Tag: "izmir", DisplayOrder: 10, SortOrder: 35},
	{ID: 36, CountryID: 1, Name: "Kars", Tag: "kars", DisplayOrder: 10, SortOrder: 36},
	{ID: 37, CountryID: 1, Name: "Kastamonu", Tag: "kastamonu", DisplayOrder: 10, SortOrder: 37},
	{ID: 38, CountryID: 1, Name: "Kayseri", Tag: "kayseri", DisplayOrder: 10, SortOrder: 38},
	{ID: 39, CountryID: 1, Name: "Kırklareli", Tag: "kirklareli", DisplayOrder: 10, SortOrder: 39},
	{ID: 40, CountryID: 1, Name: "Kırşehir", Tag: "kirsehir", DisplayOrder: 10, SortOrder: 40},
	{ID: 41, CountryID: 1, Name: "Kocaeli", Tag: "kocaeli", DisplayOrder: 10, SortOrder: 41},
	{ID: 42, CountryID: 1, Name: "Konya", Tag: "konya", DisplayOrder: 10, SortOrder: 42},
	{ID: 43, CountryID: 1, Name: "Kütahya", Tag: "kutahya", DisplayOrder: 10, SortOrder: 43},
	{ID: 44, CountryID: 1, Name: "Malatya", Tag: "malatya", DisplayOrder: 10, SortOrder: 44},
	{ID: 45, CountryID: 1, Name: "Manisa", Tag: "manisa", DisplayOrder: 10, SortOrder: 45},
	{ID: 46, CountryID: 1, Name: "Kahramanmaraş", Tag: "kahramanmaras", DisplayOrder: 10, SortOrder: 46},
	{ID: 47, CountryID: 1, Name: "Mardin", Tag: "mardin", DisplayOrder: 10, SortOrder: 47},
	{ID: 48, CountryID: 1, Name: "Muğla", Tag: "mugla", DisplayOrder: 10, SortOrder: 48},
	{ID: 49, CountryID: 1, Name: "Muş", Tag: "mus", DisplayOrder: 10, SortOrder: 49},
	{ID: 50, CountryID: 1, Name: "Nevşehir", Tag: "nevsehir", DisplayOrder: 10, SortOrder: 50},
	{ID: 51, CountryID: 1, Name: "Niğde", Tag: "nigde", DisplayOrder: 10, SortOrder: 51},
	{ID: 52, CountryID: 1, Name: "Ordu", Tag: "ordu", DisplayOrder: 10, SortOrder: 52},
	{ID: 53, CountryID: 1, Name: "Rize", Tag: "rize", DisplayOrder: 10, SortOrder: 53},
	{ID: 54, CountryID: 1, Name: "Sakarya", Tag: "sakarya", DisplayOrder: 10, SortOrder: 54},
	{ID: 55, CountryID: 1, Name: "Samsun", Tag: "samsun", DisplayOrder: 10, SortOrder: 55},
	{ID: 56, CountryID: 1, Name: "Siirt", Tag: "siirt", DisplayOrder: 10, SortOrder: 56},
	{ID: 57, CountryID: 1, Name: "Sinop", Tag: "sinop", DisplayOrder: 10, SortOrder: 57},
	{ID: 58, CountryID: 1, Name: "Sivas", Tag: "sivas", DisplayOrder: 10, SortOrder: 58},
	{ID: 59, CountryID: 1, Name: "Tekirdağ", Tag: "tekirdag", DisplayOrder: 10, SortOrder: 59},
	{ID: 60, CountryID: 1, Name: "Tokat", Tag: "tokat", DisplayOrder: 10, SortOrder: 60},
	{ID: 61, CountryID: 1, Name: "Trabzon", Tag: "trabzon", DisplayOrder: 10, SortOrder: 61},
	{ID: 62, CountryID: 1, Name: "Tunceli", Tag: "tunceli", DisplayOrder: 10, SortOrder: 62},
	{ID: 63, CountryID: 1, Name: "Şanlıurfa", Tag: "sanliurfa", DisplayOrder: 10, SortOrder: 63},
	{ID: 64, CountryID: 1, Name: "Uşak", Tag: "usak", DisplayOrder: 10, SortOrder: 64},
	{ID: 65, CountryID: 1, Name: "Van", Tag: "van", DisplayOrder: 10, SortOrder: 65},
	{ID: 66, CountryID: 1, Name: "Yozgat", Tag: "yozgat", DisplayOrder: 10, SortOrder: 66},
	{ID: 67, CountryID: 1, Name: "Zonguldak", Tag: "zonguldak", DisplayOrder: 10, SortOrder: 67},
	{ID: 68, CountryID: 1, Name: "Aksaray", Tag: "aksaray", DisplayOrder: 10, SortOrder: 68},
	{ID: 69, CountryID: 1, Name: "Bayburt", Tag: "bayburt", DisplayOrder: 10, SortOrder: 69},
	{ID: 70, CountryID: 1, Name: "Karaman", Tag: "karaman", DisplayOrder: 10, SortOrder: 70},
	{ID: 71, CountryID: 1, Name: "Kırıkkale", Tag: "kirikkale", DisplayOrder: 10, SortOrder: 71},
	{ID: 72, CountryID: 1, Name: "Batman", Tag: "batman", DisplayOrder: 10, SortOrder: 72},
	{ID: 73, CountryID: 1, Name: "Şırnak", Tag: "sirnak", DisplayOrder: 10, SortOrder: 73},
	{ID: 74, CountryID: 1, Name: "Bartın", Tag: "bartin", DisplayOrder: 10, SortOrder: 74},
	{ID: 75, CountryID: 1, Name: "Ardahan", Tag: "ardahan", DisplayOrder: 10, SortOrder: 75},
	{ID: 76, CountryID: 1, Name: "Iğdır", Tag: "igdir", DisplayOrder: 10, SortOrder: 76},
	{ID: 77, CountryID: 1, Name: "Yalova", Tag: "yalova", DisplayOrder: 10, SortOrder: 77},
	{ID: 78, CountryID: 1, Name: "Karabük", Tag: "karabuk", DisplayOrder: 10, SortOrder: 78},
	{ID: 79, CountryID: 1, Name: "Kilis", Tag: "kilis", DisplayOrder: 10, SortOrder: 79},
	{ID: 80, CountryID: 1, Name: "Osmaniye", Tag: "osmaniye", DisplayOrder: 10, SortOrder: 80},
	{ID: 81, CountryID: 1, Name: "Düzce", Tag: "duzce", DisplayOrder: 10, SortOrder: 81},
}

func GetTownsSQLFile() (string, error) {
	s, err := os.ReadFile("sql/towns.sql")
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func GetDistrictsSQLFile() (string, error) {
	s, err := os.ReadFile("sql/districts.sql")
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func GetQuartersSQLFile() (string, error) {
	s, err := os.ReadFile("sql/quarters.sql")
	if err != nil {
		return "", err
	}
	return string(s), nil
}
