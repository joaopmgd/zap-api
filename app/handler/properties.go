package handler

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	cache "github.com/patrickmn/go-cache"
	"gitlab.com/zap-api/app/model"
	"gitlab.com/zap-api/config"
)

// GetAllProperties will recover all Properties for the requested source
func GetAllProperties(config *config.Config, w http.ResponseWriter, r *http.Request) {
	source := r.Header.Get("source")
	config.Logger.Info("Recovering Properties for", source)
	datasources := *config.Datasources
	properties := &[]model.Property{}
	if _, ok := datasources[source]; !ok {
		config.Logger.Error("No property found for", source)
		respondError(w, http.StatusNotFound, "Source not accepted.")
		return
	}
	propertiesCachedString, found := config.Cache.Get(source)
	if !found {
		config.Logger.Info("There is no cache for this request.")
		config.Logger.Info("Requesting data from ZAP")
		properties = getPropertiesOr404(config.Endpoints.ZapProperties, w, r)
		if properties == nil {
			config.Logger.Error("Request NOT FOUND")
			return
		}
		properties = setCacheProperties(source, properties, config)
	} else {
		config.Logger.Info("Found a cache for this request")
		properties = propertiesCachedString.(*[]model.Property)
	}
	respondJSON(w, http.StatusOK, paginate(config, r, properties))
}

// getPropertiesOr404 gets all properties, or respond the 404 error otherwise
func getPropertiesOr404(url string, w http.ResponseWriter, r *http.Request) *[]model.Property {
	properties := []model.Property{}
	if err := requestProperties(&properties, url); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &properties
}

// the Timeout is 100s because it can take too long to make the first request without cache
func requestProperties(target interface{}, url string) error {
	var myClient = &http.Client{Timeout: 100 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

// Paginate just picksup a slice from the Response, showing just the page Requested
func paginate(config *config.Config, r *http.Request, properties *[]model.Property) *model.ListPropertyResponse {
	config.Logger.Info("Paginating the Response")
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	if offset*limit > len(*properties)-1 {
		config.Logger.Error("Offset bigger than the Response")
		return &model.ListPropertyResponse{}
	}
	if offset+limit > len(*properties) {
		limit = len(*properties)
	}
	test := (*properties)[offset*limit : (offset*limit)+limit]
	listPropertyResponse := model.ListPropertyResponse{
		Properties:           test,
		PageNumber:           offset,
		PageSize:             limit,
		PropertiesTotalCount: len(*properties),
	}
	return &listPropertyResponse
}

// setCacheProperties will create a cache for the possible requests
// since the JSON returned is too big, the next requests will all be recovered by the cache
// the first request will analyze each property, distribute in different caches and return just the selected source
func setCacheProperties(source string, properties *[]model.Property, config *config.Config) *[]model.Property {
	config.Logger.Info("Setting up the Response Cache for future Requests.")
	zapProperties := []model.Property{}
	vivaRealProperties := []model.Property{}
	for _, property := range *properties {
		businessType := property.PricingInfos.BusinessType
		price, err := strconv.Atoi(property.PricingInfos.Price)
		// Price is converted before everything, if it fails the property is rejected early
		// Verifies if lat and lon are 0, if they are the property is rejected
		if err == nil && !isLatLonValid(property.Address.GeoLocation.Location) {
			if isZapPropertyValid(property.UsableAreas, price, businessType) {
				if isInsideBoundingBox(property.Address.GeoLocation.Location.Lat, property.Address.GeoLocation.Location.Lon) &&
					businessType == "SALE" {
					property.PricingInfos.Price = strconv.FormatFloat((float64(price) * (0.9)), 'f', 6, 64)
					property.UpdatedAt = time.Now().Format(time.RFC3339)
				}
				zapProperties = append(zapProperties, property)
			}
			if isVivaRealPropertyValid(property.PricingInfos.MonthlyCondoFee, price, businessType) {
				if isInsideBoundingBox(property.Address.GeoLocation.Location.Lat, property.Address.GeoLocation.Location.Lon) &&
					businessType == "RENTAL" {
					property.PricingInfos.Price = strconv.FormatFloat((float64(price) * (1.5)), 'f', 6, 64)
					property.UpdatedAt = time.Now().Format(time.RFC3339)
				}
				vivaRealProperties = append(vivaRealProperties, property)
			}
		}
	}
	config.Cache.Set("zap", &zapProperties, cache.DefaultExpiration)
	config.Cache.Set("vivareal", &vivaRealProperties, cache.DefaultExpiration)
	config.Logger.Info("Responding with Properties for.", source)
	if source == "zap" {
		return &zapProperties
	}
	return &vivaRealProperties
}

func isZapPropertyValid(usableAreas int, price int, businessType string) bool {
	acceptableSquareMeter := true
	if businessType == "RENTAL" {
		acceptableSquareMeter = isSquareMeterAcceptable(price, usableAreas)
	}
	salePrice := businessType == "SALE" && price >= 600000
	rentalPrice := businessType == "RENTAL" && price >= 3500
	return (salePrice || rentalPrice) && acceptableSquareMeter
}

// Verifies if the square meter value is acceptable for Zap
func isSquareMeterAcceptable(price, usableAreas int) bool {
	if usableAreas > 0 {
		return (price / usableAreas) > 3500
	}
	// It must consider just the usableAreas above 0
	return false
}

// Verifies the monthlyCondoFee againts 30% of the Rental price
// both must be parsed to float because the comparison is made with 30% of the price
// And verifies the old rules of minimum price for SALE and RENTAL
func isVivaRealPropertyValid(monthlyCondoFee string, price int, businessType string) bool {
	rentalPriceCondoComparison := true
	if businessType == "RENTAL" {
		rentalPriceCondoComparison = isRentalPriceAndCondoFeeAcceptable(monthlyCondoFee, price)
	}
	salePrice := businessType == "SALE" && price >= 700000
	rentalPrice := businessType == "RENTAL" && price >= 4000
	return (salePrice || rentalPrice) && rentalPriceCondoComparison
}

func isRentalPriceAndCondoFeeAcceptable(propertyMonthlyCondoFee string, rentalPrice int) bool {
	monthlyCondoFee, err := strconv.ParseFloat(propertyMonthlyCondoFee, 64)
	// If the value is not valid and returns an error, then it is valid
	if err == nil {
		if monthlyCondoFee < float64(rentalPrice)*0.3 {
			return true
		}
		// If price is not a real value, then it will not able to compare with the monthlyCondoFee
		return false
	}
	return true
}

// Checks if Lat and Long are equal to 0
func isLatLonValid(location model.Location) bool {
	if location.Lat == 0 && location.Lon == 0 {
		return true
	}
	return false
}

func isInsideBoundingBox(x, y float64) bool {
	// Describing the points by the coordinates
	x1 := model.VivaRealBoundBox.Maxlat
	y1 := model.VivaRealBoundBox.Minlon

	x2 := model.VivaRealBoundBox.Maxlat
	y2 := model.VivaRealBoundBox.Maxlon

	x3 := model.VivaRealBoundBox.Minlat
	y3 := model.VivaRealBoundBox.Maxlon

	x4 := model.VivaRealBoundBox.Minlat
	y4 := model.VivaRealBoundBox.Minlon

	// Calculating all the areas
	// A = area total
	A := triangleArea(x1, y1, x2, y2, x3, y3) + triangleArea(x1, y1, x4, y4, x3, y3)

	// Areas of every triangle using the new point
	A1 := triangleArea(x, y, x1, y1, x2, y2)
	A2 := triangleArea(x, y, x2, y2, x3, y3)
	A3 := triangleArea(x, y, x3, y3, x4, y4)
	A4 := triangleArea(x, y, x1, y1, x4, y4)
	// Check if sum of A1, A2, A3 and A4 is same as A
	// If the sum of the areas using the new point is the same as the area of the square
	// Then the new point is inside the square
	return A == (A1 + A2 + A3 + A4)
}

func triangleArea(x1, y1, x2, y2, x3, y3 float64) float64 {
	return math.Abs(Round((Round(x1*(y2-y3)) + Round(x2*(y3-y1)) + Round(x3*(y1-y2))) / 2.0))
}

// Round every step to 6 decimal places in the float value
func Round(x float64) float64 {
	return math.Round(x/0.0000005) * 0.0000005
}
