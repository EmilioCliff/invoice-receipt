package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	doc, err := New(1, &Options{
		DocumentType: Invoice,
	})
	if err != nil {
		log.Fatal(err)
	}
	require.NoError(t, err)

	doc.SetHeaders("My header")
	doc.SetFooter("My footer")
	doc.SetDocumentData(&DocumentData{
		DocumentNumber: "2020",
		Date:           "2023-10-09",
		Note:           "My Note",
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

	for i := 0; i < 30; i++ {
		doc.AddItem(&Item{
			Description: "Test Product 1",
			Quantity:    10,
			UnitPrice:   15.25,
		})
	}

	_, err = doc.Build()
	require.NoError(t, err)
}
