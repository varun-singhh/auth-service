package models

type AuthorizedResponse struct {
	User    *User       `json:"user,omitempty"`
	Message string      `json:"message"`
	Token   string      `json:"token,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
}

type ResetPassResponse struct {
	Link string `json:"reset_link,omitempty"`
}
