package templates

import (
	"fmt"
	"github/EmilioCliff/invoice-receipt/generator"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	"github.com/go-pdf/fpdf"
)

func MysticAura(doc *generator.Document) error {
	descriptionData := map[int]map[string]interface{}{
		0: {
			"fillHeader": []interface{}{true, 200, 200, 200},
			"fillRow":    []interface{}{true, 255, 255, 255},
			"border":     []string{"1", "1"},
			"note":       true,
			"payment":    false,
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

	doc.Pdf = fpdf.New("P", "mm", "A4", "") // 210 x 297 (mm)
	doc.Pdf.SetMargins(generator.MarginX, generator.MarginY, generator.MarginX)

	docType := strings.ToUpper(doc.Type)
	if doc.Footer != "" {
		doc.Pdf.SetFooterFunc(func() {
			doc.Pdf.SetY(-generator.MarginY)
			doc.Pdf.SetFont("Arial", "I", generator.ExtraSmallTextFontSize)
			wd := doc.Pdf.GetStringWidth(doc.Footer)
			doc.Pdf.Cell(wd, generator.CellLineHeight, doc.Footer)
			doc.Pdf.SetX(-generator.MarginX)
			doc.Pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", doc.Pdf.PageNo()),
				"", 0, "C", false, 0, "")
		})
	}

	doc.Pdf.SetAutoPageBreak(true, 20)

	doc.Pdf.AddPage()

	pageW, _ := doc.Pdf.GetPageSize()
	safeAreaW := pageW - 2*generator.MarginX

	targetWidth, targetHeight, err := generator.ResizeImage(doc.CompanyContact.CompanyLogo)
	if err != nil {
		return err
	}

	doc.Pdf.ImageOptions(doc.CompanyContact.CompanyLogo, safeAreaW/2+15, generator.MarginY, targetWidth, targetHeight, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	doc.Pdf.SetXY(generator.MarginX, generator.MarginY)
	doc.Pdf.SetFont("Arial", "B", generator.LargeTextFontSize)
	doc.Pdf.Cell(40, 20, docType)
	doc.Pdf.Ln(18)

	doc.CompanyContact.LayerCompanyContact(doc)
	doc.Pdf.Ln(10)

	headerY := doc.Pdf.GetY()

	doc.Pdf.SetX(generator.MarginX)
	doc.Pdf.SetFont("Arial", "B", generator.LargeTextFontSize)
	doc.Pdf.Cell(40, 20, docType+" TO")
	doc.Pdf.Ln(18)

	doc.CustomerContact.LayerCustomerContact(doc)
	tableY := doc.Pdf.GetY()

	middleX := (safeAreaW / 2) + 15
	doc.Pdf.SetXY(middleX, headerY)
	doc.Pdf.SetFont("Arial", "B", generator.LargeTextFontSize)
	doc.Pdf.Cell(40, 20, docType)
	doc.Pdf.Ln(18)
	doc.Pdf.SetX(middleX)

	documentDataStructure := map[string]interface{}{
		"names": []string{fmt.Sprintf("%v NO: ", docType), fmt.Sprintf("%v DATE: ", docType)},
		"formats": map[string]string{
			"date_format": "2006-01-02",
			"alignment":   "",
			"font_style":  "B",
		},
	}

	doc.DocumentData.LayerDocumentData(docType, doc, documentDataStructure)

	doc.Pdf.SetXY(generator.MarginX, tableY+10)
	doc.SetTableHeadings(descriptionData)
	doc.AddItemToTable(descriptionData)
	return nil
}
