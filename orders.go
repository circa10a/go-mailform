package mailform

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	ordersEndpoint = "/orders"
)

var (
	// Service codes are the supported delivery services
	ServiceCodes = []string{
		"FEDEX_OVERNIGHT",
		"USPS_PRIORITY_EXPRESS",
		"USPS_PRIORITY",
		"USPS_CERTIFIED_PHYSICAL_RECEIPT",
		"USPS_CERTIFIED_RECEIPT",
		"USPS_CERTIFIED",
		"USPS_FIRST_CLASS",
		"USPS_STANDARD",
		"USPS_POSTCARD",
	}
)

// OrderInput is the input used to create an order.
// Representative of the example here: https://www.mailform.io/docs/api/#/orders
type OrderInput struct {
	// FilePath is the path of the file to print and mail
	// The PDF document to be mailed.
	// If this is not specified, the url parameter must be provided.
	// If both the file parameter and the url parameter are provided, the url parameter will be ignored
	FilePath string
	// The URL of the PDF document to be mailed: it will be downloaded completely before the API call completes.
	// The download must complete within 30 seconds.
	// If this is not specified, the file parameter must be provided.
	// If both the file parameter and the url parameter are provided, the url parameter will be ignored
	URL string
	// An optional customer reference to be attached to the order
	CustomerReference string
	// The delivery service to be used.
	// Must be one of FEDEX_OVERNIGHT, USPS_PRIORITY_EXPRESS, USPS_PRIORITY, USPS_CERTIFIED_PHYSICAL_RECEIPT, USPS_CERTIFIED_RECEIPT, USPS_CERTIFIED, USPS_FIRST_CLASS, USPS_STANDARD or USPS_POSTCARD.
	Service string
	// The webhook that should receive notifications about order updates to this order
	Webhook string
	// 	The company that this order should be associated with
	Company string
	// True if the document should be printed one page to a sheet, false if the document can be printed on both sides of a sheet
	Simplex bool
	// True if the document should be printed in color, false if the document should be printed in black and white
	Color bool
	// True if the document MUST be mailed in a flat envelope, false if it is acceptable to mail the document folded
	Flat bool
	// 	True if the document MUST use a real postage stamp, false if it is acceptable to mail the document using metered postage or an imprint
	Stamp bool
	// The message to be printed on the non-picture side of a postcard.
	Message string
	// The name of the recipient of this envelope or postcard
	ToName string
	// The organization or company associated with the recipient of this envelope or postcard
	ToOrganization string
	// The street number and name of the recipient of this envelope or postcard
	ToAddress1 string
	// The suite or room number of the recipient of this envelope or postcard
	ToAddress2 string
	// The address city of the recipient of this envelope or postcard
	ToCity string
	// The address state of the recipient of this envelope or postcard
	ToState string
	// The address postcode or zip code of the recipient of this envelope or postcard
	ToPostcode string
	// 	The address country of the recipient of this envelope or postcard
	ToCountry string
	// The name of the sender of this envelope or postcard
	FromName string
	// The organization or company associated with this address
	FromOrganization string
	// 	The street number and name of the sender of this envelope or postcard
	FromAddress1 string
	// The suite or room number of the sender of this envelope or postcard
	FromAddress2 string
	// The address city of the sender of this envelope or postcard
	FromCity string
	// The address state of the sender of this envelope or postcard
	FromState string
	// The address postcode or zip code of the sender of this envelope or postcard
	FromPostcode string
	// 	The address country of the sender of this envelope or postcard
	FromCountry string
	// The identifier of the bank account for the check associated with this order. Required if a check is to be included in this order.
	BankAccount string
	// The amount of the check associated with this order, in cents. Required if a check is to be included in this order.
	Amount int
	// The name of the recipient of the check associated with this order. Required if a check is to be included in this order.
	CheckName string
	// The number of the check associated with this order. Required if a check is to be included in this order.
	CheckNumber int
	// The memo line for the check associated with this order.
	CheckMemo string
}

// FormData converts order input fields to a map[string]string of form data.
func (o *OrderInput) FormData() map[string]string {
	formData := map[string]string{
		"customer_reference": o.CustomerReference,
		"service":            o.Service,
		"webhook":            o.Webhook,
		"company":            o.Company,
		"simplex":            strconv.FormatBool(o.Simplex),
		"color":              strconv.FormatBool(o.Color),
		"flat":               strconv.FormatBool(o.Flat),
		"stamp":              strconv.FormatBool(o.Stamp),
		"message":            o.Message,
		"to.name":            o.ToName,
		"to.organization":    o.ToOrganization,
		"to.address1":        o.ToAddress1,
		"to.address2":        o.ToAddress2,
		"to.city":            o.ToCity,
		"to.state":           o.ToState,
		"to.postcode":        o.ToPostcode,
		"to.country":         o.ToCountry,
		"from.name":          o.FromName,
		"from.organization":  o.FromOrganization,
		"from.address1":      o.FromAddress1,
		"from.address2":      o.FromAddress2,
		"from.city":          o.FromCity,
		"from.state":         o.FromState,
		"from.postcode":      o.FromPostcode,
		"from.country":       o.FromCountry,
	}

	// Check if URL is populated before sending
	// This can clash with file which is why we separate
	if o.URL != "" {
		formData["url"] = o.URL
	}

	// Check if bank/check related fields are populated before sending
	// These fields are typically expected to go together and are not required for a simple piece of mail
	if o.BankAccount != "" {
		formData["bank_account"] = o.BankAccount
	}
	if o.Amount != 0 {
		formData["amount"] = strconv.Itoa(o.Amount)
	}
	if o.CheckName != "" {
		formData["check_name"] = o.CheckName
	}
	if o.CheckNumber != 0 {
		formData["check_number"] = strconv.Itoa(o.CheckNumber)
	}
	if o.CheckMemo != "" {
		formData["check_memo"] = o.CheckMemo
	}

	return formData
}

// Order is the details of an order from mailform.
type Order struct {
	Success bool `json:"success"`
	Data    struct {
		Object    string    `json:"object"`
		ID        string    `json:"id"`
		Created   time.Time `json:"created"`
		Total     int       `json:"total"`
		Modified  time.Time `json:"modified"`
		Webhook   string    `json:"webhook"`
		Lineitems []struct {
			ID        string `json:"id"`
			Pagecount int    `json:"pagecount"`
			To        struct {
				Name         string `json:"name"`
				Address1     string `json:"address1"`
				Address2     string `json:"address2"`
				City         string `json:"city"`
				State        string `json:"state"`
				Postcode     string `json:"postcode"`
				Country      string `json:"country"`
				Formatted    string `json:"formatted"`
				Organization string `json:"organization"`
			} `json:"to"`
			From struct {
				Name         string `json:"name"`
				Address1     string `json:"address1"`
				Address2     string `json:"address2"`
				City         string `json:"city"`
				State        string `json:"state"`
				Postcode     string `json:"postcode"`
				Country      string `json:"country"`
				Formatted    string `json:"formatted"`
				Organization string `json:"organization"`
			} `json:"from"`
			Simplex bool   `json:"simplex"`
			Color   bool   `json:"color"`
			Service string `json:"service"`
			Pricing []struct {
				Type  string `json:"type"`
				Value int    `json:"value"`
			} `json:"pricing"`
		} `json:"lineitems"`
		Account            string    `json:"account"`
		CustomerReference  string    `json:"customer_reference"`
		Channel            string    `json:"channel"`
		TestMode           bool      `json:"test_mode"`
		State              string    `json:"state"`
		Cancelled          time.Time `json:"cancelled"`
		CancellationReason string    `json:"cancellation_reason"`
	} `json:"data"`
}

// CreateOrder creates a mailform order.
func (c *Client) CreateOrder(o OrderInput) (*Order, error) {
	order := &Order{}
	mailformErr := &ErrMailform{}

	// First validate order input
	err := o.Validate()
	if err != nil {
		return order, err
	}

	// Convert order input to form data
	formData := o.FormData()

	req := c.restClient.R()
	// If path is provided, set file form data and read local file
	if o.FilePath != "" {
		req.SetFile("file", o.FilePath)
	}

	// Send order
	resp, err := req.
		SetResult(order).
		SetError(mailformErr).
		SetFormData(formData).
		Post(ordersEndpoint)

	if err != nil {
		return order, err
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		return order, &ErrMailform{
			Err: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    strconv.Itoa(http.StatusUnauthorized),
				Message: "unauthorized",
			},
		}
	}

	if resp.IsError() {
		return order, mailformErr
	}

	// We can actually get 200 failed successfully and it's so dumb amirite
	err = checkBodyForErr(resp.Body())
	if err != nil {
		return order, err
	}

	return order, nil
}

// GetOrder gets a mailform order.
func (c *Client) GetOrder(o string) (*Order, error) {
	getOrderEndpoint := fmt.Sprintf("%s/%s", ordersEndpoint, o)
	order := &Order{}
	mailformErr := &ErrMailform{}

	resp, err := c.restClient.R().SetResult(order).SetError(mailformErr).Get(getOrderEndpoint)
	if err != nil {
		return order, err
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		return order, &ErrMailform{
			Err: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    strconv.Itoa(http.StatusUnauthorized),
				Message: "unauthorized",
			},
		}
	}

	if resp.IsError() {
		return order, mailformErr
	}

	// We can actually get 200 failed successfully and it's so dumb amirite
	err = checkBodyForErr(resp.Body())
	if err != nil {
		return order, err
	}

	return order, nil
}

// ErrOrderInvalid is returned when order input is invalid
type ErrOrderInvalid struct {
	message string
}

func (e *ErrOrderInvalid) Error() string {
	return e.message
}

// Validate validates an order input by checking all required fields.
// https://www.mailform.io/docs/api/#/orders
func (o *OrderInput) Validate() error {
	genericRejectionStr := "%s not provided, but is required"
	// Validate service code
	if supported := isServiceSupported(o.Service); !supported {
		return &ErrOrderInvalid{
			message: fmt.Sprintf("service code: '%s' not supported. Must be one of %v", o.Service, ServiceCodes),
		}
	}

	// Validate ToName
	if o.ToName == "" {
		return &ErrOrderInvalid{
			message: fmt.Sprintf(genericRejectionStr, "ToName"),
		}
	}

	// Validate ToAddress1
	if o.ToAddress1 == "" {
		return &ErrOrderInvalid{
			message: fmt.Sprintf(genericRejectionStr, "ToAddress1"),
		}
	}

	// Validate ToCity
	if o.ToCity == "" {
		return &ErrOrderInvalid{
			message: fmt.Sprintf(genericRejectionStr, "ToCity"),
		}
	}

	// Validate ToState
	if o.ToState == "" {
		return &ErrOrderInvalid{
			message: fmt.Sprintf(genericRejectionStr, "ToState"),
		}
	}

	// Validate ToPostcode
	if o.ToPostcode == "" {
		return &ErrOrderInvalid{
			message: fmt.Sprintf(genericRejectionStr, "ToPostcode"),
		}
	}

	// Validate ToCountry
	if o.ToCountry == "" {
		return &ErrOrderInvalid{
			message: fmt.Sprintf(genericRejectionStr, "ToCountry"),
		}
	}

	// Validate FromName
	if o.FromName == "" {
		return &ErrOrderInvalid{
			message: fmt.Sprintf(genericRejectionStr, "FromName"),
		}
	}

	// Validate FromAddress1
	if o.FromAddress1 == "" {
		return &ErrOrderInvalid{
			message: fmt.Sprintf(genericRejectionStr, "FromAddress1"),
		}
	}

	// Validate FromCity
	if o.FromCity == "" {
		return &ErrOrderInvalid{
			message: fmt.Sprintf(genericRejectionStr, "FromCity"),
		}
	}

	// Validate FromState
	if o.FromState == "" {
		return &ErrOrderInvalid{
			message: fmt.Sprintf(genericRejectionStr, "FromState"),
		}
	}

	// Validate FromPostcode
	if o.FromPostcode == "" {
		return &ErrOrderInvalid{
			message: fmt.Sprintf(genericRejectionStr, "FromPostcode"),
		}
	}

	// Validate FromCountry
	if o.FromCountry == "" {
		return &ErrOrderInvalid{
			message: fmt.Sprintf(genericRejectionStr, "FromCountry"),
		}
	}

	return nil
}

// isServiceSupported checks if service code string is supported or not.
func isServiceSupported(service string) bool {
	for _, v := range ServiceCodes {
		if service == v {
			return true
		}
	}

	return false
}
