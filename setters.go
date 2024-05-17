package main

func (doc *Document) SetHeaders(text string) *Document {
	doc.Header = text
	return doc
}

func (doc *Document) SetFooter(text string) *Document {
	doc.Footer = text
	return doc
}

func (doc *Document) SetCompanyAddress(contact *CompanyContact) *Document {
	doc.CompanyContact = contact
	return doc
}

func (doc *Document) SetCustomerAddress(contact *CustomerContact) *Document {
	doc.CustomerContact = contact
	return doc
}

func (doc *Document) SetDocumentData(data *DocumentData) *Document {
	doc.DocumentData = data
	return doc
}

func (doc *Document) AddItem(item *Item) *Document {
	doc.Items = append(doc.Items, item)
	return doc
}
