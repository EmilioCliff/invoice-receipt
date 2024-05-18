package main

import "fmt"

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
