package httpclient

// defaultContentType defines the default content type for requests.
const defaultContentType = "application/json"

type Headers map[string]string

// requestParams represents parameters for an HTTP request.
type requestParams struct {
	headers     Headers
	contentType string
}

// newRequestParams creates a new requestParams instance with the provided options.
func newRequestParams(opts ...RequestOption) (*requestParams, error) {
	rp := &requestParams{
		contentType: defaultContentType,
	}

	for _, opt := range opts {
		err := opt(rp)
		if err != nil {
			return nil, err
		}
	}

	return rp, nil
}

// RequestOption represents an option for configuring request parameters.
type RequestOption func(*requestParams) error

// WithHeaders sets the HTTP headers for the request.
func WithHeaders(headers Headers) RequestOption {
	return func(r *requestParams) error {
		r.headers = headers

		return nil
	}
}
