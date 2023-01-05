package client

type HTTPClient interface {
	// Get get request
	Get(uri string) ([]byte, error)
}
