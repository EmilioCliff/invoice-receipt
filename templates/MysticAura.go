package templates

import (
	"fmt"
	"github/EmilioCliff/invoice-receipt/generator"
	_ "image/jpeg"
	_ "image/png"
	"strings"
	"time"

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
			"calculations": map[string]map[string][]string{
				"Subtotal": {
					"alignment": []string{"CM", "CM"},
					"margin":    []string{"1", "1"},
					"style":     []string{"B", "B"},
					"fill":      []string{"255,255,255", "255,255,255"},
				},
				"TOTAL": {
					"alignment": []string{"CM", "CM"},
					"margin":    []string{"1", "1"},
					"style":     []string{"B", "B"},
					"fill":      []string{"255,255,255", "255,255,255"},
				},
			},
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

	doc.Pdf.SetFont("Arial", "", 12)
	_, lineHeight := doc.Pdf.GetFontSize()
	doc.Pdf.Cell(40, lineHeight, doc.CompanyContact.CompanyName)
	if len(doc.CompanyContact.CompanyEmail) > 0 {
		doc.Pdf.Ln(lineHeight + generator.SmallGapY)
		doc.Pdf.Cell(40, lineHeight, doc.CompanyContact.CompanyEmail)
	}

	doc.Pdf.Ln(lineHeight + generator.SmallGapY)
	if len(doc.CompanyContact.CompanyAddress.StreetAddress) > 0 {
		doc.Pdf.Cell(40, lineHeight, fmt.Sprintf("%s", doc.CompanyContact.CompanyAddress.StreetAddress))
	}
	doc.Pdf.Ln(lineHeight + generator.SmallGapY)
	doc.Pdf.Cell(40, lineHeight, fmt.Sprintf("%s %s, %s", doc.CompanyContact.CompanyAddress.PostalCode, doc.CompanyContact.CompanyAddress.City, doc.CompanyContact.CompanyAddress.Country))

	if len(doc.CompanyContact.CompanyPhoneNumber) > 0 {
		doc.Pdf.Ln(lineHeight)
		doc.Pdf.SetFontStyle("I")
		doc.Pdf.Cell(40, lineHeight, fmt.Sprintf("Tel: %s", doc.CompanyContact.CompanyPhoneNumber))
	}

	doc.Pdf.Ln(10)

	headerY := doc.Pdf.GetY()

	doc.Pdf.SetX(generator.MarginX)
	doc.Pdf.SetFont("Arial", "B", generator.LargeTextFontSize)
	doc.Pdf.Cell(40, 20, docType+" TO")
	doc.Pdf.Ln(18)

	doc.Pdf.SetFont("Arial", "", generator.NormalTextFontSize)
	_, lineHeight = doc.Pdf.GetFontSize()
	doc.Pdf.Cell(40, lineHeight, doc.CustomerContact.Name)
	if len(doc.CustomerContact.Email) > 0 {
		doc.Pdf.Ln(lineHeight + generator.SmallGapY)
		doc.Pdf.Cell(40, lineHeight, doc.CustomerContact.Email)
	}

	doc.Pdf.Ln(lineHeight + generator.SmallGapY)
	doc.Pdf.Cell(40, lineHeight, fmt.Sprintf("%s %s, %s", doc.CustomerContact.Address.PostalCode, doc.CustomerContact.Address.City, doc.CustomerContact.Address.Country))

	if len(doc.CustomerContact.PhoneNumber) > 0 {
		doc.Pdf.Ln(lineHeight)
		doc.Pdf.SetFontStyle("I")
		doc.Pdf.Cell(40, lineHeight, fmt.Sprintf("Tel: %s", doc.CustomerContact.PhoneNumber))
	}

	tableY := doc.Pdf.GetY()

	middleX := (safeAreaW / 2) + 15
	doc.Pdf.SetXY(middleX, headerY)
	doc.Pdf.SetFont("Arial", "B", generator.LargeTextFontSize)
	doc.Pdf.Cell(40, 20, docType)
	doc.Pdf.Ln(18)
	doc.Pdf.SetX(middleX)

	x := doc.Pdf.GetX()
	doc.Pdf.SetFont("Arial", "", generator.NormalTextFontSize)
	_, lineHeight = doc.Pdf.GetFontSize()
	doc.Pdf.SetFontStyle("B")
	doc.Pdf.Cell(42, lineHeight, fmt.Sprintf("%v NO: ", docType))
	doc.Pdf.SetFontStyle("")
	doc.Pdf.Cell(40, lineHeight, doc.DocumentData.DocumentNumber)
	doc.Pdf.Ln(lineHeight + generator.GapY + 2)

	doc.Pdf.SetX(x)
	doc.Pdf.SetFontStyle("B")
	doc.Pdf.Cell(42, lineHeight, fmt.Sprintf("%v DATE: ", docType))
	doc.Pdf.SetFontStyle("")
	doc.Pdf.Cell(40, lineHeight, fmt.Sprintf("%s", time.Now().Format("2006-01-02")))

	doc.Pdf.SetXY(generator.MarginX, tableY+10)
	doc.SetTableHeadings(descriptionData)
	err = doc.AddItemToTable(descriptionData)
	if err != nil {
		return err
	}
	return nil
}
