package mailform

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// DefaultBaseURL is the default mailform API base url, but is can be overwritten via Config
	DefaultBaseURL = "https://www.mailform.io/app/api/v1"
	DefaultTimeout = time.Second * 15
	// Order Statuses
	StatusCancelled           = "cancelled"
	StatusQueued              = "queued"
	StatusAwaitingFulfillment = "awaiting_fulfillment"
	StatusFulfilled           = "fulfilled"
)

var (
	// ErrNilConfig is returned when a nil config is being passed to New().
	ErrNilConfig = errors.New("config cannot be nil")
)

// Client is the mailform REST API client.
type Client struct {
	restClient *resty.Client
}

// Config is the configuration used to communicate with the mailform API.
type Config struct {
	Token   string
	BaseURL string
	Timeout time.Duration
}

// ErrMailform is the error returned when mailform responds with an error.
type ErrMailform struct {
	Err struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Detail string `json:"detail"`
}

func (e *ErrMailform) Error() string {
	// Return detailed error message if populated
	if e.Detail != "" {
		return e.Detail
	}

	// Otherwise send message
	return e.Err.Message
}

// New returns a new mailform REST API client.
func New(c *Config) (*Client, error) {
	if c == nil {
		return nil, ErrNilConfig
	}

	// Check for API token
	// Allow consumer to override default baseURL if needed
	baseURL := DefaultBaseURL
	if c.BaseURL != "" {
		baseURL = c.BaseURL
	}

	// Allow consumer to override default timeout if needed
	timeout := DefaultTimeout
	if c.Timeout.Seconds() != 0 {
		timeout = c.Timeout
	}

	// Create new client for mailform.io
	mailformClient := &Client{
		restClient: resty.New().
			SetBaseURL(baseURL).
			SetTimeout(timeout).
			SetAuthToken(c.Token),
	}

	return mailformClient, nil
}

// checkBodyForErr ensures the response from mailform isn't actually an error.
// cause we can get 200, failed successfully 🙄
func checkBodyForErr(b []byte) error {
	mailformErr := &ErrMailform{}

	err := json.Unmarshal(b, mailformErr)
	if err != nil {
		return err
	}

	if mailformErr.Err.Message != "" {
		return mailformErr
	}

	return nil
}
