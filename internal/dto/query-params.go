package dto

// User Feed Query Parsms
type FeedQueryParams struct {
	Page   int      `json:"page" validate:"gte=0"`
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
}
