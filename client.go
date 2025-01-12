package threecx

import (
	"fmt"

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
}

type authRequest struct {
	Username     string `json:"Username"`
	Password     string `json:"Password"`
	SecurityCode string `json:"SecurityCode"`
}

type authResponse struct {
	Status        string `json:"Status"`
	Token         token  `json:"Token"`
	TwoFactorAuth any    `json:"TwoFactorAuth"`
}

type token struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Client struct {
	config ClientConfig
	token  token
	rest   *resty.Client
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

	if len(c.config.User) == 0 {
		return &Client{}, fmt.Errorf("%s", "missing User")
	}

	if len(c.config.Passwort) == 0 {
		return &Client{}, fmt.Errorf("%s", "missing Password")
	}

	c.rest = resty.New().SetDebug(c.config.Debug).SetBaseURL(fmt.Sprintf("https://%s:%d", c.config.FQDN, c.config.Port))

	MFA, _ := getOTP(c.config.MFA)

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
