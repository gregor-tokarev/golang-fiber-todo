package provider

type Provider interface {
	GetAuthUrl() (string, error)
	GetTokens() (map[string]interface{}, error)
	AuthRefresh() error
	FetchInfo() error
}

//var Providers = map[string]Provider{
//	"google": NewGoogleProvider(),
//}

type ProviderUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
