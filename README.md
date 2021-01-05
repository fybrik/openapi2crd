# openapi2crd

`openapi2crd` is a CLI to generate Kubernetes Custom Resource Definition (CRD) from [OpenAPI 3.0](https://www.openapis.org/).

## Install

```bash
go get github.com/roee88/openapi2crd
```

This will put openapi2crd in `$(go env GOPATH)/bin`. You may need to add that directory to your `$PATH` if you encounter a "command not found" error.

## Usage

```
Outputs a CustomResourceDefinition using the `components.schemas` field of an OpenAPI 3.0 document

Usage:
  openapi2crd [flags]

Flags:
  -h, --help            help for openapi2crd
  -i, --input string    Path to directory with CustomResourceDefinition YAML files (required)
  -o, --output string   Path to output file (required)
  -s, --spec string     Path to OpenAPI 3.0 specification file (required)
```

## Limitations

- Only [structural schemas](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#specifying-a-structural-schema) are allowed with the exception that you can use `$ref` to reference objects defined in the same spec file.

## Acknowledgements

The work is inspired by https://github.com/ant31/crd-validation and https://github.com/kubeflow/crd-validation.
