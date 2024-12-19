package entity

// invoice entity
type Invoice struct {
	ID         string `json:"id" validate:"required, id"`
	InvoiceURL string `json:"invoice_url"  validate:"required, invoice_url"`
}
