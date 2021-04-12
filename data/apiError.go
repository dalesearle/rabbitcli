package data

type ApiError struct {
	Error string `json:"error"`
	Reason string `json:"reason"`
}
