package domain

import (
	"fmt"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Print(os.Getenv("GOPATH"))
}