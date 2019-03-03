// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================
package goheader

import (
	"regexp"
	"strings"
)

var ownerPattern = regexp.MustCompile("(?i)Copyright( [(]C[)])? (\\d+)(\\s+)([^-\n]+)")
var emailPattern = regexp.MustCompile("(?i)((.+))[<](.+[@].+[.].+)[>].?")
var allRightsReservedPattern = regexp.MustCompile("(?i)(.+)([\\s-]+)All rights reserved")

func ParseOwner(str string) *Owner {
	if matches := ownerPattern.FindStringSubmatch(str); len(matches) > 0 {
		str2 := strings.TrimSpace(matches[4])
		if emailMatches := emailPattern.FindStringSubmatch(str2); len(emailMatches) > 0 {
			return &Owner{Name: strings.TrimSpace(emailMatches[1]), Email: emailMatches[3]}
		}
		if allRightsReservedMatches := allRightsReservedPattern.FindStringSubmatch(str2); len(allRightsReservedMatches) > 0 {
			return &Owner{Name: strings.TrimSpace(allRightsReservedMatches[1]), Email: ""}
		}
		return &Owner{Name: str2, Email: ""}
	}
	return nil
}
