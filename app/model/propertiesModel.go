package model

// Property data model
type Property struct {
	UsableAreas   int          `json:"usableAreas"`
	ListingType   string       `json:"listingType"`
	CreatedAt     string       `json:"createdAt"`
	ListingStatus string       `json:"listingStatus"`
	Id            string       `json:"id"`
	ParkingSpaces int          `json:"parkingSpaces"`
	UpdatedAt     string       `json:"updatedAt"`
	Owner         bool         `json:"owner"`
	Images        []string     `json:"images"`
	Address       Address      `json:"address"`
	Bathrooms     int          `json:"bathrooms"`
	Bedrooms      int          `json:"bedrooms"`
	PricingInfos  PricingInfos `json:"pricingInfos"`
}

type Address struct {
	City         string      `json:"city"`
	Neighborhood string      `json:"neighborhood"`
	GeoLocation  GeoLocation `json:"geoLocation"`
}

type GeoLocation struct {
	Precision string   `json:"precision"`
	Location  Location `json:"location"`
}

type Location struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type PricingInfos struct {
	YearlyIptu       string `json:"yearlyIptu"`
	Price            string `json:"price"`
	BusinessType     string `json:"businessType"`
	MonthlyCondoFee  string `json:"monthlyCondoFee"`
	Period           string `json:"period"`
	RentalTotalPrice string `json:"rentalTotalPrice"`
}

type ListPropertyResponse struct {
	Properties           []Property `json:"listings"`
	PageNumber           int        `json:"pageNumber"`
	PageSize             int        `json:"pageSize"`
	PropertiesTotalCount int        `json:"propertiestotalCount"`
}

type BoundingBox struct {
	Minlon float64
	Minlat float64
	Maxlon float64
	Maxlat float64
}

var VivaRealBoundBox = BoundingBox{
	Minlon: -46.693419,
	Minlat: -23.568704,
	Maxlon: -46.641146,
	Maxlat: -23.546686,
}
