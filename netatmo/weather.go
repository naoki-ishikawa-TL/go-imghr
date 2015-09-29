package netatmo

import (
    "net/http"
    "net/url"
    "encoding/json"
)

type publicData struct {
    Status string
    Body []struct {
        Id string `json:"_id"`
        Place struct {
            Location []float32
            Alttude float32
            Timezone string
        }
        Mark int
        Measures interface{}
    }
}

func GetTemperatureAndHumidity(accessToken string, latitudeNe string, longitudeNe string, latitudeSw string, longitudeSw string) float32, int {
    v := url.Values{}
    v.Set("access_token", access_token)
    v.Set("lat_ne", latitudeNe)
    v.Set("lon_ne", longitudeNe)
    v.Set("lat_sw", latitudeSw)
    v.Set("lon_sw", longitudeSw)
    response, _ := http.Get("https://api.netatmo.com/api/getpublicdata"+v.Encode())
    dec := json.NewDecoder(response.Body)
    var data publicData
    dec.Decode(&data)
}
