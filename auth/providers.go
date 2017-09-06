package auth

import (
	"encoding/json"
	"io/ioutil"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type Provider interface {
	SetRedirectURL(string)
	ID() string
	Name() string
	AuthURL(string) string
	Profile(string) (Profile, error)
}

func getData(config *oauth2.Config, code, url string) ([]byte, error) {
	oauthToken, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	client := config.Client(oauth2.NoContext, oauthToken)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

////////////////

type GoogleProvider struct {
	*oauth2.Config
}

func Google(clientID, clientSecret string) GoogleProvider {
	return GoogleProvider{
		&oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (p GoogleProvider) SetRedirectURL(redirectURL string) {
	p.RedirectURL = redirectURL
}

func (p GoogleProvider) ID() string {
	return "google"
}

func (p GoogleProvider) Name() string {
	return "Google"
}

func (p GoogleProvider) AuthURL(state string) string {
	return p.AuthCodeURL(state)
}

func (p GoogleProvider) Profile(code string) (Profile, error) {
	data, err := getData(p.Config, code, "https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return Profile{}, err
	}

	v := &struct {
		Name string `json:"name"`
	}{}
	if err := json.Unmarshal(data, v); err != nil {
		return Profile{}, err
	}
	return Profile{v.Name, "email"}, nil
}

////////////////

type FacebookProvider struct {
	*oauth2.Config
}

func Facebook(clientID, clientSecret string) FacebookProvider {
	return FacebookProvider{
		&oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"public_profile"},
			Endpoint:     facebook.Endpoint,
		},
	}
}

func (p FacebookProvider) SetRedirectURL(redirectURL string) {
	p.RedirectURL = redirectURL
}

func (p FacebookProvider) ID() string {
	return "facebook"
}

func (p FacebookProvider) Name() string {
	return "Facebook"
}

func (p FacebookProvider) AuthURL(state string) string {
	return p.AuthCodeURL(state)
}

func (p FacebookProvider) Profile(code string) (Profile, error) {
	data, err := getData(p.Config, code, "https://graph.facebook.com/me")
	if err != nil {
		return Profile{}, err
	}

	v := &struct {
		Name string `json:"name"`
	}{}
	if err := json.Unmarshal(data, v); err != nil {
		return Profile{}, err
	}
	return Profile{v.Name, "email"}, nil
}

////////////////

type GitHubProvider struct {
	*oauth2.Config
}

func GitHub(clientID, clientSecret string) GitHubProvider {
	return GitHubProvider{
		&oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"user"},
			Endpoint:     github.Endpoint,
		},
	}
}

func (p GitHubProvider) SetRedirectURL(redirectURL string) {
	p.RedirectURL = redirectURL
}

func (p GitHubProvider) ID() string {
	return "github"
}

func (p GitHubProvider) Name() string {
	return "GitHub"
}

func (p GitHubProvider) AuthURL(state string) string {
	return p.AuthCodeURL(state)
}

func (p GitHubProvider) Profile(code string) (Profile, error) {
	data, err := getData(p.Config, code, "https://api.github.com/user")
	if err != nil {
		return Profile{}, err
	}

	v := &struct {
		Name string `json:"name"`
	}{}
	if err := json.Unmarshal(data, v); err != nil {
		return Profile{}, err
	}
	return Profile{v.Name, "email"}, nil
}
