package templates

import (
	"fmt"
	"github/EmilioCliff/invoice-receipt/generator"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

func UnNamed2(doc *generator.Document) error {
	doc.Pdf = fpdf.New("P", "mm", "A4", "") // 210 x 297 (mm)
	doc.Pdf.SetMargins(generator.MarginX, generator.MarginY, generator.MarginX)

	docType := strings.ToUpper(doc.Type)

	doc.Pdf.SetAutoPageBreak(true, 10)

	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	return fmt.Errorf("Failed to get current dir: %w", err)
	// }

	// doc.Pdf.SetFontLocation(filepath.Join(currentDir, "fonts"))
	// doc.Pdf.AddFont("Pacifico", "", "Pacifico-Regular.json")

	doc.Pdf.SetHeaderFunc(func() {
		doc.Pdf.SetFillColor(200, 200, 200)
		doc.Pdf.Rect(generator.MarginX, generator.MarginY, 190, 10, "F")
	})
	doc.Pdf.SetFooterFunc(func() {
		doc.Pdf.SetFillColor(200, 200, 200)
		doc.Pdf.Rect(generator.MarginX, 277, 190, 10, "F")
	})

	doc.Pdf.AddPage()

	doc.Pdf.Ln(15)

	doc.Pdf.SetFont("Arial", "B", generator.LargeTextFontSize)
	doc.Pdf.Cell(50, generator.CellLineHeight, doc.CompanyContact.CompanyName)
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
	doc.Pdf.Cell(30, generator.CellLineHeight, fmt.Sprintf("%s#", docType))
	doc.Pdf.SetFontStyle("")
	doc.Pdf.CellFormat(50, generator.CellLineHeight, doc.DocumentData.DocumentNumber, "0", 0, "RM", false, 0, "")
	doc.Pdf.Ln(8)

	doc.Pdf.SetX(generator.MarginX + 110)
	doc.Pdf.SetFontStyle("B")
	doc.Pdf.Cell(30, generator.CellLineHeight, fmt.Sprintf("%s DATE", docType))
	doc.Pdf.SetFontStyle("")
	doc.Pdf.CellFormat(50, generator.CellLineHeight, time.Now().Format("02/01/2006"), "0", 0, "RM", false, 0, "")
	doc.Pdf.Ln(8)

	doc.Pdf.SetX(generator.MarginX + 110)
	doc.Pdf.SetFontStyle("B")
	doc.Pdf.Cell(30, generator.CellLineHeight, "DUE DATE")
	doc.Pdf.SetFontStyle("")
	doc.Pdf.CellFormat(50, generator.CellLineHeight, "03/05/2023", "0", 0, "RM", false, 0, "") // Add Due Date Day

	subtotal := 0.0
	for _, item := range doc.Items {
		totalPrice := item.UnitPrice * float64(item.Quantity)
		subtotal += totalPrice
	}

	doc.Pdf.Ln(20)
	y = doc.Pdf.GetY()
	doc.Pdf.Line(generator.MarginX, y, 210-generator.MarginX, y)
	doc.Pdf.Ln(10)
	doc.Pdf.SetFont("Arial", "B", generator.LargeTextFontSize)
	doc.Pdf.Cell(50, generator.CellLineHeight, fmt.Sprintf("%s Total", doc.Type))
	wd := doc.Pdf.GetStringWidth(fmt.Sprintf("%s%v", doc.Options.CurrencySymbol, subtotal))
	doc.Pdf.SetX(-(generator.MarginX + wd))
	doc.Pdf.Cell(wd, generator.CellLineHeight, fmt.Sprintf("%s%v", doc.Options.CurrencySymbol, subtotal))
	doc.Pdf.Ln(15)
	y = doc.Pdf.GetY()
	doc.Pdf.Line(generator.MarginX, y, 210-generator.MarginX, y)

	return nil
}
