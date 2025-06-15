package integration

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Ждём, пока контейнеры запустятся
	time.Sleep(5 * time.Second)
	code := m.Run()
	os.Exit(code)
}

func assertEnv(t *testing.T) {
	assert.NotEmpty(t, os.Getenv("DB_URL"))
	assert.NotEmpty(t, os.Getenv("JWT_SECRET"))
}
