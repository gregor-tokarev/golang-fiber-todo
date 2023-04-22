package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type GoogleProvider struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
	State        string `json:"state"`
}

type GoogleProviderConfig struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

type GoogleTokens struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IdToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func NewGoogleProvider(config *GoogleProviderConfig) *GoogleProvider {
	return &GoogleProvider{
		ClientId:     config.ClientId,
		ClientSecret: config.ClientSecret,
		RedirectUri:  "http://localhost:3000/api/v1/auth/google/callback",
		Scope:        config.Scope,
		State:        "",
	}
}

func (g GoogleProvider) GetAuthUrl() (string, error) {
	if g.ClientId == "" || g.ClientSecret == "" || g.RedirectUri == "" || g.Scope == "" {
		return "", errors.New("invalid config")
	}

	return fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?scope=%s&access_type=offline&prompt=consent&include_granted_scopes=true&client_id=%s&redirect_uri=%s&response_type=code",
		g.Scope,
		g.ClientId,
		g.RedirectUri,
	), nil
}

func (g GoogleProvider) GetTokens(code string) (*GoogleTokens, error) {

	url := fmt.Sprintf(
		"https://oauth2.googleapis.com/token?code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		code,
		g.ClientId,
		g.ClientSecret,
		g.RedirectUri,
	)
	res, err := http.Post(
		url,
		"application/x-www-form-urlencoded",
		nil,
	)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var mapRes *GoogleTokens
	err = json.Unmarshal(body, &mapRes)
	if err != nil {
		return nil, err
	}

	return mapRes, nil
}

func (g GoogleProvider) FetchInfo(accessToken string) (*ProviderUser, error) {
	if accessToken == "" {
		return nil, errors.New("invalid access token")
	}

	url := "https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=" + accessToken
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	fmt.Println(string(body))

	var mapRes *ProviderUser
	err = json.Unmarshal(body, &mapRes)
	if err != nil {
		return nil, err
	}

	return mapRes, nil
}

//func (g GoogleProvider) FetchInfo(url string) (string, error) {
//
//}
