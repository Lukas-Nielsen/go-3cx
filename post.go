package threecx

import "fmt"

func (c *Client) Post(uri string, payload any, query map[string]string) error {
	resp, err := c.rest.
		R().
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
