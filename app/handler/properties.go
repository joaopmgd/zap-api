package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	cache "github.com/patrickmn/go-cache"
	"gitlab.com/api/app/model"
	"gitlab.com/api/config"
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
func paginate(config *config.Config, r *http.Request, properties *[]model.Property) *[]model.Property {
	config.Logger.Info("Paginating the Response")
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}
	if offset > len(*properties)-1 {
		config.Logger.Error("Offset bigger than the Response")
		return &[]model.Property{}
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	limit = offset + limit
	if limit > len(*properties) {
		limit = len(*properties)
	}
	test := (*properties)[offset:limit]
	return &test
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
		if err == nil {
			if (businessType == "SALE" && price >= 600000) ||
				(businessType == "RENTAL" && price >= 3500) {
				zapProperties = append(zapProperties, property)
			}
			if (businessType == "SALE" && price >= 700000) ||
				(businessType == "RENTAL" && price >= 4000) {
				vivaRealProperties = append(vivaRealProperties, property)
			}
		} else {
			config.Logger.Error("Property with wrong price.", err)
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
