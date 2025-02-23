package smtpserver

import (
	"github.com/SongZihuan/web-watcher/src/logger"
)

func printError(err error) {
	if err == nil {
		return
	}

	logger.Errorf("SMTP Send Error: %s", err.Error())
}

func SendNotify(msg string) {
	printError(Send("网站状态异常期通知", msg))
}
