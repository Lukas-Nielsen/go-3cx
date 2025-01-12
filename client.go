package threecx

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type ClientConfig struct {
	// 3cx FQDN
	FQDN string
	// 3cx port
	Port     int
	User     string
	Passwort string
	// MFA client secret
	MFA string
	// debug the requests (see https://github.com/go-resty/resty)
	Debug bool
	Token Token
}

type authRequest struct {
	Username     string `json:"Username"`
	Password     string `json:"Password"`
	SecurityCode string `json:"SecurityCode"`
}

type authResponse struct {
	Status        string `json:"Status"`
	Token         Token  `json:"Token"`
	TwoFactorAuth any    `json:"TwoFactorAuth"`
}

type Token struct {
	TokenType    string `json:"token_type"`
	Expires      int64  `json:"expires_in"` // unix timestamp when token expires
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Client struct {
	config    ClientConfig
	token     Token
	rest      *resty.Client
	tokenAuth bool
}

func NewClient(config ClientConfig) (*Client, error) {
	c := Client{
		config: config,
	}

	if c.config.Port == 0 {
		c.config.Port = 443
	}

	if len(c.config.FQDN) == 0 {
		return &Client{}, fmt.Errorf("%s", "missing FQDN")
	}

	if len(c.config.Token.TokenType) == 0 ||
		len(c.config.Token.AccessToken) == 0 ||
		len(c.config.Token.RefreshToken) == 0 ||
		c.config.Token.Expires == 0 {

		if len(c.config.User) == 0 {
			return &Client{}, fmt.Errorf("%s", "missing User")
		}

		if len(c.config.Passwort) == 0 {
			return &Client{}, fmt.Errorf("%s", "missing Password")
		}
	} else {
		c.tokenAuth = true
	}

	c.rest = resty.
		New().
		SetDebug(c.config.Debug).
		SetBaseURL(
			fmt.Sprintf(
				"https://%s:%d",
				c.config.FQDN,
				c.config.Port,
			),
		)

	MFA, _ := getOTP(c.config.MFA)

	if !c.tokenAuth {
		var res authResponse

		resp, err := c.rest.
			R().
			SetResult(&res).
			SetBody(
				authRequest{
					Username:     c.config.User,
					Password:     c.config.Passwort,
					SecurityCode: MFA,
				},
			).
			Post("/webclient/api/Login/GetAccessToken")

		if err != nil {
			return &Client{}, err
		}

		if resp.IsError() {
			return &Client{}, fmt.Errorf("%s", "error during login")
		}

		c.token = res.Token
		c.token.Expires = time.Now().Unix() + res.Token.Expires*60*1000
	} else {
		c.token = c.config.Token
	}

	if c.token.Expires < time.Now().Unix() {
		var res authResponse
		resp, err := c.rest.
			R().
			SetResult(&res).
			SetFormData(
				map[string]string{
					"client_id":     "go-3cx",
					"grant_type":    "refresh_token",
					"refresh_token": c.token.RefreshToken,
				},
			).
			Post("/connect/token")

		if err != nil {
			return &Client{}, err
		}

		if resp.IsError() {
			return &Client{}, fmt.Errorf("%s", "error during token refresh")
		}

		c.token = res.Token
		c.token.Expires = time.Now().Unix() + res.Token.Expires*60*1000
	}

	c.rest.
		SetAuthToken(c.token.AccessToken).
		SetBaseURL(
			fmt.Sprintf(
				"https://%s:%d/xapi/v1",
				c.config.FQDN,
				c.config.Port,
			),
		)

	return &c, nil
}
