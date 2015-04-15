package brewerydb

type SearchResults struct {
	c *Client
}

func (c *Client) Search( /* params */ ) (*SearchResults, error) {
	return nil, nil
}

func (sr *SearchResults) First() error {
	return nil
}

func (sr *SearchResults) Next() error {
	return nil
}
