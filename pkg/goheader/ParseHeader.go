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

func ParseHeader(raw string) (*Header, error) {
	contents := strings.TrimSpace(RemoveEmptyLines(RemoveLinePrefix(raw, "//")))
	copyright, err := ParseCopyright(contents)
	if err != nil {
		return nil, err
	}
	license := ParseLicense(contents)
	return &Header{
		Raw:       raw,
		Contents:  contents,
		Copyright: copyright,
		License:   license,
	}, nil
}
