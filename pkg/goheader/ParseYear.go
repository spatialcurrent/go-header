// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================
package goheader

import (
	"regexp"
	"strconv"
)

var yearPattern = regexp.MustCompile("(?i)Copyright( [(]C[)])? (\\d+)")

func ParseYear(str string) (int, error) {
	if yearMatches := yearPattern.FindStringSubmatch(str); len(yearMatches) > 0 {
		y, err := strconv.Atoi(yearMatches[2])
		if err != nil {
			return -1, err
		}
		return y, nil
	}
	return -1, nil
}
