package main

import (
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ulfox/sysaudit/alerts"
	"github.com/ulfox/sysaudit/audit"
	"github.com/ulfox/sysaudit/utils"
)

func main() {
	logger := logrus.New()
	logger.WithField("SSHD-AUDIT", "v0.0.1").Info("Initiated")

	sigStop := utils.SetupSignalHandler()

	slackClient := alerts.NewSlackAlert("https://hooks.slack.com/services/...")

	queueChan := make(chan string)

	sshd, err := audit.NewSSHDAudit(queueChan, sigStop, logger)
	if err != nil {
		logger.Fatal(err)
	}

	err = sshd.SeekTail()
	if err != nil {
		logger.Fatal(err)
	}

	sshd.StartReading()

slackLoop:
	for {
		select {
		case msg := <-queueChan:
			err := slackClient.SendSlackNotification(msg)
			if err != nil {
				logger.Error(err)
			}
		case <-sigStop:
			break slackLoop
		default:
			if sshd.Err != nil {
				logger.Error(sshd.Err)
				syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			}
			time.Sleep(100 * time.Microsecond)
		}
	}
}
