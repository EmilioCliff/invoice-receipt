package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-pdf/fpdf"
)

type Templates struct {
	TemplateFunc map[int8]func(doc *Document) error `json:"template_func"`
}

func (doc *Document) AddTemplate(templateIndex int8, templateFunc func(doc *Document) error) {
	if doc.Templates == nil {
		doc.Templates = make([]*Templates, 0)
	}

	template := &Templates{
		TemplateFunc: make(map[int8]func(doc *Document) error),
	}

	template.TemplateFunc[templateIndex] = templateFunc
	doc.Templates = append(doc.Templates, template)
}

func (doc *Document) GetTemplate(templateIndex int8) (func(doc *Document) error, error) {
	for _, templates := range doc.Templates {
		if templates != nil {
			if templateFunc, exists := templates.TemplateFunc[templateIndex]; exists {
				return templateFunc, nil
			}
		}
	}
	return nil, fmt.Errorf("template with index %d not found", templateIndex)
}

func MysticAura(doc *Document) error {
	doc.pdf = fpdf.New("P", "mm", "A4", "") // 210 x 297 (mm)
	doc.pdf.SetMargins(MarginX, MarginY, MarginX)

	docType := strings.ToUpper(doc.Type)
	if doc.Footer != "" {
		doc.pdf.SetFooterFunc(func() {
			doc.pdf.SetY(-MarginY)
			doc.pdf.SetFont("Arial", "I", ExtraSmallTextFontSize)
			wd := doc.pdf.GetStringWidth(doc.Footer)
			doc.pdf.Cell(wd, CellLineHeight, doc.Footer)
			doc.pdf.SetX(-MarginX)
			doc.pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", doc.pdf.PageNo()),
				"", 0, "C", false, 0, "")
		})
	}

	doc.pdf.SetAutoPageBreak(true, 20)

	doc.pdf.AddPage()

	pageW, _ := doc.pdf.GetPageSize()
	safeAreaW := pageW - 2*MarginX

	targetWidth, targetHeight, err := ResizeImage(doc.CompanyContact.CompanyLogo)
	if err != nil {
		return err
	}

	doc.pdf.ImageOptions(doc.CompanyContact.CompanyLogo, safeAreaW-80, MarginY, targetWidth, targetHeight, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	doc.pdf.SetXY(MarginX, MarginY)
	doc.pdf.SetFont("Arial", "B", ExtraLargeTextFontSize)
	doc.pdf.Cell(40, 20, docType)
	doc.pdf.Ln(15)

	doc.CompanyContact.LayerCompanyContact(doc)
	doc.pdf.Ln(10)

	headerY := doc.pdf.GetY()

	doc.pdf.SetX(MarginX)
	doc.pdf.SetFont("Arial", "B", LargeTextFontSize)
	doc.pdf.Cell(40, 20, docType+" TO")
	doc.pdf.Ln(15)

	doc.CustomerContact.LayerCustomerContact(doc)
	tableY := doc.pdf.GetY()

	middleX := safeAreaW / 2
	doc.pdf.SetXY(middleX, headerY+10)
	doc.DocumentData.LayerDocumentData(docType, doc)

	doc.pdf.SetXY(MarginX, tableY+10)
	doc.SetTableHeadings()
	doc.AddItemToTable()
	return nil
}

func CelestialDream(doc *Document) error {
	doc.pdf = fpdf.New("P", "mm", "A4", "") // 210 x 297 (mm)
	doc.pdf.SetMargins(MarginX, MarginY, MarginX)
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return err
	}

	doc.pdf.AddPage()

	doc.pdf.Image(filepath.Join(currentDir, "logo.png"), 6, 6, 30, 0, false, "", 0, "")
	doc.pdf.SetFont("Arial", "", 16)
	doc.pdf.Text(40, 20, doc.DocumentData.Note)
	doc.pdf.CellFormat(0, 6, "this is a new chapter for tpl 1", "", 1, "L", false, 0, "")
	doc.pdf.SetDrawColor(0, 100, 200)
	doc.pdf.SetLineWidth(2.5)
	doc.pdf.Line(95, 12, 105, 22)

	return nil
}
