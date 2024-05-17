package main

import (
	"bytes"
	"fmt"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
)

func New(templateIndex int8, options *Options) (*Document, error) {
	if err := defaults.Set(options); err != nil {
		return &Document{}, err
	}

	validate := validator.New()
	err := validate.Struct(options)
	if err != nil {
		return nil, err
	}

	if options.DocumentType != "Invoice" && options.DocumentType != "Receipt" && options.DocumentType == "" {
		return &Document{}, fmt.Errorf("Invalid document type")
	}

	doc := &Document{
		Type:          options.DocumentType,
		TemplateIndex: templateIndex,
		Options:       options,
	}

	doc.AddTemplate(1, MysticAura)
	doc.AddTemplate(2, CelestialDream)

	return doc, nil
}

func (doc *Document) Build() ([]byte, error) {
	template, err := doc.GetTemplate(doc.TemplateIndex)
	if err != nil {
		return nil, err
	}
	err = template(doc)

	if doc.Options.Output == "pdf" {
		if err := doc.pdf.OutputFileAndClose(fmt.Sprintf("%s.pdf", doc.DocumentData.DocumentNumber)); err != nil {
			return nil, err
		}
	} else {
		var buffer bytes.Buffer
		if err := doc.pdf.Output(&buffer); err != nil {
			return nil, err
		}

		doc.pdf.Close()
		return buffer.Bytes(), nil
	}
	return nil, nil
}

// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"

// 	"github.com/go-pdf/fpdf"
// )

//	func main() {
//		pdf := fpdf.New("P", "mm", "A4", "")
//		titleStr := "20000 Leagues Under the Seas"
//		pdf.SetTitle(titleStr, false)
//		pdf.SetAuthor("Jules Verne", false)
//		pdf.SetHeaderFunc(func() {
//			// Arial bold 15
//			pdf.SetFont("Arial", "B", 15)
//			// Calculate width of title and position
//			wd := pdf.GetStringWidth(titleStr) + 6
//			pdf.SetX((210 - wd) / 2)
//			// Colors of frame, background and text
//			pdf.SetDrawColor(0, 80, 180)
//			pdf.SetFillColor(230, 230, 0)
//			pdf.SetTextColor(220, 50, 50)
//			// Thickness of frame (1 mm)
//			pdf.SetLineWidth(1)
//			// Title
//			pdf.CellFormat(wd, 9, titleStr, "1", 1, "C", true, 0, "")
//			// Line break
//			pdf.Ln(1)
//		})
//		pdf.SetFooterFunc(func() {
//			// Position at 1.5 cm from bottom
//			pdf.SetY(-15)
//			// Arial italic 8
//			pdf.SetFont("Arial", "I", 8)
//			// Text color in gray
//			pdf.SetTextColor(128, 128, 128)
//			// Page number
//			pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
//				"", 0, "C", false, 0, "")
//		})
//		chapterTitle := func(chapNum int, titleStr string) {
//			// 	// Arial 12
//			pdf.SetFont("Arial", "", 12)
//			// Background color
//			pdf.SetFillColor(200, 220, 255)
//			// Title
//			pdf.CellFormat(0, 6, fmt.Sprintf("Chapter %d : %s", chapNum, titleStr),
//				"", 1, "L", true, 0, "")
//			// Line break
//			pdf.Ln(4)
//		}
//		chapterBody := func(fileStr string) {
//			// Read text file
//			txtStr, err := os.ReadFile(fileStr)
//			if err != nil {
//				pdf.SetError(err)
//			}
//			// Times 12
//			pdf.SetFont("Times", "", 12)
//			// Output justified text
//			pdf.MultiCell(0, 5, string(txtStr), "1", "", false)
//			// Line break
//			pdf.Ln(-1)
//			// Mention in italics
//			pdf.SetFont("", "I", 0)
//			pdf.Cell(0, 5, "(end of excerpt)")
//		}
//		printChapter := func(chapNum int, titleStr, fileStr string) {
//			pdf.AddPage()
//			chapterTitle(chapNum, titleStr)
//			chapterBody(fileStr)
//		}
//		currentDir, err := os.Getwd()
//		if err != nil {
//			fmt.Println("Error getting current working directory:", err)
//			return
//		}
//		printChapter(1, "A RUNAWAY REEF", filepath.Join(currentDir, "20k_c1.txt"))
//		printChapter(2, "THE PROS AND CONS", filepath.Join(currentDir, "20k_c1.txt"))
//		err = pdf.OutputFileAndClose("out.pdf")
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"

// 	"github.com/go-pdf/fpdf"
// )

// func main() {
// 	currentDir, err := os.Getwd()
// 	if err != nil {
// 		fmt.Println("Error getting current working directory:", err)
// 		return
// 	}
// 	pdf := fpdf.New("P", "mm", "A4", "")
// 	pdf.SetCompression(false)
// 	// pdf.SetFont("Times", "", 12)
// 	template := pdf.CreateTemplate(func(tpl *fpdf.Tpl) {
// 		tpl.Image(filepath.Join(currentDir, "logo.png"), 6, 6, 30, 0, false, "", 0, "")
// 		tpl.SetFont("Arial", "B", 16)
// 		tpl.Text(40, 20, "Template says hello")
// 		tpl.CellFormat(0, 6, "this is a new chapter for tpl 1", "", 1, "L", false, 0, "")
// 		tpl.SetDrawColor(0, 100, 200)
// 		tpl.SetLineWidth(2.5)
// 		tpl.Line(95, 12, 105, 22)
// 	})
// 	_, tplSize := template.Size()
// 	// fmt.Println("Size:", tplSize)
// 	// fmt.Println("Scaled:", tplSize.ScaleBy(1.5))

// 	template2 := pdf.CreateTemplate(func(tpl *fpdf.Tpl) {
// 		tpl.UseTemplate(template)
// 		subtemplate := tpl.CreateTemplate(func(tpl2 *fpdf.Tpl) {
// 			tpl2.Image(filepath.Join(currentDir, "logo.png"), 6, 86, 30, 0, false, "", 0, "")
// 			tpl2.SetFont("Arial", "B", 16)
// 			tpl2.Text(40, 100, "Subtemplate says hello")
// 			tpl.Ln(1)
// 			tpl.CellFormat(0, 6, "this is a new chapter for tpl 2", "", 1, "L", false, 0, "")
// 			tpl2.SetDrawColor(0, 200, 100)
// 			tpl2.SetLineWidth(2.5)
// 			tpl2.Line(102, 92, 112, 102)
// 		})
// 		tpl.UseTemplate(subtemplate)
// 	})

// 	pdf.SetDrawColor(200, 100, 0)
// 	pdf.SetLineWidth(2.5)
// 	pdf.SetFont("Arial", "B", 16)

// 	// serialize and deserialize template
// 	// b, _ := template2.Serialize()
// 	// template3, _ := fpdf.DeserializeTemplate(b)

// 	pdf.AddPage()
// 	// pdf.UseTemplate(template)
// 	// pdf.UseTemplateScaled(template, fpdf.PointType{X: 0, Y: 30}, tplSize)
// 	// pdf.Line(40, 210, 60, 210)
// 	// pdf.Text(40, 200, "Template example page 1")

// 	// pdf.AddPage()
// 	// pdf.UseTemplate(template2)
// 	pdf.UseTemplateScaled(template2, fpdf.PointType{X: 10, Y: 30}, tplSize.ScaleBy(0.5))
// 	pdf.Line(60, 210, 80, 210)
// 	pdf.Text(40, 200, "Template example page 2")

// 	err = pdf.OutputFileAndClose("outp.pdf")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
