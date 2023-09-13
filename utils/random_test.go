package utils

import (
	"testing"
)

func BenchmarkRandomBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomBool()
	}
}

func BenchmarkRandomEmail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomEmail()
	}
}

func BenchmarkRandomInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomInt(9)
	}
}

func BenchmarkRandomIntBetween(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomIntBetween(1, 999999999)
	}
}

func BenchmarkRandomName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomName()
	}
}

func BenchmarkRandomPhone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomPhone()
	}
}

func BenchmarkRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomString(16)
	}
}
