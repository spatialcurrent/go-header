// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================
package goheader

import (
	"regexp"
)

var YearPattern = regexp.MustCompile("(?i)(?s)^((?:.*)(?:Copyright(?: [(]C[)])? ))(\\d+)(.*)$")
