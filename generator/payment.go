package generator

type PaymentDetails struct {
	Bank    *BankDetails     `json:"bank"`
	Paybill *PaybillDetails  `json:"paybill"`
	Till    *BuyGoodsDetails `json:"till"`
}

type BankDetails struct {
	BankName      string `json:"bank_name" validate:"required"`
	AccountName   string `json:"bank_account" validate:"required"`
	AccountNumber string `json:"account_name" validate:"required"`
}

type PaybillDetails struct {
	PaybillNumber string `json:"paybill_number" validate:"required"`
	AccountNumber string `json:"account)number" validate:"required"`
}

type BuyGoodsDetails struct {
	TillNumber string `json:"till_number" validate:"required"`
}
