package cookiesession

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCookieStore(t *testing.T) {

	t.Run("CookieStore use Encode options that should be", func(t *testing.T) {
		assert := assert.New(t)
		defer func() {
			err := recover()
			assert.NotNil(err)
		}()

		val, err := Encode(nil)
		assert.Empty(val)
		assert.NotNil(err)
	})

	t.Run("CookieStore with Decode that should be", func(t *testing.T) {
		assert := assert.New(t)
		var i int
		err := Decode("", i)
		assert.NotNil(err)
	})

}
