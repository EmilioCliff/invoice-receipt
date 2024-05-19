package templates

import (
	"fmt"
	"github/EmilioCliff/invoice-receipt/generator"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

func CelestialDream(doc *generator.Document) error {
	if doc.Payment == nil {
		return fmt.Errorf("template requires payment methods")
	}
	doc.Pdf = fpdf.New("P", "mm", "A4", "") // 210 x 297 (mm)
	doc.Pdf.SetMargins(generator.MarginX, generator.MarginY, generator.MarginX)

	docType := strings.ToUpper(doc.Type)

	doc.Pdf.SetAutoPageBreak(true, 20)

	doc.Pdf.AddPage()

	targetWidth, targetHeight, err := generator.ResizeImage(doc.CompanyContact.CompanyLogo)
	if err != nil {
		return err
	}

	doc.Pdf.ImageOptions(doc.CompanyContact.CompanyLogo, generator.MarginX, generator.MarginY, targetWidth, targetHeight, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	doc.Pdf.SetFont("Times", "B", generator.LargeTextFontSize)
	wd := doc.Pdf.GetStringWidth(docType)
	doc.Pdf.SetXY(-(generator.MarginX + wd), generator.MarginY)

	doc.Pdf.CellFormat(wd, targetHeight, docType, "0", 0, "RM", false, 0, "")
	doc.Pdf.Ln(-1)
	doc.Pdf.Ln(10)

	y := doc.Pdf.GetY()
	doc.Pdf.SetFont("Arial", "B", generator.NormalTextFontSize)
	if doc.Type == generator.Invoice {
		doc.Pdf.Cell(40, 12, "BILLED TO:")
	} else {
		doc.Pdf.Cell(40, 12, "RECEIPT TO:")
	}
	doc.Pdf.Ln(-1)
	doc.Pdf.SetFontStyle("")

	_, lineHeight := doc.Pdf.GetFontSize()
	doc.Pdf.Cell(40, lineHeight, doc.CustomerContact.Name)

	if len(doc.CustomerContact.PhoneNumber) > 0 {
		doc.Pdf.Ln(5)
		doc.Pdf.Cell(40, lineHeight, doc.CustomerContact.PhoneNumber)
	}

	if len(doc.CustomerContact.Email) > 0 {
		doc.Pdf.Ln(5)
		doc.Pdf.Cell(40, lineHeight, doc.CustomerContact.Email)
	}
	tableY := doc.Pdf.GetY()

	data := fmt.Sprintf("%s No. %s", doc.Type, doc.DocumentData.DocumentNumber)
	wd = doc.Pdf.GetStringWidth(data)
	doc.Pdf.SetXY(-(generator.MarginX + wd), y)
	doc.Pdf.CellFormat(wd, generator.CellLineHeight, data, "0", 0, "RM", false, 0, "")
	doc.Pdf.Ln(5)
	doc.Pdf.SetX(-(generator.MarginX + wd))
	doc.Pdf.CellFormat(wd, generator.CellLineHeight, fmt.Sprintf("%s", time.Now().Format("02 January 2006")), "0", 0, "RM", false, 0, "")

	doc.Pdf.SetY(tableY)
	doc.Pdf.Ln(10)

	descriptionData := map[int]map[string]interface{}{
		0: {
			"fillHeader": []interface{}{true, 255, 255, 255},
			"fillRow":    []interface{}{true, 255, 255, 255},
			"border":     []string{"TB", "TB"},
			"calculations": map[string]map[string][]string{
				"Subtotal": {
					"alignment": []string{"CM", "CM"},
					"margin":    []string{"T", "T"},
					"style":     []string{"B", ""},
					"fill":      []string{"255,255,255", "255,255,255"},
				},
				"Tax": {
					"alignment": []string{"CM", "CM"},
					"margin":    []string{"B", "B"},
					"style":     []string{"B", ""},
					"fill":      []string{"255,255,255", "255,255,255"},
				},
				"Total": {
					"alignment": []string{"CM", "CM"},
					"margin":    []string{"T", "T"},
					"style":     []string{"B", "B"},
					"fill":      []string{"255,255,255", "255,255,255"},
				},
			},
			"note":    true,
			"payment": false,
		},
		1: {
			"columnName": doc.Options.TextItemsNameDescriptionTitle,
			"width":      90.0,
			"alignment":  []string{"LM", "LM"},
		},
		2: {
			"columnName": doc.Options.TextItemsQuantityTitle,
			"width":      30.0,
			"alignment":  []string{"CM", "CM"},
		},
		3: {
			"columnName": fmt.Sprintf("%s (%s)", doc.Options.TextItemsUnitCostTitle, doc.Options.CurrencySymbol),
			"width":      35.0,
			"alignment":  []string{"CM", "CM"},
		},
		4: {
			"columnName": fmt.Sprintf("%s (%s)", doc.Options.TextItemsTotalTitle, doc.Options.CurrencySymbol),
			"width":      35.0,
			"alignment":  []string{"CM", "CM"},
		},
	}

	doc.SetTableHeadings(descriptionData)
	err = doc.AddItemToTable(descriptionData)
	if err != nil {
		return err
	}

	doc.Pdf.Ln(20)
	doc.Pdf.SetFont("Times", "", generator.LargeTextFontSize)
	doc.Pdf.Cell(40, generator.CellLineHeight, "Thank You!")

	doc.Pdf.SetY(-50)
	y = doc.Pdf.GetY()

	doc.Pdf.SetFont("Times", "B", generator.NormalTextFontSize)
	doc.Pdf.Cell(40, generator.CellLineHeight, "PAYMENT INFORMATION")
	doc.Pdf.SetFontStyle("")
	doc.Pdf.Ln(5)
	if doc.Payment.Bank != nil {
		doc.Pdf.Cell(60, generator.CellLineHeight, doc.Payment.Bank.BankName)
		doc.Pdf.Ln(5)
		doc.Pdf.Cell(60, generator.CellLineHeight, fmt.Sprintf("Account Name: %s", doc.Payment.Bank.AccountName))
		doc.Pdf.Ln(5)
		doc.Pdf.Cell(60, generator.CellLineHeight, fmt.Sprintf("Account Number: %s", doc.Payment.Bank.AccountNumber))
	} else if doc.Payment.Paybill != nil {
		doc.Pdf.Cell(60, generator.CellLineHeight, fmt.Sprintf("Paybill Number: %s", doc.Payment.Paybill.PaybillNumber))
		doc.Pdf.Ln(5)
		doc.Pdf.Cell(60, generator.CellLineHeight, fmt.Sprintf("Account Number: %s", doc.Payment.Paybill.AccountNumber))
	} else {
		doc.Pdf.Cell(60, generator.CellLineHeight, fmt.Sprintf("Till Number: %s", doc.Payment.Till.TillNumber))
	}

	wd = doc.Pdf.GetStringWidth(doc.CompanyContact.CompanyName)
	doc.Pdf.SetXY(-(generator.MarginX + wd), y)
	doc.Pdf.CellFormat(wd, generator.CellLineHeight, doc.CompanyContact.CompanyName, "0", 0, "RM", false, 0, "")
	if len(doc.CompanyContact.CompanyEmail) > 0 {
		doc.Pdf.Ln(5)
		wd = doc.Pdf.GetStringWidth(doc.CompanyContact.CompanyEmail)
		doc.Pdf.SetX(-(generator.MarginX + wd))
		doc.Pdf.CellFormat(wd, generator.CellLineHeight, doc.CompanyContact.CompanyEmail, "0", 0, "RM", false, 0, "")
	}

	if len(doc.CompanyContact.CompanyAddress.City) > 0 {
		doc.Pdf.Ln(5)
		address := fmt.Sprintf("%s %s, %s", doc.CompanyContact.CompanyAddress.PostalCode, doc.CompanyContact.CompanyAddress.City, doc.CompanyContact.CompanyAddress.Country)
		wd = doc.Pdf.GetStringWidth(address)
		doc.Pdf.SetX(-(generator.MarginX + wd))
		doc.Pdf.CellFormat(wd, generator.CellLineHeight, address, "0", 0, "RM", false, 0, "")
	}

	if len(doc.CompanyContact.CompanyPhoneNumber) > 0 {
		doc.Pdf.Ln(5)
		doc.Pdf.SetFontStyle("I")
		wd = doc.Pdf.GetStringWidth(doc.CompanyContact.CompanyPhoneNumber)
		doc.Pdf.SetX(-(generator.MarginX + wd))
		doc.Pdf.CellFormat(wd, generator.CellLineHeight, doc.CompanyContact.CompanyPhoneNumber, "0", 0, "RM", false, 0, "")
	}
	return nil
}
