package service

import "testing"

func TestGetExceptionOverview(t *testing.T) {
	t.Run("TestGetExceptionOverview", func(t *testing.T) {
		GetExceptionOverview("watchdog_store_buying_center_control", 0)
	})
}
