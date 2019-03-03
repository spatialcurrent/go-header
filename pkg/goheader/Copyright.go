// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================
package goheader

type Copyright struct {
	Owner *Owner
	Year  int
}

func (c Copyright) Map() map[string]interface{} {
	m := map[string]interface{}{}
	if c.Owner != nil {
		m["owner"] = c.Owner.Map()
	}
	if c.Year != -1 {
		m["year"] = c.Year
	}
	return m
}
