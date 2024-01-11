package v1

var SortsMap = map[string]string{
	"id":          "id",
	"name":        "name",
	"description": "description",
	"type":        "type",
	"make":        "vehicle_make",
	"model":       "vehicle_model",
	"year":        "vehicle_year",
	"length":      "vehicle_length",
	"sleeps":      "sleeps",
	"price":       "price_per_day",
	"city":        "home_city",
	"state":       "home_state",
	"zip":         "home_zip",
	"country":     "home_country",
}

type Rental struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Type            string   `json:"type"`
	Make            string   `json:"make"`
	Model           string   `json:"model"`
	Year            int      `json:"year"`
	Length          float32  `json:"length"`
	Sleeps          int      `json:"sleeps"`
	PrimaryImageURL string   `json:"primary_image_url"`
	Price           Price    `json:"price"`
	Location        Location `json:"location"`
	User            User     `json:"user"`
}

type Price struct {
	Day int `json:"day"`
}

type Location struct {
	City    string  `json:"city"`
	State   string  `json:"state"`
	Zip     string  `json:"zip"`
	Country string  `json:"country"`
	Lat     float32 `json:"lat"`
	Lng     float32 `json:"lng"`
}
