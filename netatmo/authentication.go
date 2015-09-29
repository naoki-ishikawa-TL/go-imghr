package netatmo

import (
    "net/http"
    "net/url"
    "encoding/json"
    "time"
)

type NetatmoAuthenticator struct {
    ClientId string
    ClientSecret string
    AccessToken string
    ExpiresIn int
    RefreshToken string
    LastUpdate time.Time
    ExpireAt time.Time
}

type tokenResponse struct {
    AccessToken string `json:"access_token"`
    ExpiresIn int `json:"expires_in"`
    RefreshToken string `json:"refresh_token"`
}

func NewNetatmoAuthenticator(clientId string, clientSecret string) *NetatmoAuthenticator {
    return &NetatmoAuthenticator{ClientId: clientId, ClientSecret: clientSecret}
}

func (this *NetatmoAuthenticator) GetClientCredentialAccessToken(username string, password string, scope string) chan string {
    accessTokenChan := make(chan string)

    go func() {
        if this.AccessToken == "" {
            this.authenticateWithClientCredential(username, password, scope)
        }
        if this.IsExpire() == true {
            this.refreshAccessToken()
        }
        accessTokenChan <-this.AccessToken
        this.LastUpdate = time.Now()
        this.ExpireAt = time.Now().Add(time.Duration(this.ExpiresIn) * time.Second)
    }()

    return accessTokenChan
}

func (this *NetatmoAuthenticator) IsExpire() bool {
    if time.Now().Add(10 * time.Second).Unix() >= this.ExpireAt.Unix() {
        return true
    }
    return false
}

func (this *NetatmoAuthenticator) authenticateWithClientCredential(username string, password string, scope string) {
	v := url.Values{}
	v.Set("grant_type", "password")
	v.Set("client_id", this.ClientId)
	v.Set("client_secret", this.ClientSecret)
	v.Set("username", username)
    v.Set("password", password)
    v.Set("scope", scope)
    response, _ := http.PostForm("https://api.netatmo.net/oauth2/token", v)
    dec := json.NewDecoder(response.Body)
    var data tokenResponse
    dec.Decode(&data)
    this.AccessToken = data.AccessToken
    this.ExpiresIn = data.ExpiresIn
    this.RefreshToken = data.RefreshToken
}

func (this *NetatmoAuthenticator) refreshAccessToken() {
	v := url.Values{}
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", this.RefreshToken)
	v.Set("client_id", this.ClientId)
	v.Set("client_secret", this.ClientSecret)
    response, _ := http.PostForm("https://api.netatmo.net/oauth2/token", v)
    dec := json.NewDecoder(response.Body)
    var data tokenResponse
    dec.Decode(&data)
    this.AccessToken = data.AccessToken
    this.ExpiresIn = data.ExpiresIn
    this.RefreshToken = data.RefreshToken
}
