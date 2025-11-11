package mailform

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// DefaultAPIBaseURL is the default mailform API base url, but is can be overwritten via Config
	DefaultAPIBaseURL = "https://www.mailform.io/app/api/v1"
	// DefaultAppBaseURL is used for other endpoints like cancel... for some reason
	DefaultAppBaseURL = "https://www.mailform.io/app/v1"
	DefaultTimeout    = time.Second * 15
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
	restClient         *resty.Client
	cancellationClient *resty.Client
}

// Config is the configuration used to communicate with the mailform API.
type Config struct {
	Token      string
	APIBaseURL string
	AppBaseURL string
	Timeout    time.Duration
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
	// Allow consumer to override default baseURL(s) if needed
	apiBaseURL := DefaultAPIBaseURL
	if c.APIBaseURL != "" {
		apiBaseURL = c.APIBaseURL
	}

	appBaseURL := DefaultAppBaseURL
	if c.AppBaseURL != "" {
		appBaseURL = c.AppBaseURL
	}

	// Allow consumer to override default timeout if needed
	timeout := DefaultTimeout
	if c.Timeout.Seconds() != 0 {
		timeout = c.Timeout
	}

	// Create new client(s) for mailform.io
	mailformClient := &Client{
		restClient: resty.New().
			SetBaseURL(apiBaseURL).
			SetTimeout(timeout).
			SetAuthToken(c.Token),
		cancellationClient: resty.New().
			SetBaseURL(appBaseURL).
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
