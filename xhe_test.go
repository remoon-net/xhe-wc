package main

import (
	"testing"

	"github.com/lainio/err2/try"
)

func TestXxx(t *testing.T) {
	s := "fdd9::/24"
	ip := try.To1(parseIP(s))
	t.Log(ip)
}
