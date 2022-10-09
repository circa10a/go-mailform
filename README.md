# go-mailform

![GitHub tag (latest semver)](https://img.shields.io/github/v/tag/circa10a/go-mailform?style=plastic)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/circa10a/go-mailform)](https://pkg.go.dev/github.com/circa10a/go-mailform?tab=overview)
[![Go Report Card](https://goreportcard.com/badge/github.com/circa10a/go-mailform)](https://goreportcard.com/report/github.com/circa10a/go-mailform)

A small library to send physical mail from your Go applications using https://mailform.io

## Usage

First you will need an API Token from [mailform.io](https://mailform.io)

```bash
go get github.com/circa10a/go-mailform
```

```go
package main

import (
	"fmt"
	"os"

	"github.com/circa10a/go-mailform"
)

func main() {
	client, err := mailform.New(&mailform.Config{
		Token: "MAILFORM_API_TOKEN",
	})
	if err != nil {
		fmt.Println(err)
	}

	// Create order
	// You can send a PDF file via local filesystem path or a URL.
	// Must be one or the other, not both.
	order, err := client.CreateOrder(mailform.OrderInput{
		// Send local pdf
		FilePath: "./sample.pdf",
		// Or you can send the file via URL
		URL: "http://s3.amazonaws.com/some-bucket/sample.pdf",
		// Shipping service options:
		// FEDEX_OVERNIGHT USPS_PRIORITY_EXPRESS USPS_PRIORITY USPS_CERTIFIED_PHYSICAL_RECEIPT USPS_CERTIFIED_RECEIPT USPS_CERTIFIED USPS_FIRST_CLASS USPS_STANDARD USPS_POSTCARD
		Service:      "USPS_PRIORITY",
		ToName:       "A Name",
		ToAddress1:   "Address 1",
		ToCity:       "Seattle",
		ToState:      "WA",
		ToPostcode:   "00000",
		ToCountry:    "US",
		FromName:     "My Name",
		FromAddress1: "My Address 1",
		FromCity:     "Dallas",
		FromState:    "TX",
		FromPostcode: "00000",
		FromCountry:  "US",
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Order ID", order.Data.ID)

	// Get Order
	orderDetails, err := client.GetOrder(order.Data.ID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(orderDetails.Data)
}
```
