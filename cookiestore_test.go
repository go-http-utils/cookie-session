package cookiesession

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCookieStore(t *testing.T) {

	t.Run("CookieStore use default options that should be", func(t *testing.T) {
		assert := assert.New(t)
		defer func() {
			err := recover()
			assert.NotNil(err)
		}()

		cs := &CookieStore{
			options: nil,
		}
		val, err := cs.Encode("x", nil)
		assert.Equal("", val)
		t.Log(err)
		assert.NotNil(err)
	})

}
