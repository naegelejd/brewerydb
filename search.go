package brewerydb

type SearchService struct {
	c *Client
}

type SearchResults struct {
	c *Client
}

func (ss *SearchService) Search( /* params */ ) (*SearchResults, error) {
	return nil, nil
}
