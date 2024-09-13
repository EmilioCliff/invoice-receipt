package generator

func (doc *Document) SetFooter(text string) {
	doc.Footer = text
	return
}

func (doc *Document) SetCompanyAddress(contact *CompanyContact) {
	doc.CompanyContact = contact
	return
}

func (doc *Document) SetCustomerAddress(contact *CustomerContact) {
	doc.CustomerContact = contact
	return
}

func (doc *Document) SetDocumentData(data *DocumentData) {
	doc.DocumentData = data
	return
}

func (doc *Document) AddItem(item *Item) {
	doc.Items = append(doc.Items, item)
	return
}

func (doc *Document) SetPaymentDetails(paymentDetails *PaymentDetails) {
	doc.Payment = paymentDetails
	return
}
