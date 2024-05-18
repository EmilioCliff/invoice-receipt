package generator

import (
	"bytes"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func (doc *Document) Build() ([]byte, error) {
	validate := validator.New()
	err := validate.Struct(doc)
	if err != nil {
		return nil, err
	}

	template, err := doc.GetTemplate(doc.TemplateName)
	if err != nil {
		return nil, err
	}
	err = template(doc)

	if doc.Options.Output == "pdf" {
		if err := doc.Pdf.OutputFileAndClose(fmt.Sprintf("%s.pdf", doc.DocumentData.DocumentNumber)); err != nil {
			return nil, err
		}
	} else {
		var buffer bytes.Buffer
		if err := doc.Pdf.Output(&buffer); err != nil {
			return nil, err
		}

		doc.Pdf.Close()
		return buffer.Bytes(), nil
	}
	return nil, nil
}
