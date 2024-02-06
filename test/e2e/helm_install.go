//go:build e2e
// +build e2e

/*
Copyright 2024 The Kubernetes Authors.

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

package e2e

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	helmAction "helm.sh/helm/v3/pkg/action"
	helmRelease "helm.sh/helm/v3/pkg/release"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	addonsv1alpha1 "sigs.k8s.io/cluster-api-addon-provider-helm/api/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api/test/framework"
)

// HelmInstallInput ...
type HelmInstallInput struct {
	BootstrapClusterProxy framework.ClusterProxy
	Namespace             *corev1.Namespace
	ClusterName           string
	HelmChartProxy        *addonsv1alpha1.HelmChartProxy
}

func HelmInstallSpec(ctx context.Context, inputGetter func() HelmInstallInput) {
	var (
		specName             = "helm-install"
		input                HelmInstallInput
		workloadClusterProxy framework.ClusterProxy
		workloadClient       ctrlclient.Client
		mgmtClient           ctrlclient.Client
	)

	input = inputGetter()
	Expect(input.BootstrapClusterProxy).NotTo(BeNil(), "Invalid argument. input.BootstrapClusterProxy can't be nil when calling %s spec", specName)
	Expect(input.Namespace).NotTo(BeNil(), "Invalid argument. input.Namespace can't be nil when calling %s spec", specName)

	By("creating a Kubernetes client to the workload cluster")
	workloadClusterProxy = input.BootstrapClusterProxy.GetWorkloadCluster(ctx, input.Namespace.Name, input.ClusterName)
	Expect(workloadClusterProxy).NotTo(BeNil())

	workloadClient = workloadClusterProxy.GetClient()
	Expect(workloadClient).NotTo(BeNil())

	mgmtClient = input.BootstrapClusterProxy.GetClient()
	Expect(mgmtClient).NotTo(BeNil())

	// Create HCP
	Byf("Creating HelmChartProxy %s/%s", input.HelmChartProxy.Namespace, input.HelmChartProxy.Name)
	err := workloadClient.Create(ctx, input.HelmChartProxy)
	Expect(err).NotTo(HaveOccurred())

	// Get Cluster
	cluster := &clusterv1.Cluster{}
	key := types.NamespacedName{
		Namespace: input.Namespace.Name,
		Name:      input.ClusterName,
	}
	err = workloadClient.Get(ctx, key, cluster)
	Expect(err).NotTo(HaveOccurred())

	// Patch cluster labels, ignore match expressions for now
	selector := input.HelmChartProxy.Spec.ClusterSelector
	labels := cluster.Labels
	if labels == nil {
		labels = make(map[string]string)
	}

	for k, v := range selector.MatchLabels {
		labels[k] = v
	}

	err = workloadClient.Update(ctx, cluster)
	Expect(err).NotTo(HaveOccurred())

	// Get Cluster proxy
	workloadKubeconfigPath := workloadClusterProxy.GetKubeconfigPath()

	cliConfig := genericclioptions.NewConfigFlags(false)
	cliConfig.KubeConfig = &workloadKubeconfigPath

	actionConfig := new(helmAction.Configuration)
	err = actionConfig.Init(cliConfig, input.Namespace.Name, "secret", Logf)
	Expect(err).NotTo(HaveOccurred())

	// Watch for Helm release
	start := time.Now()

	// Workaround atm so we don't need to deal with generated random Helm release names
	releaseName := input.HelmChartProxy.Spec.ReleaseName
	Expect(releaseName).NotTo(BeEmpty())

	Log("starting to wait for Helm release to become available and deployed")
	Eventually(func() bool {
		getClient := helmAction.NewGet(actionConfig)

		if release, err := getClient.Run(releaseName); err == nil {
			if release != nil && release.Info.Status == helmRelease.StatusDeployed {
				return true
			}
		}

		return false
	}, e2eConfig.GetIntervals(specName, "wait-helm-release")...).Should(BeTrue(), "Failed to get Helm release %s on cluster %s", releaseName, input.ClusterName)
	Logf("Helm release %s is now available, took %v", releaseName, time.Since(start))

}
