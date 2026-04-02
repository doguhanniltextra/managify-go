package middleware

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func InitLogger() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}
