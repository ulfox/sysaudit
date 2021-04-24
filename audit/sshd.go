package audit

import (
	"fmt"

	"github.com/coreos/go-systemd/sdjournal"
	"github.com/sirupsen/logrus"
)

type SSHD struct {
	Journal *sdjournal.Journal
	Queue   chan<- string
	SigStop <-chan struct{}
	Logger  *logrus.Logger
	Err     error
}

func NewSSHDAudit(qC chan<- string, sk <-chan struct{}, l *logrus.Logger) (*SSHD, error) {
	auditJournal, err := sdjournal.NewJournal()
	if err != nil {
		return nil, err
	}

	sshd := &SSHD{
		Journal: auditJournal,
		Queue:   qC,
		SigStop: sk,
		Logger:  l,
	}
	return sshd, nil
}

func (s *SSHD) SeekTail() error {
	return s.Journal.SeekTail()
}

func (s *SSHD) StartReading() {
	go func() {
		for {
			select {
			case <-s.SigStop:
				return
			default:
				n, err := s.Journal.Next()
				if err != nil {
					s.Err = err
					return
				}
				if n < 1 {
					s.Journal.Wait(sdjournal.IndefiniteWait)
					continue
				}
				s.ReadJournal()
			}
		}
	}()
}

func (s *SSHD) Next() error {
	for {
		n, err := s.Journal.Next()
		if err != nil {
			return err
		}
		if n < 1 {
			s.Journal.Wait(sdjournal.IndefiniteWait)
			continue
		}
		break
	}
	return nil
}

func (s *SSHD) ReadJournal() {
	sshd, err := s.Journal.GetDataValue("_SYSTEMD_UNIT")
	if err != nil {
		s.Logger.Warn(err)
		return
	}

	if sshd != "sshd.service" && sshd != "systemd-logind.service" {
		s.Logger.Infof("Entry %s is not a whitelisted entry. Skipping", sshd)
		return
	}

	msg, err := s.Journal.GetDataValue("MESSAGE")
	if err != nil {
		s.Logger.Error(err)
		return
	}

	hostname, err := s.Journal.GetDataValue("_HOSTNAME")
	if err != nil {
		s.Logger.Error(err)
		return
	}

	pid, err := s.Journal.GetDataValue("_PID")
	if err != nil {
		s.Logger.Error(err)
		return
	}

	report := fmt.Sprintf("%v %v[%v] %v\n", hostname, pid, msg, sshd)
	s.Logger.Info(report)

	s.Queue <- report
}
