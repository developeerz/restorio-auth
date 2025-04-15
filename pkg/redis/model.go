package redis

type User struct {
	TelegramID int64  `json:"telegram_id"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Telegram   string `json:"telegram"`
	Password   string `json:"password"`
}
