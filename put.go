package threecx

import "fmt"

func (c *Client) Put(uri string, payload any, query map[string]string) error {
	resp, err := c.rest.R().SetQueryParams(query).SetBody(payload).Put(uri)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("%s", resp.String())
	}

	return nil
}
