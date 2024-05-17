package main

import (
	"fmt"
)

type Item struct {
	Description string
	Quantity    int64
	UnitPrice   float64
}

func (doc *Document) SetTableHeadings() {
	header := [NumberOfColums]string{doc.Options.TextItemsNumberTitle,
		doc.Options.TextItemsNameDescriptionTitle,
		doc.Options.TextItemsQuantityTitle,
		fmt.Sprintf("%s (%s)", doc.Options.TextItemsUnitCostTitle, doc.Options.CurrencySymbol),
		fmt.Sprintf("%s (%s)", doc.Options.TextItemsTotalTitle, doc.Options.CurrencySymbol),
	}
	colWidth := [NumberOfColums]float64{NumberColumnOffset, DescriptionColumnOffset, QuantityColumnOffset, UnitPriceColumnOffset, TotalPriceOffset}

	doc.pdf.SetFont("Arial", "B", 12)
	doc.pdf.SetFillColor(200, 200, 200)

	doc.pdf.SetX(MarginX)

	for colJ := 0; colJ < NumberOfColums; colJ++ {
		doc.pdf.CellFormat(colWidth[colJ], TableCellLineHeight, header[colJ], "1", 0, "CM", true, 0, "")
	}
	doc.pdf.Ln(-1)
}

func (doc *Document) AddItemToTable() {
	colWidth := [NumberOfColums]float64{NumberColumnOffset, DescriptionColumnOffset, QuantityColumnOffset, UnitPriceColumnOffset, TotalPriceOffset}
	doc.pdf.SetFillColor(255, 255, 255)

	doc.pdf.SetFont("Arial", "", 12)
	subtotal := 0.0

	doc.pdf.SetX(MarginX)
	for idx, item := range doc.Items {
		totalPrice := item.UnitPrice * float64(item.Quantity)
		subtotal += totalPrice

		doc.pdf.CellFormat(colWidth[0], TableCellLineHeight, fmt.Sprintf("%d", idx+1), "1", 0, "CM", true, 0, "")
		doc.pdf.CellFormat(colWidth[1], TableCellLineHeight, item.Description, "1", 0, "LM", true, 0, "")
		doc.pdf.CellFormat(colWidth[2], TableCellLineHeight, fmt.Sprintf("%d", item.Quantity), "1", 0, "CM", true, 0, "")
		doc.pdf.CellFormat(colWidth[3], TableCellLineHeight, fmt.Sprintf("%.2f", item.UnitPrice), "1", 0, "CM", true, 0, "")
		doc.pdf.CellFormat(colWidth[4], TableCellLineHeight, fmt.Sprintf("%.2f", totalPrice), "1", 0, "CM", true, 0, "")
		doc.pdf.Ln(-1)

	}

	doc.pdf.SetFontStyle("B")
	leftIndent := 0.0
	for i := 0; i < 3; i++ {
		leftIndent += colWidth[i]
	}

	doc.pdf.SetX(MarginX + leftIndent)
	doc.pdf.CellFormat(colWidth[3], TableCellLineHeight, "Subtotal", "1", 0, "CM", true, 0, "")
	doc.pdf.CellFormat(colWidth[4], TableCellLineHeight, fmt.Sprintf("%.2f", subtotal), "1", 0, "CM", true, 0, "")
	doc.pdf.Ln(-1)

	grandTotal := subtotal
	doc.pdf.SetX(MarginX + leftIndent)
	doc.pdf.CellFormat(colWidth[3], TableCellLineHeight, "Grand total", "1", 0, "CM", true, 0, "")
	doc.pdf.CellFormat(colWidth[4], TableCellLineHeight, fmt.Sprintf("%.2f", grandTotal), "1", 0, "CM", true, 0, "")
	doc.pdf.Ln(-1)
}

// Description string
// Quantity    int64
// UnitPrice   float64
