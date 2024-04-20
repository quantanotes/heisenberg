package ioutil

import (
	"log/slog"
	"time"
)

type Syncer struct {
	stream syncable
	ticker time.Ticker
	logger *slog.Logger
}

type syncable interface {
	Sync() error
}

func NewSyncer(stream syncable, interval time.Duration, logger *slog.Logger) Syncer {
	return Syncer{
		stream: stream,
		ticker: *time.NewTicker(interval),
		logger: logger,
	}
}

func (s *Syncer) Start() {
	go s.run()
}

func (s *Syncer) Stop() {
	s.ticker.Stop()
}

func (s *Syncer) run() {
	for range s.ticker.C {
		if err := s.stream.Sync(); err != nil {
			s.logger.Error(err.Error())
		}
	}
}
