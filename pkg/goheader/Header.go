// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================
package goheader

type Header struct {
	Raw       string
	Contents  string
	Copyright *Copyright
	License   string
}

func (h Header) Map() map[string]interface{} {
	m := map[string]interface{}{
		"raw":      h.Raw,
		"contents": h.Contents,
	}
	if h.Copyright != nil {
		m["copyright"] = h.Copyright.Map()
	}
	if len(h.License) > 0 {
		m["license"] = h.License
	}
	return m
}
