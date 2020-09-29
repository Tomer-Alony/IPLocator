package models

type IP struct {
	Country string `json:"country"`
	City string `json:"city"`
}

type IPService interface {
	FindCountry(ip string) (*IP, error)
}
