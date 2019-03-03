// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================
package goheader

import (
	"regexp"
)

var licensePattern = regexp.MustCompile("(?s)([\\w-]+)([\\s\n]+)license")

func ParseLicense(str string) string {
	if licenseMatches := licensePattern.FindStringSubmatch(str); len(licenseMatches) > 0 {
		return licenseMatches[1]
	}
	return ""
}
