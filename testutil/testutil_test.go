package testutil

import (
	"errors"
	"testing"
)

func testStepA() error {
	//balabala
	return nil // pass step A
}

func testStepB() error {
	return errors.New("has error") // step b has error
}
func TestTryMoreTime(t *testing.T) {
	//just show test
	TryMoreTime(testStepA, 1, "testStepA")
	TryMoreTime(testStepB, 2, "testStepA")
}
