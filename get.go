package threecx

import "fmt"

func (c *Client) Get(uri string, result *any, query map[string]string) error {
	resp, err := c.rest.R().SetQueryParams(query).SetResult(result).Get(uri)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("%s", resp.String())
	}

	return nil
}
