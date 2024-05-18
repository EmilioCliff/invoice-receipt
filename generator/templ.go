package generator

import "fmt"

type Templates struct {
	TemplateFunc map[string]func(doc *Document) error `json:"template_func"`
}

func (doc *Document) AddTemplate(templateName string, templateFunc func(doc *Document) error) {
	if doc.Templates == nil {
		doc.Templates = make([]*Templates, 0)
	}

	template := &Templates{
		TemplateFunc: make(map[string]func(doc *Document) error),
	}

	template.TemplateFunc[templateName] = templateFunc
	doc.Templates = append(doc.Templates, template)
}

func (doc *Document) GetTemplate(templateName string) (func(doc *Document) error, error) {
	for _, templates := range doc.Templates {
		if templates != nil {
			if templateFunc, exists := templates.TemplateFunc[templateName]; exists {
				return templateFunc, nil
			}
		}
	}
	return nil, fmt.Errorf("template with index %s not found", templateName)
}
