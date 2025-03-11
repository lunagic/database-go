package tester

import (
	"log"
	"testing"
)

func Logger(t testing.TB) *log.Logger {
	t.Helper()
	return log.New(testWriter{TB: t}, "", log.LstdFlags|log.LUTC)
}

type testWriter struct {
	testing.TB
}

func (tw testWriter) Write(p []byte) (int, error) {
	tw.Helper()
	tw.Logf("%s", p)
	return len(p), nil
}
