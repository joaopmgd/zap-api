package model

import "time"

// Property data model
type Property struct {
	UsableAreas   int          `json:"usableAreas"`
	ListingType   string       `json:"listingType"`
	CreatedAt     time.Time    `json:"createdAt"`
	ListingStatus string       `json:"listingStatus"`
	Id            string       `json:"id"`
	ParkingSpaces int          `json:"parkingSpaces"`
	UpdatedAt     time.Time    `json:"updatedAt"`
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
