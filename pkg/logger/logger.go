package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const modePerm os.FileMode = 0666

func InitLogger(serviceName string) error {
	file, err := os.OpenFile(
		fmt.Sprintf("/var/log/%s/log.log", serviceName),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		modePerm,
	)
	if err != nil {
		return fmt.Errorf("logger init: %w", err)
	}

	log.Logger = zerolog.New(file).With().Timestamp().Logger()

	return nil
}
