package testutil

import (
	"testing"
)

func aTest() error {
	return nil
}
func bTest() error {
	//return errors.New("has error")
	return nil
}
func TestTryMoreTime(t *testing.T) {
	//just show test
	TryMoreTime(aTest, 1, "atest")
	TryMoreTime(bTest, 2, "bTest")
}
