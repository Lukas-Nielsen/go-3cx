package tcx

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Host struct {
	// 3cx FQDN
	FQDN string
	// 3cx port
	Port int
	// debug the requests (see https://github.com/go-resty/resty)
	Debug bool
}

type Token struct {
	// normaly Bearer
	TokenType string `json:"token_type"`
	// unix timestamp when token expires
	Expires      int64  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	// username
	Username string
	// password
	Password string
	// mfa secret
	MFA string
}

type Rest struct {
	// client id
	ClientID string
	// client secret
	ClientSecret string
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

type Client struct {
	host   Host
	token  Token
	user   User
	rest   Rest
	client *resty.Client
}

func NewClient(host Host) (*Client, error) {
	c := Client{
		host: host,
	}

	c = *c.setup()

	return &c, nil
}

func (c *Client) SetHost(host Host) (*Client, error) {
	c.host = host

	return c.setup(), nil
}

func (c *Client) SetUser(user User) (*Client, error) {
	c.user = user
	MFA, _ := getOTP(c.user.MFA)

	var res authResponse

	resp, err := c.client.
		R().
		SetResult(&res).
		SetBody(
			authRequest{
				Username:     c.user.Username,
				Password:     c.user.Password,
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
	return c, nil
}

func (c *Client) SetRest(rest Rest) (*Client, error) {
	c.rest = rest

	var res Token

	resp, err := c.client.
		R().
		SetResult(&res).
		SetFormData(
			map[string]string{
				"client_id":     c.rest.ClientID,
				"client_secret": c.rest.ClientSecret,
				"grant_type":    "client_credentials",
			},
		).
		Post("/connect/token")

	if err != nil {
		return &Client{}, err
	}

	if resp.IsError() {
		return &Client{}, fmt.Errorf("%s", "error during client_credentials login")
	}

	c.token = res
	c.token.Expires = time.Now().Unix() + res.Expires*60*1000
	return c, nil
}

func (c *Client) SetToken(token Token) (*Client, error) {
	c.token = token
	if c.token.Expires < time.Now().Unix() {
		var res authResponse

		resp, err := c.client.
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
	return c, nil
}

func (c *Client) setup() *Client {
	c.client = resty.
		New().
		SetDebug(c.host.Debug).
		SetBaseURL(
			fmt.Sprintf(
				"https://%s:%d",
				c.host.FQDN,
				c.host.Port,
			),
		)
	return c
}
