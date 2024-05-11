package serr

import (
	"fmt"
	"strings"
)

// ArrayToString conveniently returns a slice of str items
// as a nicely formatted string. If a limit greater than 0 is supplied,
// the output is nicely truncated as necessary with ellipses
func ArrayToString(strArr []string, delim string, limit int, listName string) (out string) {
	lnArr := len(strArr)

	if lnArr == 0 {
		out = "0 " + listName
		return
	}

	if limit > 0 && lnArr > limit {
		out = fmt.Sprintf("%d %s: %s", lnArr, listName, strings.Join(strArr[:limit], delim)+"...")
	} else {
		out = strings.Join(strArr, delim)
	}
	return
}
