// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package convert

import (
	"fmt"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/validation"
	structuralschema "k8s.io/apiextensions-apiserver/pkg/apiserver/schema"
)

func Validate(props *apiextensions.JSONSchemaProps) error {
	if validation.SchemaHasInvalidTypes(props) {
		return fmt.Errorf("schema has invalid types")
	}

	ss, err := structuralschema.NewStructural(props)
	if err != nil {
		return err
	}

	errorList := structuralschema.ValidateStructural(nil, ss)
	if len(errorList) > 0 {
		return errorList.ToAggregate()
	}

	return nil
}
