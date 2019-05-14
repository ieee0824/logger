package main

import (
	"os"

	"github.com/ieee0824/logger"
)

func main() {
	l := logger.NewLogger()

	os.Setenv("ENV", logger.Dev.String())

	l.Infof("info: dev\n")
	l.Warnf("warn: dev\n")
	l.Errof("err: dev\n")
}
