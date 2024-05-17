package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeUtilLoadLocation(t *testing.T) {
	t.Run("return loaded location", func(t *testing.T) {
		t.Setenv("TIMEZONE", "Asia/Bangkok")
		loc := LoadLocation()

		assert.Equal(t, "Asia/Bangkok", loc.String())
	})

	t.Run("return default location", func(t *testing.T) {
		t.Setenv("TIMEZONE", "")
		loc := LoadLocation()

		assert.Equal(t, "Asia/Jakarta", loc.String())
	})

	t.Run("return local location when failed load location zone", func(t *testing.T) {
		t.Setenv("TIMEZONE", "XXXYYY-ZZZ")
		loc := LoadLocation()

		assert.Equal(t, time.Local.String(), loc.String())
	})
}
