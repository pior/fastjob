package fastjob

import (
	"context"
	"fmt"
	"log"
)

type Logger interface {
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, err error, format string, args ...interface{})
}

type standardLogger struct{}

func (s *standardLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	log.Printf("DEBUG: "+format, args...)
}

func (s *standardLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	log.Printf("INFO: "+format, args...)
}

func (s *standardLogger) Errorf(ctx context.Context, err error, format string, args ...interface{}) {
	log.Printf("ERROR: %s: %s", fmt.Sprintf(format, args...), err)
}
