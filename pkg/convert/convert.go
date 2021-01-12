// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package convert

import (
	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

// SchemaPropsToJSONProps converts openapi3.Schema to a JSONProps
func SchemaPropsToJSONProps(schemaRef *openapi3.SchemaRef, spec openapi3.Schemas) *apiextensions.JSONSchemaProps {
	var props *apiextensions.JSONSchemaProps

	if schemaRef == nil {
		return props
	}

	schemaProps := schemaRef.Value

	props = &apiextensions.JSONSchemaProps{
		// ID:               schemaProps.ID,
		// Schema:           apiextensions.JSONSchemaURL(string(schemaRef.Ref.)),
		// Ref:              ref,
		Description:          schemaProps.Description,
		Type:                 schemaProps.Type,
		Format:               schemaProps.Format,
		Title:                schemaProps.Title,
		Maximum:              schemaProps.Max,
		ExclusiveMaximum:     schemaProps.ExclusiveMax,
		Minimum:              schemaProps.Min,
		ExclusiveMinimum:     schemaProps.ExclusiveMin,
		MaxLength:            castUInt64P(schemaProps.MaxLength),
		MinLength:            castUInt64(schemaProps.MinLength),
		Pattern:              schemaProps.Pattern,
		MaxItems:             castUInt64P(schemaProps.MaxItems),
		MinItems:             castUInt64(schemaProps.MinItems),
		UniqueItems:          schemaProps.UniqueItems, // TODO: The field uniqueItems cannot be set to true.
		MultipleOf:           schemaProps.MultipleOf,
		Enum:                 enumJSON(schemaProps.Enum),
		MaxProperties:        castUInt64P(schemaProps.MaxProps),
		MinProperties:        castUInt64(schemaProps.MinProps),
		Required:             schemaProps.Required,
		Items:                schemaToJSONSchemaPropsOrArray(schemaProps.Items, spec),
		AllOf:                schemasToJSONSchemaPropsArray(schemaProps.AllOf, spec),
		OneOf:                schemasToJSONSchemaPropsArray(schemaProps.OneOf, spec),
		AnyOf:                schemasToJSONSchemaPropsArray(schemaProps.AnyOf, spec),
		Not:                  SchemaPropsToJSONProps(schemaProps.Not, spec),
		Properties:           schemasToJSONSchemaPropsMap(schemaProps.Properties, spec),
		AdditionalProperties: schemaToJSONSchemaPropsOrBool(schemaProps.AdditionalProperties, spec),
		// PatternProperties:    schemasToJSONSchemaPropsMap(schemaProps.PatternProperties, spec),
		// AdditionalItems: schemaToJSONSchemaPropsOrBool(schemaProps.AdditionalItems, spec),
	}
	return props
}

func schemasToJSONSchemaPropsArray(schemas openapi3.SchemaRefs, spec openapi3.Schemas) []apiextensions.JSONSchemaProps {
	var s []apiextensions.JSONSchemaProps
	for _, schema := range schemas {
		s = append(s, *SchemaPropsToJSONProps(schema, spec))
	}
	return s
}

// enumJSON converts []interface{} to []JSON
func enumJSON(enum []interface{}) []apiextensions.JSON {
	var s []apiextensions.JSON
	for _, elt := range enum {
		s = append(s, apiextensions.JSON(elt))
	}
	return s
}

func schemaToJSONSchemaPropsOrArray(schema *openapi3.SchemaRef, spec openapi3.Schemas) *apiextensions.JSONSchemaPropsOrArray {
	if schema == nil {
		return nil
	}
	return &apiextensions.JSONSchemaPropsOrArray{
		Schema: SchemaPropsToJSONProps(schema, spec),
	}
}

func schemaToJSONSchemaPropsOrBool(schema *openapi3.SchemaRef, spec openapi3.Schemas) *apiextensions.JSONSchemaPropsOrBool {
	if schema == nil {
		return nil
	}

	return &apiextensions.JSONSchemaPropsOrBool{
		Schema: SchemaPropsToJSONProps(schema, spec),
		Allows: true, // TODO: *schema.Value.AdditionalPropertiesAllowed
	}
}

func schemasToJSONSchemaPropsMap(schemaMap openapi3.Schemas, spec openapi3.Schemas) map[string]apiextensions.JSONSchemaProps {
	m := make(map[string]apiextensions.JSONSchemaProps)
	for key, schema := range schemaMap {
		m[key] = *SchemaPropsToJSONProps(schema, spec)
	}
	return m
}

func castUInt64P(p *uint64) *int64 {
	if p == nil {
		return nil
	}
	val := int64(*p)
	return &val
}

func castUInt64(v uint64) *int64 {
	val := int64(v)
	if val == 0 {
		return nil
	}
	return &val
}
