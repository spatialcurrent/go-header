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

func ParseCopyright(str string) (*Copyright, error) {
	text := strings.TrimSpace(str)
	owner := ParseOwner(text)
	year, err := ParseYear(text)
	if err != nil {
		return nil, err
	}
	return &Copyright{
		Owner: owner,
		Year: year,
	}, nil
}
