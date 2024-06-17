package templates

import (
	"fmt"
	"github/EmilioCliff/invoice-receipt/generator"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

func AzureEclipse(doc *generator.Document) error {
	if doc.Payment == nil {
		return fmt.Errorf("template requires payment methods")
	}
	doc.Pdf = fpdf.New("P", "mm", "A4", "") // 210 x 297 (mm)
	doc.Pdf.SetMargins(generator.MarginX, generator.MarginY, generator.MarginX)

	docType := strings.ToUpper(doc.Type)

	doc.Pdf.SetAutoPageBreak(true, 20)

	doc.Pdf.AddPage()

	doc.Pdf.SetFillColor(80, 80, 80)
	doc.Pdf.Rect(0, generator.MarginY, generator.MarginX-2, generator.ExtraLargeTextFontSize+3, "F")
	doc.Pdf.Rect(193, generator.MarginY, 15, generator.ExtraLargeTextFontSize+3, "F")
	doc.Pdf.SetFillColor(255, 255, 255)

	doc.Pdf.SetXY(generator.MarginX, generator.MarginY)
	doc.Pdf.SetFont("Arial", "B", generator.ExtraLargeTextFontSize)
	wd := doc.Pdf.GetStringWidth(docType)
	doc.Pdf.CellFormat(wd, generator.ExtraLargeTextFontSize, docType, "0", 0, "SM", false, 0, "")
	doc.Pdf.SetFontSize(generator.NormalTextFontSize)
	wd = doc.Pdf.GetStringWidth(doc.CompanyContact.CompanyName)
	doc.Pdf.SetXY(-(generator.MarginX + 5 + 5 + wd), generator.MarginY)
	doc.Pdf.CellFormat(wd, 16, doc.CompanyContact.CompanyName, "0", 0, "RB", false, 0, "")
	doc.Pdf.Ln(-1)
	doc.Pdf.Ln(1)
	doc.Pdf.SetX(-(generator.MarginX + 5 + 5 + wd))
	doc.Pdf.SetFontSize(generator.SmallTextFontSize)
	doc.Pdf.SetFontStyle("")
	doc.Pdf.CellFormat(wd, 16, "Your Business Partner", "0", 0, "RT", false, 0, "")

	doc.Pdf.Ln(32)
	lineY := doc.Pdf.GetY()
	doc.Pdf.Line(generator.MarginX, lineY, 200, lineY)
	doc.Pdf.Ln(10)

	totalY := doc.Pdf.GetY()
	doc.Pdf.SetFont("Arial", "B", generator.NormalTextFontSize)
	doc.Pdf.Cell(40, generator.CellLineHeight, fmt.Sprintf("%s TO :", docType))
	doc.Pdf.Ln(10)

	doc.Pdf.SetFontSize(generator.LargeTextFontSize)
	doc.Pdf.Cell(40, generator.CellLineHeight, doc.CustomerContact.Name)
	doc.Pdf.Ln(8)

	doc.Pdf.SetFont("Arial", "", generator.NormalTextFontSize)
	doc.Pdf.Cell(40, generator.CellLineHeight, fmt.Sprintf("P : %s", doc.CustomerContact.PhoneNumber))
	doc.Pdf.Ln(5)
	doc.Pdf.Cell(40, generator.CellLineHeight, fmt.Sprintf("E : %s", doc.CustomerContact.Email))
	doc.Pdf.Ln(5)
	doc.Pdf.Cell(40, generator.CellLineHeight, fmt.Sprintf("A : %s %s, %s", doc.CustomerContact.Address.PostalCode, doc.CustomerContact.Address.City, doc.CustomerContact.Address.Country))

	doc.Pdf.SetFont("Arial", "B", generator.NormalTextFontSize)
	wd = doc.Pdf.GetStringWidth("TOTAL PAID")
	doc.Pdf.SetXY(-(generator.MarginX + wd), totalY)
	if doc.Type == generator.Invoice {
		doc.Pdf.CellFormat(wd, generator.CellLineHeight, fmt.Sprintf("TOTAL DUE"), "0", 0, "RM", false, 0, "")
	} else {
		doc.Pdf.CellFormat(wd, generator.CellLineHeight, fmt.Sprintf("TOTAL PAID"), "0", 0, "RM", false, 0, "")
	}
	doc.Pdf.Ln(10)

	totalY = doc.Pdf.GetY()
	doc.Pdf.SetFontSize(generator.LargeTextFontSize)
	wd = doc.Pdf.GetStringWidth(" Ksh TOTAL PAID")

	// If tax is added find a way to calculate this
	// subtotal := 0.0
	// for _, item := range doc.Items {
	// 	totalPrice := item.UnitPrice * float64(item.Quantity)
	// 	subtotal += totalPrice
	// }
	totalPaidX, totalPaidY := doc.Pdf.GetXY()

	// doc.Pdf.SetFont("Arial", "B", generator.LargeTextFontSize)
	// doc.Pdf.SetXY(-(generator.MarginX + wd), totalY)
	// grandTotal := doc.CalculateTotalWithDiscount()
	// doc.Pdf.CellFormat(wd, generator.CellLineHeight, fmt.Sprintf("%s %.2f", doc.Options.CurrencySymbol, grandTotal), "0", 0, "RM", false, 0, "")
	doc.Pdf.Ln(8)

	// TODO: ADD A Line
	lineY = doc.Pdf.GetY()
	doc.Pdf.Line(190, lineY+5, 200, lineY+5)
	doc.Pdf.Ln(5)
	doc.Pdf.SetFont("Arial", "", generator.NormalTextFontSize)
	wd = doc.Pdf.GetStringWidth(fmt.Sprintf("No: %s", doc.DocumentData.DocumentNumber))
	doc.Pdf.SetX(-(generator.MarginX + wd))
	doc.Pdf.CellFormat(wd, generator.CellLineHeight, fmt.Sprintf("No: %s", doc.DocumentData.DocumentNumber), "0", 0, "RM", false, 0, "")
	doc.Pdf.Ln(5)
	wd = doc.Pdf.GetStringWidth(fmt.Sprintf("Date: %s", time.Now().Format("02/01/2006")))
	doc.Pdf.SetX(-(generator.MarginX + wd))
	doc.Pdf.CellFormat(wd, generator.CellLineHeight, fmt.Sprintf("Date: %s", time.Now().Format("02/01/2006")), "0", 0, "RM", false, 0, "")
	doc.Pdf.Ln(20)

	descriptionData := map[int]map[string]interface{}{
		0: {
			"fillHeader": []interface{}{true, 80, 80, 80},
			"fillRow":    []interface{}{true, 255, 255, 255},
			"border":     []string{"1", "1"},
			"note":       false,
			"payment":    true,
			"calculations": map[string]map[string][]string{
				"Subtotal": {
					"alignment": []string{"CM", "CM"},
					"margin":    []string{"1", "1"},
					"style":     []string{"B", ""},
					"fill":      []string{"255,255,255", "255,255,255"},
				},
			},
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

	if doc.DocumentData.Tax != 0 {
		descriptionData[0]["calculations"].(map[string]map[string][]string)["Tax"] = map[string][]string{
			"alignment": {"CM", "CM"},
			"margin":    {"1", "1"},
			"style":     {"B", ""},
			"fill":      {"255,255,255", "255,255,255"},
		}
	}

	if doc.DocumentData.Discount != 0 {
		descriptionData[0]["calculations"].(map[string]map[string][]string)["Discount"] = map[string][]string{
			"alignment": {"CM", "CM"},
			"margin":    {"1", "1"},
			"style":     {"B", ""},
			"fill":      {"255,255,255", "255,255,255"},
		}
	}

	descriptionData[0]["calculations"].(map[string]map[string][]string)["TOTAL"] = map[string][]string{
		"alignment": {"CM", "CM"},
		"margin":    {"1", "1"},
		"style":     {"B", "B"},
		"fill":      {"255,255,255", "80,80,80"},
	}

	grandTotal := 0.0
	if doc.DocumentData.Discount != 0 && doc.DocumentData.Tax != 0 {
		grandTotal = doc.CalculateTotalWithTaxAndDiscount()
	} else if doc.DocumentData.Discount != 0 {
		grandTotal = doc.CalculateTotalWithDiscount()
	} else if doc.DocumentData.Tax != 0 {
		grandTotal = doc.CalculateTotalWithTax()
	} else {
		grandTotal = doc.CalculateTotalWithoutTaxAndDiscount()
	}

	tableX, tableY := doc.Pdf.GetXY()
	doc.Pdf.SetXY(totalPaidX, totalPaidY)
	doc.Pdf.SetFont("Arial", "B", generator.LargeTextFontSize)
	doc.Pdf.SetXY(-(generator.MarginX + wd), totalY)
	doc.Pdf.CellFormat(wd, generator.CellLineHeight, fmt.Sprintf("%s %.2f", doc.Options.CurrencySymbol, grandTotal), "0", 0, "RM", false, 0, "")

	doc.Pdf.SetXY(tableX, tableY)
	doc.SetTableHeadings(descriptionData)
	err := doc.AddItemToTable(descriptionData)
	if err != nil {
		return err
	}

	doc.Pdf.SetY(-45)
	y := doc.Pdf.GetY()

	doc.Pdf.SetFont("Arial", "B", generator.LargeTextFontSize)
	doc.Pdf.Cell(100, generator.CellLineHeight, "Thank you for purchase!")

	doc.Pdf.SetFontSize(generator.NormalTextFontSize)
	doc.Pdf.SetX(130)
	doc.Pdf.CellFormat(70, generator.CellLineHeight, doc.DocumentData.IssuedBy, "0", 0, "CM", false, 0, "")
	doc.Pdf.Ln(10)
	y = doc.Pdf.GetY()

	doc.Pdf.Line(130, y, 200, y)
	doc.Pdf.Ln(3)
	doc.Pdf.SetX(130)
	doc.Pdf.SetFont("Arial", "", generator.SmallTextFontSize)
	doc.Pdf.CellFormat(70, generator.CellLineHeight, doc.DocumentData.IssuedByPosition, "0", 0, "CM", false, 0, "")

	return nil
}
