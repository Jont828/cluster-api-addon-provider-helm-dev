/*
Copyright 2022 The Kubernetes Authors.

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
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/cluster-api/util"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var helmreleaseproxylog = logf.Log.WithName("helmreleaseproxy-resource")

func (r *HelmReleaseProxy) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-addons-cluster-x-k8s-io-v1alpha1-helmreleaseproxy,mutating=true,failurePolicy=fail,sideEffects=None,groups=addons.cluster.x-k8s.io,resources=helmreleaseproxies,verbs=create;update,versions=v1alpha1,name=helmreleaseproxy.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &HelmReleaseProxy{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (p *HelmReleaseProxy) Default() {
	helmreleaseproxylog.Info("default", "name", p.Name)

	if p.Spec.ReleaseName == "" {
		p.Spec.Options.Install.GenerateReleaseName = true
	}

	if p.Spec.ReleaseNamespace == "" {
		p.Spec.ReleaseNamespace = "default"
	}

	if p.Spec.Options.Install.GenerateReleaseName {
		p.Spec.ReleaseName = fmt.Sprintf("%s-%s", p.Spec.ChartName, util.RandomString(6))
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-addons-cluster-x-k8s-io-v1alpha1-helmreleaseproxy,mutating=false,failurePolicy=fail,sideEffects=None,groups=addons.cluster.x-k8s.io,resources=helmreleaseproxies,verbs=create;update,versions=v1alpha1,name=vhelmreleaseproxy.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &HelmReleaseProxy{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (p *HelmReleaseProxy) ValidateCreate() (admission.Warnings, error) {
	helmreleaseproxylog.Info("validate create", "name", p.Name)

	// They can both be set here because of the defaulting webhook.
	if p.Spec.ReleaseName == "" && !p.Spec.Options.Install.GenerateReleaseName {
		return nil, fmt.Errorf("spec.releaseName and spec.options.install.generateReleaseName cannot be empty/unset at the same time")
	}

	if err := isUrlValid(p.Spec.RepoURL); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *HelmReleaseProxy) ValidateUpdate(oldRaw runtime.Object) (admission.Warnings, error) {
	helmreleaseproxylog.Info("validate update", "name", r.Name)

	var allErrs field.ErrorList
	old, ok := oldRaw.(*HelmReleaseProxy)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a HelmReleaseProxy but got a %T", old))
	}

	if !reflect.DeepEqual(r.Spec.RepoURL, old.Spec.RepoURL) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "RepoURL"),
				r.Spec.RepoURL, "field is immutable"),
		)
	}

	if !reflect.DeepEqual(r.Spec.ChartName, old.Spec.ChartName) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "ChartName"),
				r.Spec.ChartName, "field is immutable"),
		)
	}

	if !reflect.DeepEqual(r.Spec.ReleaseName, old.Spec.ReleaseName) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "ReleaseName"),
				r.Spec.ReleaseName, "field is immutable"),
		)
	}

	if !reflect.DeepEqual(r.Spec.Options.Install.GenerateReleaseName, old.Spec.Options.Install.GenerateReleaseName) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "options", "install", "generateReleaseName"),
				r.Spec.Options.Install.GenerateReleaseName, "field is immutable"),
		)
	}

	if !reflect.DeepEqual(r.Spec.ReleaseNamespace, old.Spec.ReleaseNamespace) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "ReleaseNamespace"),
				r.Spec.ReleaseNamespace, "field is immutable"),
		)
	}

	if len(allErrs) > 0 {
		return nil, apierrors.NewInvalid(GroupVersion.WithKind("HelmReleaseProxy").GroupKind(), r.Name, allErrs)
	}

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *HelmReleaseProxy) ValidateDelete() (admission.Warnings, error) {
	helmreleaseproxylog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
