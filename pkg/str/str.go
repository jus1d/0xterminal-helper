package str

import (
	"strconv"
)

func Unescape(s string) string {
	quoted := `"` + s + `"`
	unescaped, _ := strconv.Unquote(quoted)
	return unescaped
}
