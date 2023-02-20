package domain

import (
	"time"
)

type Offer struct {
	OfferID     string     `json:"offer_id"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
