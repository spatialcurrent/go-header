// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================
package goheader

import (
	"strings"
)

func RemoveEmptyLines(str string) string {
	output := make([]string, 0)
	for _, line := range strings.Split(str, "\n") {
		if len(line) > 0 {
			output = append(output, line)
		}
	}
	return strings.Join(output, "\n")
}
