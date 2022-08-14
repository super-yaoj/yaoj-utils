package utils_test

import (
	"bytes"
	"testing"

	utils "github.com/super-yaoj/yaoj-utils"
)

func TestChecksum(t *testing.T) {
	sum := utils.ReaderChecksum(bytes.NewReader([]byte("hello1")))
	t.Log(sum)
}
