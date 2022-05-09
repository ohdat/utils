package utils

import "testing"

func TestZap(t *testing.T) {
	var logger = NewZapLogger("test")
	logger.Info("testInfo")
}
