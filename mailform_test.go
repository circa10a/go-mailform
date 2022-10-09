package mailform

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name            string
		input           *Config
		expectErr       bool
		expectedErr     error
		expectedBaseURL string
		expectedToken   string
	}{
		{
			name:        "EnsureNilConfigErrReturned",
			expectErr:   true,
			expectedErr: ErrNilConfig,
		},
		{
			name: "EnsureCustomBaseURLIsSet",
			input: &Config{
				BaseURL: "customBaseURL",
			},
			expectedBaseURL: "customBaseURL",
		},
		{
			name: "EnsureTokenURLIsSet",
			input: &Config{
				Token: "someToken",
			},
			expectedToken:   "someToken",
			expectedBaseURL: DefaultBaseURL,
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
			assert.Equal(t, test.expectedBaseURL, actual.restClient.BaseURL)
			// Ensure baseURL is passed correctly
			assert.Equal(t, test.expectedToken, actual.restClient.Token)
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
