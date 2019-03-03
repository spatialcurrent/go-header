// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================
package goheader

type File struct {
	Path      string
	Name      string
	Extension string
	Header    *Header
}

func (f File) Map() map[string]interface{} {
	m := map[string]interface{}{
		"path":      f.Path,
		"name":      f.Name,
		"extension": f.Extension,
	}
	if f.Header != nil {
		m["header"] = f.Header.Map()
	}
	return m
}
