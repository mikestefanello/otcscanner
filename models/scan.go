package models

// Scan describes input provided on the scan form which is used to update orders
type Scan struct {
	Barcode   string `json:"barcode" validate:"required"`
	Country   string `json:"country" validate:"required"`
	Weight    string `json:"weight" validate:"required,numeric,gt=0"`
	Length    string `json:"length" validate:"required,numeric,gt=0"`
	Width     string `json:"width" validate:"required,numeric,gt=0"`
	Height    string `json:"height" validate:"required,numeric,gt=0"`
	Date      string `json:"date" validate:"required"`
	Service   string `json:"service" validate:"required"`
	Account   string `json:"account" validate:"required"`
	CreateNew bool   `json:"createNew"`
}
