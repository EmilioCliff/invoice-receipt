package main

import (
	"github/EmilioCliff/invoice-receipt/generator"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	doc, err := New("UnNamed2", &generator.Options{
		DocumentType:   generator.Invoice,
		CurrencySymbol: "KES",
	})
	if err != nil {
		log.Fatal(err)
	}
	require.NoError(t, err)

	doc.SetHeaders("My header")
	doc.SetFooter("Your satisfaction is our best feeling")
	doc.SetDocumentData(&generator.DocumentData{
		DocumentNumber:   "20324334532342232",
		Discount:         10,
		DueDate:          15,
		TermNConditions:  "Payment is due within 15 days",
		Tax:              17,
		IssuedBy:         "Drew Feig",
		IssuedByPosition: "Administrator",
		Note:             "This is toscascaskcnsjkcnsdjkbd bjkdb sjb sdbsbvj bj jsd jhsd sd be payed before a certain date",
	})

	currentDir, err := os.Getwd()
	require.NoError(t, err)

	doc.SetCompanyAddress(&generator.CompanyContact{
		CompanyName:        "My Company Name",
		CompanyEmail:       "company@gmail.com",
		CompanyPhoneNumber: "07070707070",
		CompanyLogo:        filepath.Join(currentDir, "logo.png"),
		CompanyAddress: &generator.Address{
			PostalCode:    "00200",
			City:          "Nairobi",
			Country:       "Kenya",
			StreetAddress: "BuildingName, buildingRoom, Road",
		},
	})

	doc.SetCustomerAddress(&generator.CustomerContact{
		Name:        "Emilio Cliff",
		Email:       "emiliocliff@gmail.com",
		PhoneNumber: "0707070707",
		Address: &generator.Address{
			PostalCode:    "00100",
			City:          "Kitengela",
			Country:       "Kenya",
			StreetAddress: "BuildingName, buildingRoom, Road",
		},
	})

	doc.SetPaymentDetails(&generator.PaymentDetails{
		Bank: &generator.BankDetails{
			BankName:      "KCB",
			AccountName:   "Emilio Cliff Bank Of Kenya",
			AccountNumber: "324235364565675",
		},
		Paybill: &generator.PaybillDetails{
			PaybillNumber: "3435346545645",
			AccountNumber: "32434534645",
		},
		Till: &generator.BuyGoodsDetails{
			TillNumber: "35234534",
		},
	})

	for i := 0; i < 3; i++ {
		doc.AddItem(&generator.Item{
			Description: "Test Product 1",
			Quantity:    10,
			UnitPrice:   15.25,
		})
	}

	_, err = doc.Build()
	require.NoError(t, err)
}
