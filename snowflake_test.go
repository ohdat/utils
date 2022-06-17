package utils

import "testing"

func TestNextSnowflakeID(t *testing.T) {
	id := NextSnowflakeID()
	t.Log(id)
}
