package models


type CookieUpdateRequest struct {
	Cookie string `json:"cookie"`
}

type CookieUpdateResponse struct {
	Message string  `json:"message"`
	Success bool 	`json:"success"`
}

