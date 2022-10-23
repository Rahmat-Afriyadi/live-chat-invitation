package dto

type Status struct {
	Present int64 `json:"present"`
	NPresent int64 `json:"npresent"`
	Hesitant int64 `json:"hesitant"`
}