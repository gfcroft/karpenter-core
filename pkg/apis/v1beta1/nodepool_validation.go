/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"context"
	"fmt"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/util/validation"
	"knative.dev/pkg/apis"
)

func (in *NodePool) SupportedVerbs() []admissionregistrationv1.OperationType {
	return []admissionregistrationv1.OperationType{
		admissionregistrationv1.Create,
		admissionregistrationv1.Update,
	}
}

func (in *NodePool) Validate(_ context.Context) (errs *apis.FieldError) {
	return errs.Also(
		apis.ValidateObjectMetadata(in).ViaField("metadata"),
		in.Spec.validate().ViaField("spec"),
	)
}

// RuntimeValidate will be used to validate any part of the CRD that can not be validated at CRD creation
func (in *NodePool) RuntimeValidate() (errs *[]error) {
	return &[]error{
		in.Spec.Template.validateLabels().ViaField("spec.template.metadata"),
		in.Spec.Template.Spec.validateTaints().ViaField("spec.template.spec"),
		in.Spec.Template.Spec.validateRequirements().ViaField("spec.template.spec"),
		in.Spec.Template.validateRequirementsNodePoolKeyDoesNotExist().ViaField("spec.template.spec"),
	}
}

func (in *NodePoolSpec) validate() (errs *[]error) {
	return &[]error{
		in.Template.validate().ViaField("template"),
		in.Disruption.validate().ViaField("deprovisioning"),
	}
}

func (in *NodeClaimTemplate) validate() (errs *apis.FieldError) {
	if len(in.Spec.Resources.Requests) > 0 {
		errs = errs.Also(apis.ErrDisallowedFields("resources.requests"))
	}
	return errs.Also(
		in.validateLabels().ViaField("metadata"),
		in.validateRequirementsNodePoolKeyDoesNotExist().ViaField("spec.requirements"),
		in.Spec.validate().ViaField("spec"),
	)
}

// func (in *NodeClaimTemplate) validate() (errs *[]error) {
// 	if len(in.Spec.Resources.Requests) > 0 {
// 		errs = &[]error{apis.ErrDisallowedFields("resources.requests")}
// 	}
// 	//TODO fix up
// 	*errs = append(*errs,([]error{in.validateLabels().ViaField("metadata"),
// 		in.validateRequirementsNodePoolKeyDoesNotExist().ViaField("spec.requirements"),
// 		in.Spec.validate().ViaField("spec"),
// 		})...)
// 	return errs
// }


func (in *NodeClaimTemplate) validateLabels() (errs *apis.FieldError) {
	for key, value := range in.Labels {
		if key == NodePoolLabelKey {
			errs = errs.Also(apis.ErrInvalidKeyName(key, "labels", "restricted"))
		}
		for _, err := range validation.IsQualifiedName(key) {
			errs = errs.Also(apis.ErrInvalidKeyName(key, "labels", err))
		}
		for _, err := range validation.IsValidLabelValue(value) {
			errs = errs.Also(apis.ErrInvalidValue(fmt.Sprintf("%s, %s", value, err), fmt.Sprintf("labels[%s]", key)))
		}
		if err := IsRestrictedLabel(key); err != nil {
			errs = errs.Also(apis.ErrInvalidKeyName(key, "labels", err.Error()))
		}
	}
	return errs
}

func (in *NodeClaimTemplate) validateRequirementsNodePoolKeyDoesNotExist() (errs *apis.FieldError) {
	for i, requirement := range in.Spec.Requirements {
		if requirement.Key == NodePoolLabelKey {
			errs = errs.Also(apis.ErrInvalidArrayValue(fmt.Sprintf("%s is restricted", requirement.Key), "requirements", i))
		}
	}
	return errs
}

func (in *Disruption) validate() (errs *apis.FieldError) {
	if in.ExpireAfter.Duration != nil && *in.ExpireAfter.Duration < 0 {
		return errs.Also(apis.ErrInvalidValue("cannot be negative", "expirationTTL"))
	}
	if in.ConsolidateAfter != nil && in.ConsolidateAfter.Duration != nil && *in.ConsolidateAfter.Duration < 0 {
		return errs.Also(apis.ErrInvalidValue("cannot be negative", "consolidationTTL"))
	}
	if in.ConsolidateAfter != nil && in.ConsolidateAfter.Duration != nil && in.ConsolidationPolicy == ConsolidationPolicyWhenUnderutilized {
		return errs.Also(apis.ErrGeneric("consolidateAfter cannot be combined with consolidationPolicy=WhenUnderutilized"))
	}
	if in.ConsolidateAfter == nil && in.ConsolidationPolicy == ConsolidationPolicyWhenEmpty {
		return errs.Also(apis.ErrGeneric("consolidateAfter must be specified with consolidationPolicy=WhenEmpty"))
	}
	return errs
}
