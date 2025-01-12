package threecx

import "fmt"

func (c *Client) Delete(uri string, query map[string]string) error {
	resp, err := c.rest.
		R().
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
