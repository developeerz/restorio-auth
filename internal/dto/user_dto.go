package dto

type SignUpRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Telegram  string `json:"telegram"`
	Password  string `json:"password"`
}

type VerificationRequest struct {
	Code     int    `json:"code"`
	Telegram string `json:"telegram"`
}

type LoginRequest struct {
	Telegram string `json:"telegram"`
	Password string `json:"password"`
}

type JwtAccess struct {
	Access string `json:"access"`
}
