// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================
package goheader

import (
	"strconv"
)

func FixHeader(str string, year int) string {
	if year > 0 {
		if yearMatches := YearPattern.FindStringSubmatch(str); len(yearMatches) > 0 {
			str = yearMatches[1] + strconv.Itoa(year) + yearMatches[3]
		}
	}
	return str
}
