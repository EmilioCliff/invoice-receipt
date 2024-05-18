package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	doc, err := New(2, &Options{
		DocumentType: Receipt,
	})
	if err != nil {
		log.Fatal(err)
	}
	require.NoError(t, err)

	doc.SetHeaders("My header")
	doc.SetFooter("Your satisfaction is our best feeling")
	doc.SetDocumentData(&DocumentData{
		DocumentNumber: "20324334532342232",
		// Note:           "This is toscascaskcnsjkcnsdjkbd bjkdb sjb sdbsbvj bj jsd jhsd sd be payed before a certain date",
	})

	currentDir, err := os.Getwd()
	require.NoError(t, err)

	doc.SetCompanyAddress(&CompanyContact{
		CompanyName:        "My Company Name",
		CompanyEmail:       "company@gmail.com",
		CompanyPhoneNumber: "07070707070",
		CompanyLogo:        filepath.Join(currentDir, "logo.png"),
		CompanyAddress: &Address{
			PostalCode: "00200",
			City:       "Nairobi",
			Country:    "Kenya",
		},
	})

	doc.SetCustomerAddress(&CustomerContact{
		Name:        "Emilio Cliff",
		Email:       "emiliocliff@gmail.com",
		PhoneNumber: "0707070707",
		Address: &Address{
			PostalCode:    "00100",
			City:          "Kitengela",
			Country:       "Kenya",
			StreetAddress: "BuildingName, buildingRoom, Road",
		},
	})

	doc.SetPaymentDetails(&PaymentDetails{
		Bank: &BankDetails{
			BankName:      "KCB",
			AccountName:   "Emilio Cliff Bank",
			AccountNumber: "324235364565675",
		},
		Paybill: &PaybillDetails{
			PaybillNumber: "3435346545645",
			AccountNumber: "32434534645",
		},
		Till: &BuyGoodsDetails{
			TillNumber: "35234534",
		},
	})

	for i := 0; i < 3; i++ {
		doc.AddItem(&Item{
			Description: "Test Product 1",
			Quantity:    10,
			UnitPrice:   15.25,
		})
	}

	_, err = doc.Build()
	require.NoError(t, err)
}
