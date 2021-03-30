package version

import (
	"strconv"
	"strings"
)

// Recent return the more recent version
func Recent(v1 string, v2 string) (string, error) {
	var err error
	var n1, n2 int

	if v1 == "" {
		return v2, nil
	}
	v1s := strings.Split(v1, ".")
	v2s := strings.Split(v2, ".")

	for i, part1 := range v1s {
		if i >= len(v2s) {
			return v1, nil
		}
		if n1, err = strconv.Atoi(part1); err != nil {
			return "", err
		}
		if n2, err = strconv.Atoi(v2s[i]); err != nil {
			return "", err
		}
		if n1 > n2 {
			return v1, nil
		} else if n2 > n1 {
			return v2, nil
		}
	}
	return v2, nil
}
