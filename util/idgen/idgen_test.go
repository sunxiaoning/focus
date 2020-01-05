package idgenutil

import (
	"testing"
)

func TestIdGen(t *testing.T) {
	for {
		t.Log(NextId())
		break
	}
}
