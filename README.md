# openapi2crd

`openapi2crd` is a CLI to generate Kubernetes Custom Resource Definition (CRD) resources from [OpenAPI 3.0](https://www.openapis.org/).

## Install

Download the appropriate version for your platform from [Releases](https://github.com/mesh-for-data/openapi2crd/releases/latest). You may want to install the binary to somewhere in your system's PATH such as `/usr/local/bin`.

Alternatively, if you have go 1.13 or later then you can also use `go get`. This will put `crdoc` in `$(go env GOPATH)/bin`:

```bash
go get github.com/mesh-for-data/openapi2crd
```

## Usage

1. Create an input directory with YAML files of `CustomResourceDefinition` resources without schema information (see [example/input](example/input)).
1. Create an OpenAPI 3.0 document with `components.schemas` (see [example/spec.yaml](example/spec.yaml))
    * The document must include a schema with the name identical to the `kind` of each input `CustomResourceDefinition`. 
    * The document must comply with the listed [limitations](#limitations)
1. Invoke `openapi2crd` command:
    ```bash
    openapi2crd SPEC_FILE --input INPUT_DIR --output OUTPUT_FILE
    ```

An output YAML file will be generated in the specified output location (see [example/output/output.yaml](example/output/output.yaml))

## Limitations

- Only [structural schemas](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#specifying-a-structural-schema) are allowed with the exception that you can use `$ref`.

## Acknowledgements

The work is inspired by https://github.com/ant31/crd-validation and https://github.com/kubeflow/crd-validation.
