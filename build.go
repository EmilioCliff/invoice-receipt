package main

import (
	"fmt"
	"github/EmilioCliff/invoice-receipt/generator"
	"github/EmilioCliff/invoice-receipt/templates"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
)

// MakeFont generates a font definition file in JSON format. A definition file of this type is
// required to use non-core fonts in the PDF documents that gofpdf generates.
// See the makefont utility in the gofpdf package for a command line interface to this function.

// fontFileStr is the name of the TrueType file (extension .ttf), OpenType file (extension .otf)
//  or binary Type1 file (extension .pfb) from which to generate a definition file.
// If an OpenType file is specified, it must be one that is based on TrueType outlines,
//  not PostScript outlines; this cannot be determined from the file extension alone.
// If a Type1 file is specified, a metric file with the same pathname except with the extension .afm must be present.

// encodingFileStr is the name of the encoding file that corresponds to the font.

// dstDirStr is the name of the directory in which to save the definition file and, if
// embed is true, the compressed font file.

// msgWriter is the writer that is called to display messages throughout the process. Use nil to turn off messages.

// embed is true if the font is to be embedded in the PDF files.

func main() {
	fmt.Printf("Hello World")
}

func New(templateName string, options *generator.Options) (*generator.Document, error) {
	if err := defaults.Set(options); err != nil {
		return &generator.Document{}, err
	}

	validate := validator.New()
	err := validate.Struct(options)
	if err != nil {
		return nil, err
	}

	if options.DocumentType != "Invoice" && options.DocumentType != "Receipt" && options.DocumentType == "" {
		return &generator.Document{}, fmt.Errorf("Invalid document type")
	}

	doc := &generator.Document{
		Type:         options.DocumentType,
		TemplateName: templateName,
		Options:      options,
	}

	doc.AddTemplate("MysticAura", templates.MysticAura)
	doc.AddTemplate("CelestialDream", templates.CelestialDream)
	doc.AddTemplate("AzureEclipse", templates.AzureEclipse)
	doc.AddTemplate("CrimsonWhisper", templates.CrimsonWhisper)
	doc.AddTemplate("UnNamed", templates.UnNamed)

	return doc, nil
}
