package generator

import (
	"fmt"
	"image"
	"os"
	"time"
)

type Address struct {
	PostalCode    string `json:"postal_code,omitempty" validate:"required"`
	City          string `json:"city,omitempty" validate:"required"`
	Country       string `json:"country,omitempty" validate:"required"`
	StreetAddress string `json:"street_address,omitempty"`
}

type CustomerContact struct {
	Name        string   `json:"name" validate:"required"`
	Email       string   `json:"email"`
	PhoneNumber string   `json:"phone_number"`
	Address     *Address `json:"address"`
}

type CompanyContact struct {
	CompanyName        string   `json:"company_name" validate:"required"`
	CompanyEmail       string   `json:"company_email" validate:"required"`
	CompanyPhoneNumber string   `json:"company_phone_number" validate:"required"`
	CompanyLogo        string   `json:"company_logo" validate:"required"`
	CompanyAddress     *Address `json:"company_address" validate:"required"`
}

type DocumentData struct {
	DocumentNumber   string  `json:"document_data,omitempty"`
	Note             string  `json:"note,omitempty"`
	IssuedBy         string  `json:"issued_by,omitempty"`
	IssuedByPosition string  `json:"issued_by_position,omitempty"`
	Tax              float64 `json:"tax"`
	Discount         float64 `json:"discount"`
	DueDate          int8    `json:"due_date"`
	TermNConditions  string  `json:"terms_n_conditions"`
}

func ResizeImage(imagePath string) (float64, float64, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Error opening image:", err)
		return 0, 0, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return 0, 0, err
	}

	imgWidth := float64(img.Bounds().Dx())
	imgHeight := float64(img.Bounds().Dy())

	aspectRatio := imgWidth / imgHeight

	var targetWidth, targetHeight float64
	if MaximunImageWidth/aspectRatio <= MaximumImageHeight {
		targetWidth = MaximunImageWidth
		targetHeight = MaximunImageWidth / aspectRatio
	} else {
		targetWidth = MaximumImageHeight * aspectRatio
		targetHeight = MaximumImageHeight
	}

	return targetWidth, targetHeight, nil
}

func structureAddress(a *Address) string {
	var addrString string
	if len(a.StreetAddress) > 0 {
		addrString += fmt.Sprintf("%s\n", a.StreetAddress)
	}

	if len(a.PostalCode) > 0 {
		addrString += fmt.Sprintf("P.O Box - %s\n", a.PostalCode)
	}

	if len(a.City) > 0 {
		addrString += fmt.Sprintf("%s, %s\n", a.City, a.Country)
	}

	return addrString
}

func (c *CustomerContact) LayerCustomerContact(doc *Document) {
	doc.Pdf.SetFont("Arial", "", NormalTextFontSize)
	_, lineHeight := doc.Pdf.GetFontSize()
	doc.Pdf.Cell(40, lineHeight, c.Name)
	if len(c.Email) > 0 {
		doc.Pdf.Ln(lineHeight + SmallGapY)
		doc.Pdf.Cell(40, lineHeight, c.Email)
	}

	if len(c.Address.City) > 0 {
		doc.Pdf.Ln(lineHeight + SmallGapY)
		structuredAddress := structureAddress(c.Address)
		fn := doc.Pdf.UnicodeTranslatorFromDescriptor("")
		doc.Pdf.MultiCell(100, lineHeight+SmallGapY, fn(structuredAddress), "0", "", false)
	}

	if len(c.PhoneNumber) > 0 {
		doc.Pdf.Ln(lineHeight)
		doc.Pdf.SetFontStyle("I")
		doc.Pdf.Cell(40, lineHeight, fmt.Sprintf("Tel: %s", c.PhoneNumber))
	}
}

func (c *CompanyContact) LayerCompanyContact(doc *Document) {
	doc.Pdf.SetFont("Arial", "", 12)
	_, lineHeight := doc.Pdf.GetFontSize()
	doc.Pdf.Cell(40, lineHeight, c.CompanyName)
	if len(c.CompanyEmail) > 0 {
		doc.Pdf.Ln(lineHeight + SmallGapY)
		doc.Pdf.Cell(40, lineHeight, c.CompanyEmail)
	}

	if len(c.CompanyAddress.City) > 0 {
		doc.Pdf.Ln(lineHeight + SmallGapY)
		structuredAddress := structureAddress(c.CompanyAddress)
		fn := doc.Pdf.UnicodeTranslatorFromDescriptor("")
		doc.Pdf.MultiCell(100, lineHeight+SmallGapY, fn(structuredAddress), "0", "", false)
	}

	if len(c.CompanyPhoneNumber) > 0 {
		doc.Pdf.Ln(lineHeight)
		doc.Pdf.SetFontStyle("I")
		doc.Pdf.Cell(40, lineHeight, fmt.Sprintf("Tel: %s", c.CompanyPhoneNumber))
	}

}

func (d *DocumentData) LayerDocumentData(docType string, doc *Document, formats map[string]interface{}) {
	names := formats["names"].([]string)
	formart := formats["formats"].(map[string]string)
	x := doc.Pdf.GetX()
	doc.Pdf.SetFont("Arial", "", NormalTextFontSize)
	_, lineHeight := doc.Pdf.GetFontSize()
	doc.Pdf.SetFontStyle(formart["font_style"])
	doc.Pdf.Cell(42, lineHeight, names[0])
	doc.Pdf.SetFontStyle("")
	doc.Pdf.Cell(40, lineHeight, d.DocumentNumber)
	doc.Pdf.Ln(lineHeight + GapY + 2)

	doc.Pdf.SetX(x)
	doc.Pdf.SetFontStyle(formart["font_style"])
	doc.Pdf.Cell(42, lineHeight, names[1])
	doc.Pdf.SetFontStyle("")
	doc.Pdf.Cell(40, lineHeight, fmt.Sprintf("%s", time.Now().Format(formart["date_format"])))
}

func (doc *Document) SetPageFooter() {
	if doc.Footer != "" {
		doc.Pdf.SetFooterFunc(func() {
			doc.Pdf.SetY(-MarginY)
			doc.Pdf.SetFont("Arial", "I", ExtraSmallTextFontSize)
			wd := doc.Pdf.GetStringWidth(doc.Footer)
			doc.Pdf.Cell(wd, CellLineHeight, doc.Footer)
			doc.Pdf.SetX(-MarginX)
			doc.Pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", doc.Pdf.PageNo()),
				"", 0, "C", false, 0, "")
		})
	}
}
