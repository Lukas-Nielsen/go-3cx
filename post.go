package threecx

import "fmt"

func (c *Client) Post(uri string, payload any, query map[string]string) error {
	resp, err := c.client.
		R().
		SetAuthToken(c.token.AccessToken).
		SetQueryParams(query).
		SetBody(payload).
		Post(uri)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("%s", resp.String())
	}

	return nil
}
