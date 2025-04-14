package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const modePerm os.FileMode = 0666

func InitLogger() error {
	file, err := os.OpenFile(
		"/var/log/restorio-auth/log.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		modePerm,
	)
	if err != nil {
		return err
	}

	log.Logger = zerolog.New(file).With().Timestamp().Logger()

	return nil
}
