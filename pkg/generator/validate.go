// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/validation"
)

func ValidateCRD(crd *apiextensions.CustomResourceDefinition) error {
	errorList := validation.ValidateCustomResourceDefinition(crd, apiextensionsv1.SchemeGroupVersion)
	if len(errorList) > 0 {
		return errorList.ToAggregate()
	}
	return nil
}
