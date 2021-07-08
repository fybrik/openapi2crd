// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"

	"github.com/mesh-for-data/openapi2crd/pkg/convert"
)

type Generator struct {
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(crd *apiextensions.CustomResourceDefinition, spec openapi3.Schemas) (*apiextensions.CustomResourceDefinition, error) {
	validation, err := getCustomResourceValidation(crd.Spec.Names.Kind, spec)
	if err != nil {
		return nil, err
	}
	crd.Spec.Validation = validation

	// A workaround because ValidateCRD requires at least one stored version in the status.
	// Otherwise the following error is raised:
	// status.storedVersions: Invalid value: []string{}: must have at least one stored version
	for _, version := range crd.Spec.Versions {
		if version.Storage {
			crd.Status.StoredVersions = append(crd.Status.StoredVersions, version.Name)
		}
	}

	if err := ValidateCRD(crd); err != nil {
		return nil, err
	}
	// TODO: yaml.Marshal creates an empty status field that we should remove
	// StoredVersions is set to empty array instead of nil to bypass the following issue:
	// https://github.com/mesh-for-data/openapi2crd/issues/1
	crd.Status.StoredVersions = []string{}

	return crd, nil
}

// getCustomResourceValidation returns the validation definition for a CRD name
func getCustomResourceValidation(name string, spec openapi3.Schemas) (*apiextensions.CustomResourceValidation, error) {
	// Fix known types (ref: https://github.com/kubernetes/kubernetes/issues/62329)
	spec["k8s.io/apimachinery/pkg/util/intstr.IntOrString"] = openapi3.NewSchemaRef("", &openapi3.Schema{
		AnyOf: openapi3.SchemaRefs{
			{
				Value: openapi3.NewStringSchema(),
			},
			{
				Value: openapi3.NewIntegerSchema(),
			},
		},
	})

	schema := spec[name]
	props := convert.SchemaPropsToJSONProps(schema)

	return &apiextensions.CustomResourceValidation{
		OpenAPIV3Schema: props,
	}, nil
}
