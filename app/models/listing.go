package models

type Listing struct {
	Kind string      `json:"kind"`
	Data ListingData `json:"data"`
}

type ListingData struct {
	After     *string `json:"after"`
	Dist      *int    `json:"dist"`
	ModHash   *string `json:"modhash"`
	GeoFilter *string `json:"geo_filter"`
	Children  []Post  `json:"children"`
	Before    *string `json:"before"`
}
