// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================
package goheader

type Owner struct {
	Name string
  Email string
}

func (o Owner) Map() map[string]interface{} {
	m := map[string]interface{}{
		"name": o.Name,
	}
	if len(o.Email) > 0 {
		m["email"] = o.Email
	}
	return m
}
