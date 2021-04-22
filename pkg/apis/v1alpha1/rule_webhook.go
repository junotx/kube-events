/*
Copyright 2020 The KubeSphere Authors.

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

package v1alpha1

import (
	"fmt"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

type ruleScopeConfig struct {
	ScopeLabelKey            string
	ScopeLabelValueCluster   string
	ScopeLabelValueWorkspace string
	ScopeLabelValueNamespace string

	ScopeWorkspaceLabelKey string

	NamespaceScopeCluster   string
	NamespaceScopeWorkspace string
}

func DefaultRuleScopeConfig() *ruleScopeConfig {
	return &ruleScopeConfig{
		ScopeLabelKey:            "kubesphere.io/rule-scope",
		ScopeLabelValueCluster:   "cluster",
		ScopeLabelValueWorkspace: "workspace",
		ScopeLabelValueNamespace: "namespace",

		ScopeWorkspaceLabelKey: "kubesphere.io/workspace",

		NamespaceScopeCluster:   "kubesphere-logging-system",
		NamespaceScopeWorkspace: "kubesphere-logging-system",
	}
}

var _ruleScopeConfig = DefaultRuleScopeConfig()

func SetRuleScopeConfig(c *ruleScopeConfig) {
	_ruleScopeConfig = c
}
func GetRuleScopeConfig() *ruleScopeConfig {
	return _ruleScopeConfig
}

// log is for logging in this package.
var rulelog = logf.Log.WithName("rule-resource")

func (r *Rule) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-events-kubesphere-io-v1alpha1-rule,mutating=true,failurePolicy=fail,groups=events.kubesphere.io,resources=rules,verbs=create;update,versions=v1alpha1,name=mrule.kb.io

var _ webhook.Defaulter = &Rule{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Rule) Default() {
	rulelog.Info("default", "name", r.Name)

	switch r.Namespace {
	case _ruleScopeConfig.NamespaceScopeCluster, _ruleScopeConfig.NamespaceScopeWorkspace:
	default:
		if r.Labels == nil {
			r.Labels = make(map[string]string)
		}
		if _, ok := r.Labels[_ruleScopeConfig.ScopeLabelKey]; !ok {
			r.Labels[_ruleScopeConfig.ScopeLabelKey] = "namespace"
		}
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:path=/validate-events-kubesphere-io-v1alpha1-rule,mutating=false,failurePolicy=fail,groups=events.kubesphere.io,resources=rules,verbs=create;update,versions=v1alpha1,name=vrule.kb.io

var _ webhook.Validator = &Rule{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Rule) ValidateCreate() error {
	rulelog.Info("validate create", "name", r.Name)

	return r.validateRule()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Rule) ValidateUpdate(old runtime.Object) error {
	rulelog.Info("validate update", "name", r.Name)

	return r.validateRule()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Rule) ValidateDelete() error {
	rulelog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *Rule) validateRule() error {
	var allErrs field.ErrorList
	if err := r.validateRuleScope(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrs.NewInvalid(
		schema.GroupKind{Group: GroupVersion.Group, Kind: "Rule"},
		r.Name, allErrs)
}

func (r *Rule) validateRuleScope() *field.Error {
	rs, ok := r.Labels[_ruleScopeConfig.ScopeLabelKey]
	labelsFieldPath := field.NewPath("metadata").Child("labels", _ruleScopeConfig.ScopeLabelKey)
	if !ok {
		return field.NotFound(labelsFieldPath, rs)
	}

	switch rs {
	case _ruleScopeConfig.ScopeLabelValueCluster:
		if r.Namespace != _ruleScopeConfig.NamespaceScopeCluster {
			return field.Invalid(labelsFieldPath, rs, fmt.Sprintf(
				"cluster rule is only supported in namespace %s", _ruleScopeConfig.NamespaceScopeCluster))
		}
	case _ruleScopeConfig.ScopeLabelValueWorkspace:
		if r.Namespace != _ruleScopeConfig.NamespaceScopeWorkspace {
			return field.Invalid(labelsFieldPath, rs, fmt.Sprintf(
				"workspace rule is only supported in namespace %s", _ruleScopeConfig.NamespaceScopeWorkspace))
		}
		if r.Labels[rs] == "" {
			return field.Invalid(field.NewPath("metadata").Child("labels"), rs,
				fmt.Sprintf("workspace rule should be tagged a label called \"%s\"", _ruleScopeConfig.ScopeWorkspaceLabelKey))
		}
	case _ruleScopeConfig.ScopeLabelValueNamespace:
	default:
		return field.NotSupported(labelsFieldPath, rs, []string{
			_ruleScopeConfig.ScopeLabelValueCluster,
			_ruleScopeConfig.ScopeLabelValueWorkspace,
			_ruleScopeConfig.ScopeLabelValueNamespace})
	}

	return nil
}
