package examples

import (
	"os"
	"testing"

	tms "github.com/testit-tms/adapters-go"
)

func TestMain(m *testing.M) {
	os.Exit(tms.Run(m))
}
