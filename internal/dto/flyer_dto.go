package dto

import "time"

// FlyerData is the DTO for the data extracted from the flyer image by Gemini.
type FlyerData struct {
	StoreInfo         Store       `json:"store"`
	CampaignInfo      Campaign    `json:"campaign"`
	FlyerItemsInfo    []FlyerItem `json:"flyer_items"`
	DisplayExpiryDate *time.Time  `json:"display_expiry_date,omitempty"`
}

type Store struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Prefecture string `json:"prefecture"`
	City       string `json:"city"`
	Street     string `json:"street"`
}

type Campaign struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type FlyerItem struct {
	Product             Product `json:"product"`
	PriceExcludingTax   int     `json:"price_excluding_tax"`
	PriceIncludingTax   int     `json:"price_including_tax"`
	Unit                string  `json:"unit"`
	RestrictionNote     string  `json:"restriction_note"`
}

type Product struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}
