// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package convert

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stoewer/go-strcase"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

// SchemaPropsToJSONProps converts openapi3.Schema to a JSONProps
func SchemaPropsToJSONProps(schemaRef *openapi3.SchemaRef) *apiextensions.JSONSchemaProps {
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
		UniqueItems:          false, // The field uniqueItems cannot be set to true.
		MultipleOf:           schemaProps.MultipleOf,
		Enum:                 enumJSON(schemaProps.Enum),
		MaxProperties:        castUInt64P(schemaProps.MaxProps),
		MinProperties:        castUInt64(schemaProps.MinProps),
		Required:             schemaProps.Required,
		Items:                schemaToJSONSchemaPropsOrArray(schemaProps.Items),
		AllOf:                schemasToJSONSchemaPropsArray(schemaProps.AllOf),
		OneOf:                schemasToJSONSchemaPropsArray(schemaProps.OneOf),
		AnyOf:                schemasToJSONSchemaPropsArray(schemaProps.AnyOf),
		Not:                  SchemaPropsToJSONProps(schemaProps.Not),
		Properties:           schemasToJSONSchemaPropsMap(schemaProps.Properties),
		AdditionalProperties: schemaToJSONSchemaPropsOrBool(schemaProps.AdditionalProperties),
		// PatternProperties:    schemasToJSONSchemaPropsMap(schemaProps.PatternProperties),
		// AdditionalItems: schemaToJSONSchemaPropsOrBool(schemaProps.AdditionalItems),
	}

	// Apply custom transformations
	props = transformations(props, schemaRef)

	return props
}

func transformations(props *apiextensions.JSONSchemaProps, schemaRef *openapi3.SchemaRef) *apiextensions.JSONSchemaProps {
	result := props
	result = oneOfRefsTransform(result, schemaRef.Value.OneOf)
	result = removeUnknownFormats(result)
	return result
}

func removeUnknownFormats(props *apiextensions.JSONSchemaProps) *apiextensions.JSONSchemaProps {
	switch props.Format {
	case "int32", "int64", "float", "double", "byte", "date", "date-time", "password":
		// Valid formats https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#format
	default:
		props.Format = ""
	}
	return props
}

// oneOfRefsTransform transforms oneOf with a list of $ref to a list of nullable properties
func oneOfRefsTransform(props *apiextensions.JSONSchemaProps, oneOf openapi3.SchemaRefs) *apiextensions.JSONSchemaProps {
	if props.OneOf != nil && len(props.Properties) == 0 && props.AdditionalProperties == nil {
		result := props.DeepCopy()
		result.Type = "object"
		result.OneOf = nil

		options := []apiextensions.JSON{}
		for _, v := range oneOf {
			if v.Ref == "" {
				// this transform does not apply here
				// return the original props
				return props
			}
			name := v.Ref
			name = name[strings.LastIndex(name, "/")+1:]
			name = strcase.LowerCamelCase(name)
			options = append(options, name)
			result.Properties[name] = *SchemaPropsToJSONProps(v)
		}

		result.Properties["type"] = apiextensions.JSONSchemaProps{
			Type:        "string",
			Enum:        options,
			Description: "Type is the discriminator for the different possible values",
		}

		return result
	}
	return props
}

func schemasToJSONSchemaPropsArray(schemas openapi3.SchemaRefs) []apiextensions.JSONSchemaProps {
	var s []apiextensions.JSONSchemaProps
	for _, schema := range schemas {
		s = append(s, *SchemaPropsToJSONProps(schema))
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

func schemaToJSONSchemaPropsOrArray(schema *openapi3.SchemaRef) *apiextensions.JSONSchemaPropsOrArray {
	if schema == nil {
		return nil
	}
	return &apiextensions.JSONSchemaPropsOrArray{
		Schema: SchemaPropsToJSONProps(schema),
	}
}

func schemaToJSONSchemaPropsOrBool(schema *openapi3.SchemaRef) *apiextensions.JSONSchemaPropsOrBool {
	if schema == nil {
		return nil
	}

	return &apiextensions.JSONSchemaPropsOrBool{
		Schema: SchemaPropsToJSONProps(schema),
		Allows: true, // TODO: *schema.Value.AdditionalPropertiesAllowed
	}
}

func schemasToJSONSchemaPropsMap(schemaMap openapi3.Schemas) map[string]apiextensions.JSONSchemaProps {
	m := make(map[string]apiextensions.JSONSchemaProps)
	for key, schema := range schemaMap {
		m[key] = *SchemaPropsToJSONProps(schema)
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
