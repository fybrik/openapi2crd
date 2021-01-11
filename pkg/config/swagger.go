// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"net/url"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
)

func LoadSwagger(filePath string) (*openapi3.Swagger, error) {

	loader := &openapi3.SwaggerLoader{
		IsExternalRefsAllowed: true,
	}

	uri, err := url.Parse(filePath)
	if err == nil && uri.Scheme != "" && uri.Host != "" {
		return loader.LoadSwaggerFromURI(uri)
	}

	return loader.LoadSwaggerFromFile(filepath.Clean(filePath))
}
