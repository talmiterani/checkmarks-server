package models

type SearchReq struct {
	PageSize int    `json:"page_size"`
	Page     int    `json:"page"`
	Query    string `json:"query,omitempty"`
}
