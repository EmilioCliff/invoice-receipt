package templates

import (
	"fmt"
	"github/EmilioCliff/invoice-receipt/generator"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

func UnNamed(doc *generator.Document) error {
	doc.Pdf = fpdf.New("P", "mm", "A4", "") // 210 x 297 (mm)
	doc.Pdf.SetMargins(generator.MarginX, generator.MarginY, generator.MarginX)

	docType := strings.ToUpper(doc.Type)

	doc.Pdf.SetAutoPageBreak(true, 10)

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Failed to get current dir: %w", err)
	}

	doc.Pdf.SetFontLocation(filepath.Join(currentDir, "fonts"))
	doc.Pdf.AddFont("Pacifico", "", "Pacifico-Regular.json")

	doc.Pdf.AddPage()

	targetWidth, targetHeight, err := generator.ResizeImage(doc.CompanyContact.CompanyLogo)
	if err != nil {
		return err
	}

	doc.Pdf.ImageOptions(doc.CompanyContact.CompanyLogo, generator.MarginX, generator.MarginY, targetWidth, targetHeight, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	doc.Pdf.Ln(50)
	y := doc.Pdf.GetY()
	doc.Pdf.SetXY(-(generator.MarginX + 50), generator.MarginY)
	doc.Pdf.SetFont("Arial", "B", generator.LargeTextFontSize)
	doc.Pdf.CellFormat(50, generator.CellLineHeight, docType, "0", 0, "RM", false, 0, "")

	doc.Pdf.SetY(y)
	doc.Pdf.Ln(15)
	doc.Pdf.SetFontSize(generator.NormalTextFontSize)
	y = doc.Pdf.GetY()
	doc.Pdf.Cell(50, generator.CellLineHeight, doc.CustomerContact.Name)
	doc.Pdf.Ln(5)
	if doc.CustomerContact.Email != "" {
		doc.Pdf.Cell(50, generator.CellLineHeight, doc.CustomerContact.Email)
		doc.Pdf.Ln(5)
	}
	if doc.CustomerContact.PhoneNumber != "" {
		doc.Pdf.Cell(50, generator.CellLineHeight, doc.CustomerContact.PhoneNumber)
		doc.Pdf.Ln(5)
	}
	if doc.CustomerContact.Address.StreetAddress != "" {
		street := strings.Split(doc.CustomerContact.Address.StreetAddress, ",")
		doc.Pdf.Cell(50, generator.CellLineHeight, fmt.Sprintf("%s %s", street[0], street[2]))
		doc.Pdf.Ln(5)
	}
	if doc.CustomerContact.Address != nil {
		doc.Pdf.Cell(50, generator.CellLineHeight, fmt.Sprintf("%s %s", doc.CustomerContact.Address.PostalCode, doc.CustomerContact.Address.City))
		doc.Pdf.Ln(5)
		doc.Pdf.Cell(50, generator.CellLineHeight, doc.CustomerContact.Address.Country)
	}

	wd := doc.Pdf.GetStringWidth(fmt.Sprintf("%s number: %s", doc.Type, doc.DocumentData.DocumentNumber))
	doc.Pdf.SetXY(-(wd + generator.MarginX), y)

	ln := doc.Pdf.GetStringWidth(fmt.Sprintf("%s number:", doc.Type))
	doc.Pdf.CellFormat(ln, generator.CellLineHeight, fmt.Sprintf("%s number:", doc.Type), "0", 0, "RM", false, 0, "")
	doc.Pdf.SetFontStyle("")
	doc.Pdf.CellFormat(wd-ln, generator.CellLineHeight, doc.DocumentData.DocumentNumber, "0", 0, "RM", false, 0, "")
	doc.Pdf.Ln(5)

	doc.Pdf.SetFontStyle("B")
	wd = doc.Pdf.GetStringWidth(fmt.Sprintf("%s date: %s", doc.Type, time.Now().Format("02/01/2006")))
	doc.Pdf.SetX(-(wd + generator.MarginX))
	ln = doc.Pdf.GetStringWidth(fmt.Sprintf("%s date:", doc.Type))
	doc.Pdf.CellFormat(ln, generator.CellLineHeight, fmt.Sprintf("%s date:", doc.Type), "0", 0, "RM", false, 0, "")
	doc.Pdf.SetFontStyle("")
	doc.Pdf.CellFormat(wd-ln, generator.CellLineHeight, time.Now().Format("02/01/2006"), "0", 0, "RM", false, 0, "")
	doc.Pdf.Ln(5)

	doc.Pdf.SetFontStyle("B")
	wd = doc.Pdf.GetStringWidth(fmt.Sprintf("Payment terms: %v days", doc.DocumentData.DueDate))
	doc.Pdf.SetX(-(wd + generator.MarginX))
	ln = doc.Pdf.GetStringWidth("Payment terms:")
	doc.Pdf.CellFormat(ln, generator.CellLineHeight, "Payment terms:", "0", 0, "RM", false, 0, "")
	doc.Pdf.SetFontStyle("")
	doc.Pdf.CellFormat(wd-ln, generator.CellLineHeight, fmt.Sprintf("%v days", doc.DocumentData.DueDate), "0", 0, "RM", false, 0, "")
	doc.Pdf.Ln(5)

	doc.Pdf.SetFontStyle("B")
	dueDate := time.Now().Add(time.Duration(doc.DocumentData.DueDate))
	wd = doc.Pdf.GetStringWidth(fmt.Sprintf("Due date: %s", dueDate.Format("02/01/2006")))
	doc.Pdf.SetX(-(wd + generator.MarginX))
	ln = doc.Pdf.GetStringWidth("Due date:")
	doc.Pdf.CellFormat(ln, generator.CellLineHeight, "Due date:", "0", 0, "RM", false, 0, "")
	doc.Pdf.SetFontStyle("")
	doc.Pdf.CellFormat(wd-ln, generator.CellLineHeight, fmt.Sprintf("%v", dueDate.Format("02/01/2006")), "0", 0, "RM", false, 0, "")
	doc.Pdf.Ln(30)

	if doc.DocumentData.Note != "" {
		doc.Pdf.MultiCell(100, 4, doc.DocumentData.Note, "0", "SL", false)
		doc.Pdf.Ln(10)
	}

	descriptionData := map[int]map[string]interface{}{
		0: {
			"fillHeader": []interface{}{true, 200, 200, 200},
			"fillRow":    []interface{}{true, 255, 255, 255},
			"border":     []string{"0", "B"},
			"calculations": map[string]map[string][]string{
				"Subtotal": {
					"alignment": []string{"RM", "RM"},
					"margin":    []string{"B", "B"},
					"style":     []string{"B", ""},
					"fill":      []string{"255,255,255", "255,255,255"},
				},
				"Tax": {
					"alignment": []string{"RM", "RM"},
					"margin":    []string{"B", "B"},
					"style":     []string{"B", ""},
					"fill":      []string{"255,255,255", "255,255,255"},
				},
				"TOTAL": {
					"alignment": []string{"RM", "RM"},
					"margin":    []string{"B", "B"},
					"style":     []string{"B", "B"},
					"fill":      []string{"255,255,255", "255,255,255"},
				},
			},
			"note":    false,
			"payment": true,
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
	doc.AddItemToTable(descriptionData)

	doc.Pdf.SetY(-45)
	y = doc.Pdf.GetY()

	doc.Pdf.SetFont("Arial", "B", generator.LargeTextFontSize)
	doc.Pdf.Cell(100, generator.CellLineHeight, "Thank you for purchase!")

	doc.Pdf.SetFontSize(generator.NormalTextFontSize)
	doc.Pdf.SetX(130)
	doc.Pdf.SetFont("Pacifico", "", generator.LargeTextFontSize)
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
