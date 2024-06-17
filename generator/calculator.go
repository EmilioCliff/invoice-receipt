package generator

func (doc *Document) calculateSubtotal() float64 {
	subtotal := 0.0
	for _, item := range doc.Items {
		totalPrice := item.UnitPrice * float64(item.Quantity)
		subtotal += totalPrice
	}
	return subtotal
}

func (doc *Document) calculateTax() float64 {
	subtotal := doc.calculateSubtotal()
	tax := 0.0
	if doc.DocumentData.Tax != 0.0 {
		tax = subtotal * (doc.DocumentData.Tax / 100)
	}

	return tax
}

func (doc *Document) calculateDiscount() float64 {
	subtotal := doc.calculateSubtotal()
	discount := 0.0
	if doc.DocumentData.Discount != 0.0 {
		discount = subtotal * (doc.DocumentData.Discount / 100)
	}
	return discount
}

func (doc *Document) calculateDiscountWithTax() float64 {
	totalWithTax := doc.CalculateTotalWithTax()
	discount := 0.0
	if doc.DocumentData.Discount != 0.0 {
		discount = totalWithTax * (doc.DocumentData.Discount / 100)
	}
	return discount
}

func (doc *Document) CalculateTotalWithDiscount() float64 {
	return doc.calculateSubtotal() - doc.calculateDiscount()
}

func (doc *Document) CalculateTotalWithTax() float64 {
	return doc.calculateSubtotal() + doc.calculateTax()
}

func (doc *Document) CalculateTotalWithTaxAndDiscount() float64 {
	return (doc.calculateSubtotal() + doc.calculateTax()) - doc.calculateDiscountWithTax()
}

func (doc *Document) CalculateTotalWithoutTaxAndDiscount() float64 {
	return doc.calculateSubtotal()
}
