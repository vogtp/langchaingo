package milvus

import (
	"os"
	"testing"

	"github.com/vogtp/langchaingo/internal/testutil/testctr"
)

func TestMain(m *testing.M) {
	testctr.EnsureTestEnv()
	os.Exit(m.Run())
}
