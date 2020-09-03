package main

import (
	"testing"
)

func BenchmarkGetTime(b *testing.B) {
	GetTime()
}
