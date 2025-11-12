package mailform

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	customHTTPClient := &http.Client{Timeout: 99 * time.Second}

	tests := []struct {
		name               string
		input              *Config
		expectErr          bool
		expectedErr        error
		expectedAPIBaseURL string
		expectedAppBaseURL string
		expectedToken      string
		expectedHTTPClient *http.Client
	}{
		{
			name:               "EnsureNilConfigErrReturned",
			expectErr:          true,
			expectedErr:        ErrNilConfig,
			expectedHTTPClient: http.DefaultClient,
		},
		{
			name: "EnsureCustomAPIBaseURLIsSet",
			input: &Config{
				APIBaseURL: "customBaseURL",
			},
			expectedAPIBaseURL: "customBaseURL",
			expectedAppBaseURL: DefaultAppBaseURL,
			expectedHTTPClient: http.DefaultClient,
		},
		{
			name: "EnsureCustomAppBaseURLIsSet",
			input: &Config{
				AppBaseURL: "customAppBaseURL",
			},
			expectedAPIBaseURL: DefaultAPIBaseURL,
			expectedAppBaseURL: "customAppBaseURL",
			expectedHTTPClient: http.DefaultClient,
		},
		{
			name: "EnsureTokenURLIsSet",
			input: &Config{
				Token: "someToken",
			},
			expectedToken:      "someToken",
			expectedAPIBaseURL: DefaultAPIBaseURL,
			expectedAppBaseURL: DefaultAppBaseURL,
			expectedHTTPClient: http.DefaultClient,
		},
		{
			name: "EnsureCustomHTTPClientIsUsed",
			input: &Config{
				HTTPClient: customHTTPClient,
			},
			expectedAPIBaseURL: DefaultAPIBaseURL,
			expectedAppBaseURL: DefaultAppBaseURL,
			expectedHTTPClient: customHTTPClient,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := New(test.input)
			if test.expectErr {
				assert.ErrorIs(t, err, test.expectedErr)
				return
			}
			// Ensure token is passed correctly
			assert.Equal(t, test.expectedAPIBaseURL, actual.apiClient.BaseURL)
			// Ensure baseURL is passed correctly
			assert.Equal(t, test.expectedToken, actual.apiClient.Token)

			// Ensure token is passed correctly
			assert.Equal(t, test.expectedAppBaseURL, actual.appClient.BaseURL)
			// Ensure baseURL is passed correctly
			assert.Equal(t, test.expectedToken, actual.appClient.Token)

			// Ensure no unexpected error
			assert.NoError(t, err)
		})
	}
}

func TestCheckBodyForErr(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectErr   bool
		expectedErr error
	}{
		{
			name:        "EnsureUnmarshalErrorWitheEmptyBytes",
			expectErr:   true,
			expectedErr: &json.SyntaxError{},
		},
		{
			name:        "EnsureErrorIsReturnedNotDetailed",
			input:       []byte(`{"error":{"code":"erroroccurred","message":"no_file_uploaded"}}`),
			expectErr:   true,
			expectedErr: &ErrMailform{},
		},
		{
			name:        "EnsureErrorIsReturnedDetailed",
			input:       []byte(`{"error":{"code":"erroroccurred","message":"unknown_error"},"detail":"Error: Not enough funds (2274:0)"}`),
			expectErr:   true,
			expectedErr: &ErrMailform{},
		},
		{
			name:  "EnsureNoErrorIsReturned",
			input: []byte(`{"someKey":"someValue"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := checkBodyForErr(test.input)
			if test.expectErr {
				assert.ErrorAs(t, err, &test.expectedErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
