package generator

import (
	"fmt"
	"strconv"
	"strings"
)

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
		doc.Pdf.SetFillColor(r, g, b)
	}

	doc.Pdf.SetFont("Arial", "B", 12)

	doc.Pdf.SetX(MarginX)

	for i := 1; i < len(columnsDescriptions); i++ {
		alignment := columnsDescriptions[i]["alignment"].([]string)
		doc.Pdf.CellFormat(columnsDescriptions[i]["width"].(float64), TableCellLineHeight, columnsDescriptions[i]["columnName"].(string), border[0], 0, alignment[0], true, 0, "")
	}
	doc.Pdf.Ln(-1)
}

func (doc *Document) AddItemToTable(columnsDescriptions map[int]map[string]interface{}) error {
	metaData := columnsDescriptions[0]
	fill := metaData["fillRow"].([]interface{})
	border := metaData["border"].([]string)
	note := metaData["note"].(bool)
	payment := metaData["payment"].(bool)
	calculations := metaData["calculations"].(map[string]map[string][]string)

	if fill[0] == true {
		r := fill[1].(int)
		g := fill[2].(int)
		b := fill[3].(int)
		doc.Pdf.SetFillColor(r, g, b)
	}

	doc.Pdf.SetFont("Arial", "", 12)
	// subtotal := 0.0
	doc.Pdf.SetX(MarginX)

	for idx, item := range doc.Items {
		totalPrice := item.UnitPrice * float64(item.Quantity)
		// subtotal += totalPrice

		if len(columnsDescriptions) > 5 {
			doc.Pdf.CellFormat(columnsDescriptions[1]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%d", idx+1), border[1], 0, columnsDescriptions[1]["alignment"].([]string)[1], true, 0, "")
			doc.Pdf.CellFormat(columnsDescriptions[2]["width"].(float64), TableCellLineHeight, item.Description, border[1], 0, columnsDescriptions[2]["alignment"].([]string)[1], true, 0, "")
			doc.Pdf.CellFormat(columnsDescriptions[3]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%d", item.Quantity), border[1], 0, columnsDescriptions[3]["alignment"].([]string)[1], true, 0, "")
			doc.Pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", item.UnitPrice), border[1], 0, columnsDescriptions[4]["alignment"].([]string)[1], true, 0, "")
			doc.Pdf.CellFormat(columnsDescriptions[5]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", totalPrice), border[1], 0, columnsDescriptions[5]["alignment"].([]string)[1], true, 0, "")
			doc.Pdf.Ln(-1)
		} else {
			doc.Pdf.CellFormat(columnsDescriptions[1]["width"].(float64), TableCellLineHeight, item.Description, border[1], 0, columnsDescriptions[1]["alignment"].([]string)[1], true, 0, "")
			doc.Pdf.CellFormat(columnsDescriptions[2]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%d", item.Quantity), border[1], 0, columnsDescriptions[2]["alignment"].([]string)[1], true, 0, "")
			doc.Pdf.CellFormat(columnsDescriptions[3]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", item.UnitPrice), border[1], 0, columnsDescriptions[3]["alignment"].([]string)[1], true, 0, "")
			doc.Pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", totalPrice), border[1], 0, columnsDescriptions[4]["alignment"].([]string)[1], true, 0, "")
			doc.Pdf.Ln(-1)
		}
	}

	y := doc.Pdf.GetY()

	if doc.DocumentData.Note != "" && note {
		doc.Pdf.Ln(5)
		_, height := doc.Pdf.GetFontSize()

		doc.Pdf.SetFont("Arial", "", ExtraSmallTextFontSize)
		doc.Pdf.MultiCell(75.0, height, doc.DocumentData.Note, "0", "", false)
	} else if payment {
		doc.Pdf.Ln(2)
		doc.Pdf.SetFont("Times", "B", NormalTextFontSize)
		doc.Pdf.Cell(40, CellLineHeight, "PAYMENT Method :")
		doc.Pdf.SetFontStyle("")
		doc.Pdf.Ln(5)

		if doc.Payment.Bank != nil {
			doc.Pdf.Cell(60, CellLineHeight, doc.Payment.Bank.BankName)
			doc.Pdf.Ln(5)
			doc.Pdf.Cell(60, CellLineHeight, fmt.Sprintf("Account Name: %s", doc.Payment.Bank.AccountName))
			doc.Pdf.Ln(5)
			doc.Pdf.Cell(60, CellLineHeight, fmt.Sprintf("Account Number: %s", doc.Payment.Bank.AccountNumber))
		} else if doc.Payment.Paybill != nil {
			doc.Pdf.Cell(60, CellLineHeight, fmt.Sprintf("Paybill Number: %s", doc.Payment.Paybill.PaybillNumber))
			doc.Pdf.Ln(5)
			doc.Pdf.Cell(60, CellLineHeight, fmt.Sprintf("Account Number: %s", doc.Payment.Paybill.AccountNumber))
		} else if doc.Payment.Till != nil {
			doc.Pdf.Cell(60, CellLineHeight, fmt.Sprintf("Till Number: %s", doc.Payment.Till.TillNumber))
		}
	}

	leftIndent := 0.0
	if len(columnsDescriptions) > 5 {
		leftIndent += columnsDescriptions[1]["width"].(float64) + columnsDescriptions[2]["width"].(float64) + columnsDescriptions[3]["width"].(float64)
	} else {
		leftIndent += columnsDescriptions[1]["width"].(float64) + columnsDescriptions[2]["width"].(float64)
	}

	doc.Pdf.SetFontStyle("B")
	doc.Pdf.SetFontSize(12)
	subtotal := doc.calculateSubtotal()
	tax := 0.0
	discount := 0.0
	for key := range calculations {
		alignment := calculations[key]["alignment"]
		margin := calculations[key]["margin"]
		style := calculations[key]["style"]
		fillOne := strings.Split(calculations[key]["fill"][0], ",")
		fillTwo := strings.Split(calculations[key]["fill"][1], ",")
		var r, g, b int

		if strings.ToLower(key) == "subtotal" {
			doc.Pdf.SetXY(MarginX+leftIndent, y)
			if len(columnsDescriptions) > 5 {
				doc.Pdf.SetFontStyle(style[0])

				r, _ = strconv.Atoi(fillOne[0])
				g, _ = strconv.Atoi(fillOne[0])
				b, _ = strconv.Atoi(fillOne[0])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, key, margin[0], 0, alignment[0], true, 0, "")
				doc.Pdf.SetFontStyle(style[1])

				r, _ = strconv.Atoi(fillTwo[1])
				g, _ = strconv.Atoi(fillTwo[1])
				b, _ = strconv.Atoi(fillTwo[1])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[5]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", subtotal), margin[1], 0, alignment[1], true, 0, "")
			} else {
				doc.Pdf.SetFontStyle(style[0])

				r, _ = strconv.Atoi(fillOne[0])
				g, _ = strconv.Atoi(fillOne[0])
				b, _ = strconv.Atoi(fillOne[0])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[3]["width"].(float64), TableCellLineHeight, key, margin[0], 0, alignment[0], true, 0, "")
				doc.Pdf.SetFontStyle(style[1])

				r, _ = strconv.Atoi(fillTwo[1])
				g, _ = strconv.Atoi(fillTwo[1])
				b, _ = strconv.Atoi(fillTwo[1])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", subtotal), margin[1], 0, alignment[1], true, 0, "")
			}
			doc.Pdf.Ln(-1)
			// grandTotal = subtotal
			continue
		}

		// It can either be tax or discount. Now here we change the value of the grandTotal if either is correct
		if strings.ToLower(key) == "tax" {
			tax = doc.calculateTax()
			// tax := 0.0
			// if doc.DocumentData.Tax != 0.0 {
			// 	tax = subtotal * (doc.DocumentData.Tax / 100)
			// 	grandTotal = subtotal + tax
			// }

			doc.Pdf.SetX(MarginX + leftIndent)
			if len(columnsDescriptions) > 5 {
				doc.Pdf.SetFontStyle(style[0])

				r, _ = strconv.Atoi(fillOne[0])
				g, _ = strconv.Atoi(fillOne[0])
				b, _ = strconv.Atoi(fillOne[0])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%v (%.2f%s)", key, doc.DocumentData.Tax, "%"), margin[0], 0, alignment[0], true, 0, "")
				doc.Pdf.SetFontStyle(style[1])

				r, _ = strconv.Atoi(fillTwo[1])
				g, _ = strconv.Atoi(fillTwo[1])
				b, _ = strconv.Atoi(fillTwo[1])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[5]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", tax), margin[1], 0, alignment[1], true, 0, "")
			} else {
				doc.Pdf.SetFontStyle(style[0])

				r, _ = strconv.Atoi(fillOne[0])
				g, _ = strconv.Atoi(fillOne[0])
				b, _ = strconv.Atoi(fillOne[0])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[3]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%v (%.2f%s)", key, doc.DocumentData.Tax, "%"), margin[0], 0, alignment[0], true, 0, "")
				doc.Pdf.SetFontStyle(style[1])

				r, _ = strconv.Atoi(fillTwo[1])
				g, _ = strconv.Atoi(fillTwo[1])
				b, _ = strconv.Atoi(fillTwo[1])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", tax), margin[1], 0, alignment[1], true, 0, "")
			}
			doc.Pdf.Ln(-1)
			continue

		}

		if strings.ToLower(key) == "discount" {
			if tax != 0.0 {
				discount = doc.calculateDiscountWithTax()
			} else {
				discount = doc.calculateDiscount()
			}
			// discount := 0.0
			// if doc.DocumentData.Discount != 0.0 {
			// 	discount = subtotal * (doc.DocumentData.Discount / 100)
			// 	grandTotal = subtotal - discount
			// }

			doc.Pdf.SetX(MarginX + leftIndent)
			if len(columnsDescriptions) > 5 {
				doc.Pdf.SetFontStyle(style[0])

				r, _ = strconv.Atoi(fillOne[0])
				g, _ = strconv.Atoi(fillOne[0])
				b, _ = strconv.Atoi(fillOne[0])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%v (%.2f%s)", key, doc.DocumentData.Discount, "%"), margin[0], 0, alignment[0], true, 0, "")
				doc.Pdf.SetFontStyle(style[1])

				r, _ = strconv.Atoi(fillTwo[1])
				g, _ = strconv.Atoi(fillTwo[1])
				b, _ = strconv.Atoi(fillTwo[1])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[5]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", discount), margin[1], 0, alignment[1], true, 0, "")
			} else {
				doc.Pdf.SetFontStyle(style[0])

				r, _ = strconv.Atoi(fillOne[0])
				g, _ = strconv.Atoi(fillOne[0])
				b, _ = strconv.Atoi(fillOne[0])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[3]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%v (%.2f%s)", key, doc.DocumentData.Discount, "%"), margin[0], 0, alignment[0], true, 0, "")
				doc.Pdf.SetFontStyle(style[1])

				r, _ = strconv.Atoi(fillTwo[1])
				g, _ = strconv.Atoi(fillTwo[1])
				b, _ = strconv.Atoi(fillTwo[1])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", discount), margin[1], 0, alignment[1], true, 0, "")
			}
			doc.Pdf.Ln(-1)
			continue
		}

		if strings.ToLower(key) == "total" {

			grandTotal := (subtotal + tax) - discount
			doc.Pdf.SetX(MarginX + leftIndent)
			if len(columnsDescriptions) > 5 {
				doc.Pdf.SetFontStyle(style[0])

				r, _ = strconv.Atoi(fillOne[0])
				g, _ = strconv.Atoi(fillOne[0])
				b, _ = strconv.Atoi(fillOne[0])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, key, margin[0], 0, alignment[0], true, 0, "")
				doc.Pdf.SetFontStyle(style[1])

				r, _ = strconv.Atoi(fillTwo[1])
				g, _ = strconv.Atoi(fillTwo[1])
				b, _ = strconv.Atoi(fillTwo[1])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[5]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", grandTotal), margin[1], 0, alignment[1], true, 0, "")
			} else {
				doc.Pdf.SetFontStyle(style[0])

				r, _ = strconv.Atoi(fillOne[0])
				g, _ = strconv.Atoi(fillOne[0])
				b, _ = strconv.Atoi(fillOne[0])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[3]["width"].(float64), TableCellLineHeight, key, margin[0], 0, alignment[0], true, 0, "")
				doc.Pdf.SetFontStyle(style[1])

				r, _ = strconv.Atoi(fillTwo[1])
				g, _ = strconv.Atoi(fillTwo[1])
				b, _ = strconv.Atoi(fillTwo[1])
				doc.Pdf.SetFillColor(r, g, b)

				doc.Pdf.CellFormat(columnsDescriptions[4]["width"].(float64), TableCellLineHeight, fmt.Sprintf("%.2f", grandTotal), margin[1], 0, alignment[1], true, 0, "")
			}
			doc.Pdf.Ln(-1)
			continue
		}
	}

	return nil
}
