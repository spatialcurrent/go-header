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

func ParseHeader(str string) (*Header, error) {
	text := strings.TrimSpace(str)
	copyright, err := ParseCopyright(text)
	if err != nil {
		return nil, err
	}
	license := ParseLicense(text)
	return &Header{
		Text: text,
		Copyright: copyright,
		License: license,
	}, nil
}
