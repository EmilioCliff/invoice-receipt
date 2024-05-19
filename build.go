package main

import (
	"fmt"
	"github/EmilioCliff/invoice-receipt/generator"
	templates "github/EmilioCliff/invoice-receipt/templates"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
)

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

	return doc, nil
}
