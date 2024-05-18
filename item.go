package main

import "fmt"

type Item struct {
	Description string  `json:"description"`
	Quantity    int64   `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

func (doc *Document) SetTableHeadings(columnsDescriptions map[int]map[string]interface{}) {
	metaData := columnsDescriptions[0]
	fill := metaData["fillHeader"].([]interface{})
	border := metaData["border"].([]string)
	if fill[0] == true {
		r := fill[1].(int)
		g := fill[2].(int)
		b := fill[3].(int)
		doc.pdf.SetFillColor(r, g, b)
	}

	doc.pdf.SetFont("Arial", "B", 12)

	doc.pdf.SetX(MarginX)

	for i := 1; i < len(columnsDescriptions); i++ {
		alignment := columnsDescriptions[i]["alignment"].([]string)
		doc.pdf.CellFormat(columnsDescriptions[i]["width"].(float64), TableCellLineHeight, columnsDescriptions[i]["columnName"].(string), border[0], 0, alignment[0], true, 0, "")
	}
	doc.pdf.Ln(-1)
}

func (doc *Document) AddItemToTable(columnsDescriptions map[int]map[string]interface{}) {
	metaData := columnsDescriptions[0]
	fill := metaData["fillRow"].([]interface{})
	border := metaData["border"].([]string)
	if fill[0] == true {
		r := fill[1].(int)
		g := fill[2].(int)
		b := fill[3].(int)
		doc.pdf.SetFillColor(r, g, b)
	}

	doc.pdf.SetFont("Arial", "", 12)
	subtotal := 0.0
	doc.pdf.SetX(MarginX)

	for idx, item := range doc.Items {
		totalPrice := item.UnitPrice * float64(item.Quantity)
		subtotal += totalPrice

		if len(columnsDescriptions) > 5 {
			doc.pdf.CellFormat(columnsDescriptions[1]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%d", idx+1), border[1], 0, columnsDescriptions[1]["alignment"].([]string)[1], true, 0, "")
			doc.pdf.CellFormat(columnsDescriptions[2]["width"].(float64), TableCellLineHeight, item.Description, border[1], 0, columnsDescriptions[2]["alignment"].([]string)[1], true, 0, "")
			doc.pdf.CellFormat(columnsDescriptions[3]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%d", item.Quantity), border[1], 0, columnsDescriptions[3]["alignment"].([]string)[1], true, 0, "")
			doc.pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", item.UnitPrice), border[1], 0, columnsDescriptions[4]["alignment"].([]string)[1], true, 0, "")
			doc.pdf.CellFormat(columnsDescriptions[5]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", totalPrice), border[1], 0, columnsDescriptions[5]["alignment"].([]string)[1], true, 0, "")
			doc.pdf.Ln(-1)
		} else {
			doc.pdf.CellFormat(columnsDescriptions[1]["width"].(float64), TableCellLineHeight, item.Description, border[1], 0, columnsDescriptions[1]["alignment"].([]string)[1], true, 0, "")
			doc.pdf.CellFormat(columnsDescriptions[2]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%d", item.Quantity), border[1], 0, columnsDescriptions[2]["alignment"].([]string)[1], true, 0, "")
			doc.pdf.CellFormat(columnsDescriptions[3]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", item.UnitPrice), border[1], 0, columnsDescriptions[3]["alignment"].([]string)[1], true, 0, "")
			doc.pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", totalPrice), border[1], 0, columnsDescriptions[4]["alignment"].([]string)[1], true, 0, "")
			doc.pdf.Ln(-1)
		}
	}

	y := doc.pdf.GetY()

	if doc.DocumentData.Note != "" {
		doc.pdf.Ln(5)
		_, height := doc.pdf.GetFontSize()

		doc.pdf.SetFont("Arial", "", ExtraSmallTextFontSize)
		doc.pdf.MultiCell(75.0, height, doc.DocumentData.Note, "0", "", false)
	}

	doc.pdf.SetFontStyle("B")
	doc.pdf.SetFontSize(12)
	leftIndent := 0.0
	if len(columnsDescriptions) > 5 {
		leftIndent += columnsDescriptions[1]["width"].(float64) + columnsDescriptions[2]["width"].(float64) + columnsDescriptions[3]["width"].(float64)
	} else {
		leftIndent += columnsDescriptions[1]["width"].(float64) + columnsDescriptions[2]["width"].(float64)
	}

	doc.pdf.SetXY(MarginX+leftIndent, y)
	if len(columnsDescriptions) > 5 {
		doc.pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, "Subtotal", border[1], 0, columnsDescriptions[4]["alignment"].([]string)[1], true, 0, "")
		doc.pdf.CellFormat(columnsDescriptions[5]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", subtotal), border[1], 0, columnsDescriptions[5]["alignment"].([]string)[1], true, 0, "")
	} else {
		doc.pdf.CellFormat(columnsDescriptions[3]["width"].(float64), TableCellLineHeight, "Subtotal", border[1], 0, columnsDescriptions[3]["alignment"].([]string)[1], true, 0, "")
		doc.pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", subtotal), border[1], 0, columnsDescriptions[4]["alignment"].([]string)[1], true, 0, "")
	}
	doc.pdf.Ln(-1)

	grandTotal := subtotal
	doc.pdf.SetX(MarginX + leftIndent)
	if len(columnsDescriptions) > 5 {
		doc.pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, "Grand total", border[1], 0, columnsDescriptions[4]["alignment"].([]string)[1], true, 0, "")
		doc.pdf.CellFormat(columnsDescriptions[5]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", grandTotal), border[1], 0, columnsDescriptions[5]["alignment"].([]string)[1], true, 0, "")
	} else {
		doc.pdf.CellFormat(columnsDescriptions[3]["width"].(float64), TableCellLineHeight, "Grand total", border[1], 0, columnsDescriptions[3]["alignment"].([]string)[1], true, 0, "")
		doc.pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", grandTotal), border[1], 0, "CM", true, 0, "")
	}
	doc.pdf.Ln(-1)
}
