package templates

import (
	"fmt"
	"github/EmilioCliff/invoice-receipt/generator"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

func CrimsonWhisper(doc *generator.Document) error {
	if doc.Payment == nil {
		return fmt.Errorf("template requires payment methods")
	}
	doc.Pdf = fpdf.New("P", "mm", "A4", "") // 210 x 297 (mm)
	doc.Pdf.SetMargins(generator.MarginX, generator.MarginY, generator.MarginX)

	docType := strings.ToUpper(doc.Type)

	doc.Pdf.SetAutoPageBreak(true, 20)

	doc.Pdf.AddPage()

	doc.Pdf.SetHomeXY()
	doc.Pdf.SetFont("Arial", "B", generator.LargeTextFontSize)
	doc.Pdf.Cell(50, generator.CellLineHeight, doc.CompanyContact.CompanyName)
	wd := doc.Pdf.GetStringWidth(docType)
	doc.Pdf.SetX(-(generator.MarginX + wd))
	doc.Pdf.CellFormat(wd, generator.CellLineHeight, docType, "0", 0, "RM", false, 0, "")
	doc.Pdf.SetFont("Arial", "", generator.SmallTextFontSize)
	doc.Pdf.Ln(10)
	if doc.CompanyContact.CompanyAddress.StreetAddress != "" {
		doc.Pdf.Cell(50, generator.CellLineHeight, doc.CompanyContact.CompanyAddress.StreetAddress)
		doc.Pdf.Ln(5)
	}
	doc.Pdf.Cell(50, generator.CellLineHeight, fmt.Sprintf("%s %s, %s", doc.CompanyContact.CompanyAddress.PostalCode, doc.CompanyContact.CompanyAddress.City, doc.CompanyContact.CompanyAddress.Country))

	doc.Pdf.Ln(25)

	y := doc.Pdf.GetY()
	doc.Pdf.SetFont("Arial", "B", generator.NormalTextFontSize)
	doc.Pdf.Cell(55, generator.CellLineHeight, "Bill To")
	doc.Pdf.Ln(8)
	doc.Pdf.SetFontStyle("")
	doc.Pdf.Cell(55, generator.CellLineHeight, doc.CustomerContact.Name)
	doc.Pdf.Ln(5)
	doc.Pdf.Cell(55, generator.CellLineHeight, fmt.Sprintf("%s %s", doc.CustomerContact.Address.PostalCode, doc.CustomerContact.Address.City))
	doc.Pdf.Ln(5)
	doc.Pdf.Cell(55, generator.CellLineHeight, doc.CustomerContact.Address.Country)

	doc.Pdf.SetXY(generator.MarginX+50, y)
	doc.Pdf.SetFontStyle("B")
	doc.Pdf.Cell(55, generator.CellLineHeight, "Ship To")
	doc.Pdf.Ln(8)
	doc.Pdf.SetX(generator.MarginX + 50)
	doc.Pdf.SetFontStyle("")
	doc.Pdf.Cell(55, generator.CellLineHeight, doc.CustomerContact.Name)
	doc.Pdf.Ln(5)
	doc.Pdf.SetX(generator.MarginX + 50)
	addres := strings.Split(doc.CustomerContact.Address.StreetAddress, ",")
	// doc.Pdf.Cell(55, generator.CellLineHeight, fmt.Sprintf("%s %s", doc.CustomerContact.Address.PostalCode, doc.CustomerContact.Address.City))
	doc.Pdf.Cell(55, generator.CellLineHeight, fmt.Sprintf("%s %s", addres[0], addres[1]))
	doc.Pdf.Ln(5)
	doc.Pdf.SetX(generator.MarginX + 50)
	// doc.Pdf.Cell(55, generator.CellLineHeight, doc.CustomerContact.Address.Country)
	doc.Pdf.Cell(55, generator.CellLineHeight, addres[2])

	doc.Pdf.SetXY(generator.MarginX+110, y)
	doc.Pdf.SetFontStyle("B")
	doc.Pdf.Cell(30, generator.CellLineHeight, fmt.Sprintf("%s#", doc.Type))
	doc.Pdf.SetFontStyle("")
	doc.Pdf.CellFormat(50, generator.CellLineHeight, doc.DocumentData.DocumentNumber, "0", 0, "RM", false, 0, "")
	doc.Pdf.Ln(8)

	doc.Pdf.SetX(generator.MarginX + 110)
	doc.Pdf.SetFontStyle("B")
	doc.Pdf.Cell(30, generator.CellLineHeight, fmt.Sprintf("%s Date", doc.Type))
	doc.Pdf.SetFontStyle("")
	doc.Pdf.CellFormat(50, generator.CellLineHeight, time.Now().Format("02/01/2006"), "0", 0, "RM", false, 0, "")
	doc.Pdf.Ln(8)

	doc.Pdf.SetX(generator.MarginX + 110)
	doc.Pdf.SetFontStyle("B")
	doc.Pdf.Cell(30, generator.CellLineHeight, "P.O#")
	doc.Pdf.SetFontStyle("")
	doc.Pdf.CellFormat(50, generator.CellLineHeight, "03/05/2023", "0", 0, "RM", false, 0, "") // Add P.O date
	doc.Pdf.Ln(8)

	doc.Pdf.SetX(generator.MarginX + 110)
	doc.Pdf.SetFontStyle("B")
	doc.Pdf.Cell(30, generator.CellLineHeight, "Due Date")
	doc.Pdf.SetFontStyle("")
	doc.Pdf.CellFormat(50, generator.CellLineHeight, "03/05/2023", "0", 0, "RM", false, 0, "") // Add Due Date Day

	doc.Pdf.Ln(20)

	descriptionData := map[int]map[string]interface{}{
		0: {
			"fillHeader": []interface{}{true, 100, 100, 100},
			"fillRow":    []interface{}{true, 255, 255, 255},
			"border":     []string{"1", "1"},
			"calculations": map[string]map[string][]string{
				"Subtotal": {
					"alignment": []string{"RM", "RM"},
					"margin":    []string{"T", "LRT"},
					"style":     []string{"B", ""},
					"fill":      []string{"255,255,255", "255,255,255"},
				},
				"Tax": {
					"alignment": []string{"RM", "RM"},
					"margin":    []string{"0", "LR"},
					"style":     []string{"B", ""},
					"fill":      []string{"255,255,255", "255,255,255"},
				},
				"TOTAL": {
					"alignment": []string{"RM", "RM"},
					"margin":    []string{"0", "1"},
					"style":     []string{"B", "B"},
					"fill":      []string{"255,255,255", "100,100,100"},
				},
			},
			"note":    false,
			"payment": false,
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
			"alignment":  []string{"CM", "RM"},
		},
		5: {
			"columnName": fmt.Sprintf("%s (%s)", doc.Options.TextItemsTotalTitle, doc.Options.CurrencySymbol),
			"width":      40.0,
			"alignment":  []string{"CM", "RM"},
		},
	}

	doc.SetTableHeadings(descriptionData)
	err := doc.AddItemToTable(descriptionData)
	if err != nil {
		return err
	}

	return nil
}