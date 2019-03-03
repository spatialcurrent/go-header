// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================
package goheader

import (
	"strings"
)

func RemoveLinePrefix(str string, prefix string) string {
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, prefix) {
			lines[i] = strings.TrimSpace(line[len(prefix):])
		}
	}
	return strings.Join(lines, "\n")
}
