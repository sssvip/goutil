package strutil

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewUUID(t *testing.T) {
	assert.Equal(t, 36, len(NewUUID()))
}

func TestMd5(t *testing.T) {
	assert.Equal(t, "5d41402abc4b2a76b9719d911017c592", Md5("hello"))
}

func TestAtoIWithDefault(t *testing.T) {
	assert.Equal(t, 12321, AtoIWithDefault("12321", 1))
	assert.Equal(t, 1, AtoIWithDefault("12321t", 1))
}

func TestAtoFloat64WithDefault(t *testing.T) {
	assert.Equal(t, float64(12321), AtoFloat64WithDefault("12321", 1))
	assert.Equal(t, float64(1), AtoFloat64WithDefault("12321t", 1))
}

func TestSafeCutString(t *testing.T) {
	assert.Equal(t, "100", SafeCutString("100testsa", 3))
	assert.Equal(t, "100testsa", SafeCutString("100testsa", 100))
	assert.Equal(t, "100testsa", SafeCutString("100testsa", -1))
	assert.Equal(t, "", SafeCutString("", 3))
}
