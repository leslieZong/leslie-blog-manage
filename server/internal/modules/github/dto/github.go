package dto

import "time"

type Activity struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Repo      string    `json:"repo"`
	CreatedAt time.Time `json:"createdAt"`
}

type ActivityResponse struct {
	Items []*Activity `json:"items"`
}
