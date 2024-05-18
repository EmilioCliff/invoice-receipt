package templates

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"strings"
	"time"

	gen "github.com/EmilioCliff/main"
	"github.com/go-pdf/fpdf"
)

// header := [NumberOfColums]string{doc.Options.TextItemsNumberTitle,
//
//		doc.Options.TextItemsNameDescriptionTitle,
//		doc.Options.TextItemsQuantityTitle,
//		fmt.Sprintf("%s (%s)", doc.Options.TextItemsUnitCostTitle, doc.Options.CurrencySymbol),
//		fmt.Sprintf("%s (%s)", doc.Options.TextItemsTotalTitle, doc.Options.CurrencySymbol),
//	}
//
// colWidth := [NumberOfColums]float64{NumberColumnOffset, DescriptionColumnOffset, QuantityColumnOffset, UnitPriceColumnOffset, TotalPriceOffset}
func MysticAura(doc *gen.Document) error {
	descriptionData := map[int]map[string]interface{}{
		0: {
			"fillHeader": []interface{}{true, 200, 200, 200},
			"fillRow":    []interface{}{true, 255, 255, 255},
			"border":     []string{"1", "1"},
		},
		1: {
			"columnName": doc.Options.TextItemsNumberTitle,
			"width":      10.0,
			"alignment":  []string{"CM", "CM"},
		},
		2: {
			"columnName": doc.Options.TextItemsNameDescriptionTitle,
			"width":      75.0,
			"alignment":  []string{"CM", "LM"},
		},
		3: {
			"columnName": doc.Options.TextItemsQuantityTitle,
			"width":      25.0,
			"alignment":  []string{"CM", "CM"},
		},
		4: {
			"columnName": fmt.Sprintf("%s (%s)", doc.Options.TextItemsUnitCostTitle, doc.Options.CurrencySymbol),
			"width":      40.0,
			"alignment":  []string{"CM", "CM"},
		},
		5: {
			"columnName": fmt.Sprintf("%s (%s)", doc.Options.TextItemsTotalTitle, doc.Options.CurrencySymbol),
			"width":      40.0,
			"alignment":  []string{"CM", "CM"},
		},
	}

	doc.pdf = fpdf.New("P", "mm", "A4", "") // 210 x 297 (mm)
	doc.pdf.SetMargins(gen.MarginX, gen.MarginY, gen.MarginX)

	docType := strings.ToUpper(doc.Type)
	if doc.Footer != "" {
		doc.pdf.SetFooterFunc(func() {
			doc.pdf.SetY(-gen.MarginY)
			doc.pdf.SetFont("Arial", "I", gen.ExtraSmallTextFontSize)
			wd := doc.pdf.GetStringWidth(doc.Footer)
			doc.pdf.Cell(wd, gen.CellLineHeight, doc.Footer)
			doc.pdf.SetX(-gen.MarginX)
			doc.pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", doc.pdf.PageNo()),
				"", 0, "C", false, 0, "")
		})
	}

	doc.pdf.SetAutoPageBreak(true, 20)

	doc.pdf.AddPage()

	pageW, _ := doc.pdf.GetPageSize()
	safeAreaW := pageW - 2*gen.MarginX

	targetWidth, targetHeight, err := gen.ResizeImage(doc.CompanyContact.CompanyLogo)
	if err != nil {
		return err
	}

	doc.pdf.ImageOptions(doc.CompanyContact.CompanyLogo, safeAreaW/2+15, gen.MarginY, targetWidth, targetHeight, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	doc.pdf.SetXY(gen.MarginX, gen.MarginY)
	doc.pdf.SetFont("Arial", "B", gen.Extragen.LargeTextFontSize)
	doc.pdf.Cell(40, 20, docType)
	doc.pdf.Ln(18)

	doc.CompanyContact.LayerCompanyContact(doc)
	doc.pdf.Ln(10)

	headerY := doc.pdf.GetY()

	doc.pdf.SetX(gen.MarginX)
	doc.pdf.SetFont("Arial", "B", gen.LargeTextFontSize)
	doc.pdf.Cell(40, 20, docType+" TO")
	doc.pdf.Ln(18)

	doc.CustomerContact.LayerCustomerContact(doc)
	tableY := doc.pdf.GetY()

	middleX := (safeAreaW / 2) + 15
	doc.pdf.SetXY(middleX, headerY)
	doc.pdf.SetFont("Arial", "B", gen.LargeTextFontSize)
	doc.pdf.Cell(40, 20, docType)
	doc.pdf.Ln(18)
	doc.pdf.SetX(middleX)

	documentDataStructure := map[string]interface{}{
		"names": []string{fmt.Sprintf("%v NO: ", docType), fmt.Sprintf("%v DATE: ", docType)},
		"formats": map[string]string{
			"date_format": "2006-01-02",
			"alignment":   "",
			"font_style":  "B",
		},
	}

	doc.DocumentData.LayerDocumentData(docType, doc, documentDataStructure)

	doc.pdf.SetXY(gen.MarginX, tableY+10)
	doc.SetTableHeadings(descriptionData)
	doc.AddItemToTable(descriptionData)
	return nil
}

func CelestialDream(doc *gen.Document) error {
	if doc.Payment == nil {
		return fmt.Errorf("template requires payment methods")
	}
	doc.pdf = fpdf.New("P", "mm", "A4", "") // 210 x 297 (mm)
	doc.pdf.SetMargins(gen.MarginX, gen.MarginY, gen.MarginX)

	docType := strings.ToUpper(doc.Type)

	doc.pdf.SetAutoPageBreak(true, 20)

	doc.pdf.AddPage()

	targetWidth, targetHeight, err := gen.ResizeImage(doc.CompanyContact.CompanyLogo)
	if err != nil {
		return err
	}

	doc.pdf.ImageOptions(doc.CompanyContact.CompanyLogo, gen.MarginX, gen.MarginY, targetWidth, targetHeight, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	doc.pdf.SetFont("Times", "B", gen.Extragen.LargeTextFontSize)
	wd := doc.pdf.GetStringWidth(docType)
	doc.pdf.SetXY(-(gen.MarginX + wd), gen.MarginY)

	doc.pdf.CellFormat(wd, targetHeight, docType, "0", 0, "RM", false, 0, "")
	doc.pdf.Ln(-1)
	doc.pdf.Ln(10)

	y := doc.pdf.GetY()
	doc.pdf.SetFont("Arial", "B", gen.NormalTextFontSize)
	if doc.Type == gen.Invoice {
		doc.pdf.Cell(40, 12, "BILLED TO:")
	} else {
		doc.pdf.Cell(40, 12, "RECEIPT TO:")
	}
	doc.pdf.Ln(-1)
	doc.pdf.SetFontStyle("")

	_, lineHeight := doc.pdf.GetFontSize()
	doc.pdf.Cell(40, lineHeight, doc.CustomerContact.Name)

	if len(doc.CustomerContact.PhoneNumber) > 0 {
		doc.pdf.Ln(5)
		doc.pdf.Cell(40, lineHeight, doc.CustomerContact.PhoneNumber)
	}

	if len(doc.CustomerContact.Email) > 0 {
		doc.pdf.Ln(5)
		doc.pdf.Cell(40, lineHeight, doc.CustomerContact.Email)
	}
	tableY := doc.pdf.GetY()

	data := fmt.Sprintf("%s No. %s", doc.Type, doc.DocumentData.DocumentNumber)
	wd = doc.pdf.GetStringWidth(data)
	doc.pdf.SetXY(-(gen.MarginX + wd), y)
	doc.pdf.CellFormat(wd, gen.CellLineHeight, data, "0", 0, "RM", false, 0, "")
	doc.pdf.Ln(5)
	doc.pdf.SetX(-(gen.MarginX + wd))
	doc.pdf.CellFormat(wd, gen.CellLineHeight, fmt.Sprintf("%s", time.Now().Format("02 January 2006")), "0", 0, "RM", false, 0, "")

	doc.pdf.SetY(tableY)
	doc.pdf.Ln(10)

	descriptionData := map[int]map[string]interface{}{
		0: {
			"fillHeader": []interface{}{true, 255, 255, 255},
			"fillRow":    []interface{}{true, 255, 255, 255},
			"border":     []string{"TB", "TB"},
		},
		1: {
			"columnName": doc.Options.TextItemsNameDescriptionTitle,
			"width":      100.0,
			"alignment":  []string{"LM", "LM"},
		},
		2: {
			"columnName": doc.Options.TextItemsQuantityTitle,
			"width":      30.0,
			"alignment":  []string{"CM", "CM"},
		},
		3: {
			"columnName": fmt.Sprintf("%s (%s)", doc.Options.TextItemsUnitCostTitle, doc.Options.CurrencySymbol),
			"width":      30.0,
			"alignment":  []string{"CM", "CM"},
		},
		4: {
			"columnName": fmt.Sprintf("%s (%s)", doc.Options.TextItemsTotalTitle, doc.Options.CurrencySymbol),
			"width":      30.0,
			"alignment":  []string{"CM", "CM"},
		},
	}

	doc.SetTableHeadings(descriptionData)
	doc.AddItemToTable(descriptionData)

	doc.pdf.Ln(20)
	doc.pdf.SetFont("Times", "", gen.LargeTextFontSize)
	doc.pdf.Cell(40, gen.CellLineHeight, "Thank You!")

	doc.pdf.SetY(-50)
	y = doc.pdf.GetY()

	doc.pdf.SetFont("Times", "B", gen.NormalTextFontSize)
	doc.pdf.Cell(40, gen.CellLineHeight, "PAYMENT INFORMATION")
	doc.pdf.SetFontStyle("")
	doc.pdf.Ln(5)
	if doc.Payment.Bank != nil {
		doc.pdf.Cell(60, gen.CellLineHeight, doc.Payment.Bank.BankName)
		doc.pdf.Ln(5)
		doc.pdf.Cell(60, gen.CellLineHeight, fmt.Sprintf("Account Name: %s", doc.Payment.Bank.AccountName))
		doc.pdf.Ln(5)
		doc.pdf.Cell(60, gen.CellLineHeight, fmt.Sprintf("Account Number: %s", doc.Payment.Bank.AccountNumber))
	} else if doc.Payment.Paybill != nil {
		doc.pdf.Cell(60, gen.CellLineHeight, fmt.Sprintf("Paybill Number: %s", doc.Payment.Paybill.PaybillNumber))
		doc.pdf.Ln(5)
		doc.pdf.Cell(60, gen.CellLineHeight, fmt.Sprintf("Account Number: %s", doc.Payment.Paybill.AccountNumber))
	} else {
		doc.pdf.Cell(60, gen.CellLineHeight, fmt.Sprintf("Till Number: %s", doc.Payment.Till.TillNumber))
	}

	wd = doc.pdf.GetStringWidth(doc.CompanyContact.CompanyName)
	doc.pdf.SetXY(-(gen.MarginX + wd), y)
	doc.pdf.CellFormat(wd, gen.CellLineHeight, doc.CompanyContact.CompanyName, "0", 0, "RM", false, 0, "")
	if len(doc.CompanyContact.CompanyEmail) > 0 {
		doc.pdf.Ln(5)
		wd = doc.pdf.GetStringWidth(doc.CompanyContact.CompanyEmail)
		doc.pdf.SetX(-(gen.MarginX + wd))
		doc.pdf.CellFormat(wd, gen.CellLineHeight, doc.CompanyContact.CompanyEmail, "0", 0, "RM", false, 0, "")
	}

	if len(doc.CompanyContact.CompanyAddress.City) > 0 {
		doc.pdf.Ln(5)
		address := fmt.Sprintf("%s %s, %s", doc.CompanyContact.CompanyAddress.PostalCode, doc.CompanyContact.CompanyAddress.City, doc.CompanyContact.CompanyAddress.Country)
		wd = doc.pdf.GetStringWidth(address)
		doc.pdf.SetX(-(gen.MarginX + wd))
		doc.pdf.CellFormat(wd, gen.CellLineHeight, address, "0", 0, "RM", false, 0, "")
	}

	if len(doc.CompanyContact.CompanyPhoneNumber) > 0 {
		doc.pdf.Ln(5)
		doc.pdf.SetFontStyle("I")
		wd = doc.pdf.GetStringWidth(doc.CompanyContact.CompanyPhoneNumber)
		doc.pdf.SetX(-(gen.MarginX + wd))
		doc.pdf.CellFormat(wd, gen.CellLineHeight, doc.CompanyContact.CompanyPhoneNumber, "0", 0, "RM", false, 0, "")
	}
	return nil
}
