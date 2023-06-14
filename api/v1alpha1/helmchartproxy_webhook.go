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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var helmchartproxylog = logf.Log.WithName("helmchartproxy-resource")

func (r *HelmChartProxy) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-addons-cluster-x-k8s-io-v1alpha1-helmchartproxy,mutating=true,failurePolicy=fail,sideEffects=None,groups=addons.cluster.x-k8s.io,resources=helmchartproxies,verbs=create;update,versions=v1alpha1,name=helmchartproxy.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &HelmChartProxy{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (p *HelmChartProxy) Default() {
	helmchartproxylog.Info("default", "name", p.Name)

	if p.Spec.ReleaseNamespace == "" {
		p.Spec.ReleaseNamespace = "default"
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-addons-cluster-x-k8s-io-v1alpha1-helmchartproxy,mutating=false,failurePolicy=fail,sideEffects=None,groups=addons.cluster.x-k8s.io,resources=helmchartproxies,verbs=create;update,versions=v1alpha1,name=vhelmchartproxy.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &HelmChartProxy{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (c *HelmChartProxy) ValidateCreate() (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (c *HelmChartProxy) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	helmchartproxylog.Info("validate update", "name", c.Name)

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (c *HelmChartProxy) ValidateDelete() (admission.Warnings, error) {
	helmchartproxylog.Info("validate delete", "name", c.Name)

	return nil, nil
}
