// Package validate_test
package validate_test

import (
	"github.com/gogf/gf/v2/test/gtest"
	"hotgo/utility/validate"
	"testing"
)

func TestIsEmail(t *testing.T) {
	b := validate.IsEmail("QTT123456@163.com")
	gtest.Assert(true, b)
}
