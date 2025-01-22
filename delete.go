package threecx

import "fmt"

func (c *Client) Delete(uri string, query map[string]string) error {
	resp, err := c.client.
		R().
		SetAuthToken(c.token.AccessToken).
		SetQueryParams(query).
		Delete(uri)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("%s", resp.String())
	}

	return nil
}
