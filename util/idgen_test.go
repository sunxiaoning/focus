package util

import (
	"testing"
)

func TestIdGen(t *testing.T) {
	for {
		t.Log(IdGenerator.NextID())
		break
	}
}
