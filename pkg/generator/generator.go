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

func (g *Generator) Generate(original *apiextensions.CustomResourceDefinition, spec openapi3.Schemas) *apiextensions.CustomResourceDefinition {
	original.Spec.Validation = getCustomResourceValidation(original.Spec.Names.Kind, spec)
	return original
}

// getCustomResourceValidation returns the validation definition for a CRD name
func getCustomResourceValidation(name string, spec openapi3.Schemas) *apiextensions.CustomResourceValidation {
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
	return &apiextensions.CustomResourceValidation{
		OpenAPIV3Schema: convert.SchemaPropsToJSONProps(schema, spec),
	}
}
