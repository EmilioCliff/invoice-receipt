package main

import (
	"github.com/go-pdf/fpdf"
)

type Document struct {
	pdf           *fpdf.Fpdf
	Type          string       `json:"type" validate:"required"`
	TemplateIndex int8         `json:"template_index"`
	Templates     []*Templates `json:"templates"`

	CustomerContact *CustomerContact `json:"cutomer_contact"`
	CompanyContact  *CompanyContact  `json:"company_contact"`
	DocumentData    *DocumentData    `json:"document_data"`
	Header          string           `json:"header"`
	Footer          string           `json:"footer"`
	Options         *Options         `json:"options"`
	Items           []*Item          `json:"items"`
}

type Options struct {
	CurrencySymbol string `default:"Ksh" json:"currency_symbol,omitempty"`

	DocumentType string `json:"document_type" validate:"required,oneof=Invoice Receipt"`

	TextItemsNumberTitle          string `default:"No" json:"text_items_number_title,omitempty"`
	TextItemsNameDescriptionTitle string `default:"Name" json:"text_items_name_description_title,omitempty"`
	TextItemsQuantityTitle        string `default:"Quantity" json:"text_items_quantity_title,omitempty"`
	TextItemsUnitCostTitle        string `default:"Unit price" json:"text_items_unit_cost_title,omitempty"`
	TextItemsTotalTitle           string `default:"Total" json:"text_items_total_title,omitempty"`

	BaseTextColor []int `default:"[35,35,35]" json:"base_text_color,omitempty"`
	GreyTextColor []int `default:"[82,82,82]" json:"grey_text_color,omitempty"`
	GreyBgColor   []int `default:"[232,232,232]" json:"grey_bg_color,omitempty"`
	DarkBgColor   []int `default:"[212,212,212]" json:"dark_bg_color,omitempty"`

	Font     string `default:"Helvetica"`
	BoldFont string `default:"Helvetica"`

	Output string `default:"pdf" json:"output" validate:"oneof=pdf bytes"`
}
