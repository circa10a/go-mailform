package mailform

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestFormData(t *testing.T) {
	tests := []struct {
		name       string
		orderInput *OrderInput
		expected   map[string]string
	}{
		{
			name: "EnsureStandardFieldsAreSetCorrectly",
			orderInput: &OrderInput{
				CustomerReference: "some_customer_reference",
				Service:           "some_service",
				Webhook:           "some_webhook",
				Company:           "some_company",
				Simplex:           true,
				Color:             true,
				Flat:              true,
				Stamp:             true,
				Message:           "some_message",
				ToName:            "some_to.name",
				ToOrganization:    "some_to.organization",
				ToAddress1:        "some_address1",
				ToAddress2:        "some_address2",
				ToCity:            "some_to.city",
				ToState:           "some_to.state",
				ToPostcode:        "some_to.postcode",
				ToCountry:         "some_to.country",
				FromName:          "some_from.name",
				FromOrganization:  "some_from.organization",
				FromAddress1:      "some_from.address1",
				FromAddress2:      "some_from.address2",
				FromCity:          "some_from.city",
				FromState:         "some_from.state",
				FromPostcode:      "some_from.postcode",
				FromCountry:       "some_from.country",
				URL:               "some_url",
				BankAccount:       "123456",
				Amount:            1,
				CheckName:         "some_checkname",
				CheckNumber:       123456,
				CheckMemo:         "some_memo",
			},
			expected: map[string]string{
				"customer_reference": "some_customer_reference",
				"service":            "some_service",
				"webhook":            "some_webhook",
				"company":            "some_company",
				"simplex":            strconv.FormatBool(true),
				"color":              strconv.FormatBool(true),
				"flat":               strconv.FormatBool(true),
				"stamp":              strconv.FormatBool(true),
				"message":            "some_message",
				"to.name":            "some_to.name",
				"to.organization":    "some_to.organization",
				"to.address1":        "some_address1",
				"to.address2":        "some_address2",
				"to.city":            "some_to.city",
				"to.state":           "some_to.state",
				"to.postcode":        "some_to.postcode",
				"to.country":         "some_to.country",
				"from.name":          "some_from.name",
				"from.organization":  "some_from.organization",
				"from.address1":      "some_from.address1",
				"from.address2":      "some_from.address2",
				"from.city":          "some_from.city",
				"from.state":         "some_from.state",
				"from.postcode":      "some_from.postcode",
				"from.country":       "some_from.country",
				"url":                "some_url",
				"bank_account":       "123456",
				"amount":             "1",
				"check_name":         "some_checkname",
				"check_number":       "123456",
				"check_memo":         "some_memo",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.orderInput.FormData()
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestCreateOrderError(t *testing.T) {
	fakeEndpoint := fmt.Sprintf("%s%s", DefaultBaseURL, ordersEndpoint)
	mailformClient, err := New(&Config{})
	assert.NoError(t, err)

	httpmock.ActivateNonDefault(mailformClient.restClient.GetClient())
	defer httpmock.DeactivateAndReset()

	// mock json rsponse
	response := `{"error":{"code":"erroroccurred","message":"unknown_error"},"detail":"Error: Not enough funds (2274:0)"}`
	// mock mailform.io response
	httpmock.RegisterResponder(http.MethodPost, fakeEndpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, response), nil

		})

	// check that addresses is nearby
	_, err = mailformClient.CreateOrder(OrderInput{
		Service:      "USPS_STANDARD",
		ToName:       "some_name",
		ToAddress1:   "some_address1",
		ToCity:       "some_city",
		ToState:      "some_state",
		ToPostcode:   "some_postcode",
		ToCountry:    "some_country",
		FromName:     "some_fromname",
		FromAddress1: "some_fromaddress1",
		FromCity:     "some_fromcity",
		FromState:    "some_fromstate",
		FromPostcode: "some_frompostcode",
		FromCountry:  "some_fromcountry",
	})

	// get count info
	httpmock.GetTotalCallCount()
	// get the amount of calls for the registered responder
	info := httpmock.GetCallCountInfo()
	// Check total calls
	assert.Equal(t, info[fmt.Sprintf("%s %s", http.MethodPost, fakeEndpoint)], 1)
	// Ensure error is returned
	mailErr := &ErrMailform{}
	assert.ErrorAs(t, err, &mailErr)
}

func TestCreateOrder(t *testing.T) {
	fakeEndpoint := fmt.Sprintf("%s%s", DefaultBaseURL, ordersEndpoint)
	fakeOrderID := "someID"
	mailformClient, err := New(&Config{})
	assert.NoError(t, err)

	httpmock.ActivateNonDefault(mailformClient.restClient.GetClient())
	defer httpmock.DeactivateAndReset()

	// mock json response
	response := &Order{
		Success: false,
		Data: struct {
			Object    string    "json:\"object\""
			ID        string    "json:\"id\""
			Created   time.Time "json:\"created\""
			Total     int       "json:\"total\""
			Modified  time.Time "json:\"modified\""
			Webhook   string    "json:\"webhook\""
			Lineitems []struct {
				ID        string "json:\"id\""
				Pagecount int    "json:\"pagecount\""
				To        struct {
					Name         string "json:\"name\""
					Address1     string "json:\"address1\""
					Address2     string "json:\"address2\""
					City         string "json:\"city\""
					State        string "json:\"state\""
					Postcode     string "json:\"postcode\""
					Country      string "json:\"country\""
					Formatted    string "json:\"formatted\""
					Organization string "json:\"organization\""
				} "json:\"to\""
				From struct {
					Name         string "json:\"name\""
					Address1     string "json:\"address1\""
					Address2     string "json:\"address2\""
					City         string "json:\"city\""
					State        string "json:\"state\""
					Postcode     string "json:\"postcode\""
					Country      string "json:\"country\""
					Formatted    string "json:\"formatted\""
					Organization string "json:\"organization\""
				} "json:\"from\""
				Simplex bool   "json:\"simplex\""
				Color   bool   "json:\"color\""
				Service string "json:\"service\""
				Pricing []struct {
					Type  string "json:\"type\""
					Value int    "json:\"value\""
				} "json:\"pricing\""
			} "json:\"lineitems\""
			Account            string    "json:\"account\""
			CustomerReference  string    "json:\"customer_reference\""
			Channel            string    "json:\"channel\""
			TestMode           bool      "json:\"test_mode\""
			State              string    "json:\"state\""
			Cancelled          time.Time "json:\"cancelled\""
			CancellationReason string    "json:\"cancellation_reason\""
		}{
			ID: fakeOrderID,
		},
	}

	// mock mailform.io response
	httpmock.RegisterResponder(http.MethodPost, fakeEndpoint,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, response)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})

	// check that addresses is nearby
	order, err := mailformClient.CreateOrder(OrderInput{
		Service:      "USPS_STANDARD",
		ToName:       "some_name",
		ToAddress1:   "some_address1",
		ToCity:       "some_city",
		ToState:      "some_state",
		ToPostcode:   "some_postcode",
		ToCountry:    "some_country",
		FromName:     "some_fromname",
		FromAddress1: "some_fromaddress1",
		FromCity:     "some_fromcity",
		FromState:    "some_fromstate",
		FromPostcode: "some_frompostcode",
		FromCountry:  "some_fromcountry",
	})
	assert.NoError(t, err)
	assert.Equal(t, order.Data.ID, fakeOrderID)

	// get count info
	httpmock.GetTotalCallCount()
	// get the amount of calls for the registered responder
	info := httpmock.GetCallCountInfo()
	// Check total calls
	assert.Equal(t, info[fmt.Sprintf("%s %s", http.MethodPost, fakeEndpoint)], 1)
}

func TestGetOrderError(t *testing.T) {
	fakeOrderID := "someID"
	fakeEndpoint := fmt.Sprintf("%s%s/%s", DefaultBaseURL, ordersEndpoint, fakeOrderID)
	mailformClient, err := New(&Config{})
	assert.NoError(t, err)

	httpmock.ActivateNonDefault(mailformClient.restClient.GetClient())
	defer httpmock.DeactivateAndReset()

	// mock json rsponse
	response := `{"error":{"code":"erroroccurred","message":"not found"}}`
	// mock mailform.io response
	httpmock.RegisterResponder(http.MethodGet, fakeEndpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, response), nil

		})

	// check that addresses is nearby
	_, err = mailformClient.GetOrder(fakeOrderID)

	// get count info
	httpmock.GetTotalCallCount()
	// get the amount of calls for the registered responder
	info := httpmock.GetCallCountInfo()
	// Check total calls
	assert.Equal(t, info[fmt.Sprintf("%s %s", http.MethodGet, fakeEndpoint)], 1)
	// Ensure error is returned
	mailErr := &ErrMailform{}
	assert.ErrorAs(t, err, &mailErr)
}

func TestGetOrder(t *testing.T) {
	fakeOrderID := "someID"
	fakeEndpoint := fmt.Sprintf("%s%s/%s", DefaultBaseURL, ordersEndpoint, fakeOrderID)

	mailformClient, err := New(&Config{})
	assert.NoError(t, err)

	httpmock.ActivateNonDefault(mailformClient.restClient.GetClient())
	defer httpmock.DeactivateAndReset()

	// mock json response
	response := &Order{
		Success: false,
		Data: struct {
			Object    string    "json:\"object\""
			ID        string    "json:\"id\""
			Created   time.Time "json:\"created\""
			Total     int       "json:\"total\""
			Modified  time.Time "json:\"modified\""
			Webhook   string    "json:\"webhook\""
			Lineitems []struct {
				ID        string "json:\"id\""
				Pagecount int    "json:\"pagecount\""
				To        struct {
					Name         string "json:\"name\""
					Address1     string "json:\"address1\""
					Address2     string "json:\"address2\""
					City         string "json:\"city\""
					State        string "json:\"state\""
					Postcode     string "json:\"postcode\""
					Country      string "json:\"country\""
					Formatted    string "json:\"formatted\""
					Organization string "json:\"organization\""
				} "json:\"to\""
				From struct {
					Name         string "json:\"name\""
					Address1     string "json:\"address1\""
					Address2     string "json:\"address2\""
					City         string "json:\"city\""
					State        string "json:\"state\""
					Postcode     string "json:\"postcode\""
					Country      string "json:\"country\""
					Formatted    string "json:\"formatted\""
					Organization string "json:\"organization\""
				} "json:\"from\""
				Simplex bool   "json:\"simplex\""
				Color   bool   "json:\"color\""
				Service string "json:\"service\""
				Pricing []struct {
					Type  string "json:\"type\""
					Value int    "json:\"value\""
				} "json:\"pricing\""
			} "json:\"lineitems\""
			Account            string    "json:\"account\""
			CustomerReference  string    "json:\"customer_reference\""
			Channel            string    "json:\"channel\""
			TestMode           bool      "json:\"test_mode\""
			State              string    "json:\"state\""
			Cancelled          time.Time "json:\"cancelled\""
			CancellationReason string    "json:\"cancellation_reason\""
		}{
			ID: fakeOrderID,
		},
	}

	// mock mailform.io response
	httpmock.RegisterResponder(http.MethodGet, fakeEndpoint,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, response)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})

	// check that addresses is nearby
	order, err := mailformClient.GetOrder(fakeOrderID)
	assert.NoError(t, err)
	assert.Equal(t, order.Data.ID, fakeOrderID)

	// get count info
	httpmock.GetTotalCallCount()
	// get the amount of calls for the registered responder
	info := httpmock.GetCallCountInfo()
	// Check total calls
	assert.Equal(t, info[fmt.Sprintf("%s %s", http.MethodGet, fakeEndpoint)], 1)
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name           string
		input          *OrderInput
		expectErr      bool
		expectedErrStr string
	}{
		{
			name: "EnsureServiceCodeErrorIsReturned",
			input: &OrderInput{
				Service: "Unsupported",
			},
			expectErr:      true,
			expectedErrStr: "service code:",
		},
		{
			name: "EnsureToNameErrorIsReturned",
			input: &OrderInput{
				Service: "USPS_STANDARD",
			},
			expectErr:      true,
			expectedErrStr: "ToName not provided",
		},
		{
			name: "EnsureToAddress1ErrorIsReturned",
			input: &OrderInput{
				Service: "USPS_STANDARD",
				ToName:  "some_name",
			},
			expectErr:      true,
			expectedErrStr: "ToAddress1 not provided",
		},
		{
			name: "EnsureToCityErrorIsReturned",
			input: &OrderInput{
				Service:    "USPS_STANDARD",
				ToName:     "some_name",
				ToAddress1: "some_address1",
			},
			expectErr:      true,
			expectedErrStr: "ToCity not provided",
		},
		{
			name: "EnsureToStateErrorIsReturned",
			input: &OrderInput{
				Service:    "USPS_STANDARD",
				ToName:     "some_name",
				ToAddress1: "some_address1",
				ToCity:     "some_city",
			},
			expectErr:      true,
			expectedErrStr: "ToState not provided",
		},
		{
			name: "EnsureToPostcodeErrorIsReturned",
			input: &OrderInput{
				Service:    "USPS_STANDARD",
				ToName:     "some_name",
				ToAddress1: "some_address1",
				ToCity:     "some_city",
				ToState:    "some_state",
			},
			expectErr:      true,
			expectedErrStr: "ToPostcode not provided",
		},
		{
			name: "EnsureToCountryErrorIsReturned",
			input: &OrderInput{
				Service:    "USPS_STANDARD",
				ToName:     "some_name",
				ToAddress1: "some_address1",
				ToCity:     "some_city",
				ToState:    "some_state",
				ToPostcode: "some_postcode",
			},
			expectErr:      true,
			expectedErrStr: "ToCountry not provided",
		},
		{
			name: "EnsureFromNameErrorIsReturned",
			input: &OrderInput{
				Service:    "USPS_STANDARD",
				ToName:     "some_name",
				ToAddress1: "some_address1",
				ToCity:     "some_city",
				ToState:    "some_state",
				ToPostcode: "some_postcode",
				ToCountry:  "some_country",
			},
			expectErr:      true,
			expectedErrStr: "FromName not provided",
		},
		{
			name: "EnsureFromAddress1ErrorIsReturned",
			input: &OrderInput{
				Service:    "USPS_STANDARD",
				ToName:     "some_name",
				ToAddress1: "some_address1",
				ToCity:     "some_city",
				ToState:    "some_state",
				ToPostcode: "some_postcode",
				ToCountry:  "some_country",
				FromName:   "some_fromname",
			},
			expectErr:      true,
			expectedErrStr: "FromAddress1 not provided",
		},
		{
			name: "EnsureFromCityErrorIsReturned",
			input: &OrderInput{
				Service:      "USPS_STANDARD",
				ToName:       "some_name",
				ToAddress1:   "some_address1",
				ToCity:       "some_city",
				ToState:      "some_state",
				ToPostcode:   "some_postcode",
				ToCountry:    "some_country",
				FromName:     "some_fromname",
				FromAddress1: "some_fromaddress1",
			},
			expectErr:      true,
			expectedErrStr: "FromCity not provided",
		},
		{
			name: "EnsureFromStateErrorIsReturned",
			input: &OrderInput{
				Service:      "USPS_STANDARD",
				ToName:       "some_name",
				ToAddress1:   "some_address1",
				ToCity:       "some_city",
				ToState:      "some_state",
				ToPostcode:   "some_postcode",
				ToCountry:    "some_country",
				FromName:     "some_fromname",
				FromAddress1: "some_fromaddress1",
				FromCity:     "some_fromcity",
			},
			expectErr:      true,
			expectedErrStr: "FromState not provided",
		},
		{
			name: "EnsureFromPostcodeErrorIsReturned",
			input: &OrderInput{
				Service:      "USPS_STANDARD",
				ToName:       "some_name",
				ToAddress1:   "some_address1",
				ToCity:       "some_city",
				ToState:      "some_state",
				ToPostcode:   "some_postcode",
				ToCountry:    "some_country",
				FromName:     "some_fromname",
				FromAddress1: "some_fromaddress1",
				FromCity:     "some_fromcity",
				FromState:    "some_fromstate",
			},
			expectErr:      true,
			expectedErrStr: "FromPostcode not provided",
		},
		{
			name: "EnsureFromCountryErrorIsReturned",
			input: &OrderInput{
				Service:      "USPS_STANDARD",
				ToName:       "some_name",
				ToAddress1:   "some_address1",
				ToCity:       "some_city",
				ToState:      "some_state",
				ToPostcode:   "some_postcode",
				ToCountry:    "some_country",
				FromName:     "some_fromname",
				FromAddress1: "some_fromaddress1",
				FromCity:     "some_fromcity",
				FromState:    "some_fromstate",
				FromPostcode: "some_frompostcode",
			},
			expectErr:      true,
			expectedErrStr: "FromCountry not provided",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.input.Validate()
			if test.expectErr {
				assert.ErrorContains(t, err, test.expectedErrStr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestIsServiceSupported(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "EnsureFedexOvernightSupported",
			input:    "FEDEX_OVERNIGHT",
			expected: true,
		},
		{
			name:     "EnsureUSPSPriorityExpressSupported",
			input:    "USPS_PRIORITY_EXPRESS",
			expected: true,
		},
		{
			name:     "EnsureUSPSPrioritySupported",
			input:    "USPS_PRIORITY",
			expected: true,
		},
		{
			name:     "EnsureUSPSCertifiedPhysicalReceiptSupported",
			input:    "USPS_CERTIFIED_PHYSICAL_RECEIPT",
			expected: true,
		},
		{
			name:     "EnsureUSPSCertifiedReceiptSupported",
			input:    "USPS_CERTIFIED_RECEIPT",
			expected: true,
		},
		{
			name:     "EnsureUSPSCertifiedSupported",
			input:    "USPS_CERTIFIED",
			expected: true,
		},
		{
			name:     "EnsureUSPSFirstClassSupported",
			input:    "USPS_FIRST_CLASS",
			expected: true,
		},
		{
			name:     "EnsureUSPSStandardSupported",
			input:    "USPS_STANDARD",
			expected: true,
		},
		{
			name:     "EnsureUSPSPostcardSupported",
			input:    "USPS_POSTCARD",
			expected: true,
		},
		{
			name:     "EnsureUnsupportedServiceIsFalse",
			input:    "FAKE",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := isServiceSupported(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}
