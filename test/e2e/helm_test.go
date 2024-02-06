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
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	addonsv1alpha1 "sigs.k8s.io/cluster-api-addon-provider-helm/api/v1alpha1"
	capi_e2e "sigs.k8s.io/cluster-api/test/e2e"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util"
)

var _ = Describe("Workload cluster creation", func() {
	var (
		ctx       = context.TODO()
		specName  = "create-workload-cluster"
		namespace *corev1.Namespace
		// cancelWatches     context.CancelFunc
		result            *clusterctl.ApplyClusterTemplateAndWaitResult
		clusterName       string
		clusterNamePrefix string
		// additionalCleanup func()
		specTimes = map[string]time.Time{}
	)

	BeforeEach(func() {
		logCheckpoint(specTimes)

		Expect(ctx).NotTo(BeNil(), "ctx is required for %s spec", specName)
		Expect(e2eConfig).NotTo(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(clusterctlConfigPath).To(BeAnExistingFile(), "Invalid argument. clusterctlConfigPath must be an existing file when calling %s spec", specName)
		Expect(bootstrapClusterProxy).NotTo(BeNil(), "Invalid argument. bootstrapClusterProxy can't be nil when calling %s spec", specName)
		Expect(os.MkdirAll(artifactFolder, 0o755)).To(Succeed(), "Invalid argument. artifactFolder can't be created for %s spec", specName)
		Expect(e2eConfig.Variables).To(HaveKey(capi_e2e.KubernetesVersion))

		// CLUSTER_NAME and CLUSTER_NAMESPACE allows for testing existing clusters
		// if CLUSTER_NAMESPACE is set don't generate a new prefix otherwise
		// the correct namespace won't be found and a new cluster will be created
		clusterNameSpace := os.Getenv("CLUSTER_NAMESPACE")
		if clusterNameSpace == "" {
			clusterNamePrefix = fmt.Sprintf("caaph-e2e-%s", util.RandomString(6))
		} else {
			clusterNamePrefix = clusterNameSpace
		}

		// Setup a Namespace where to host objects for this spec and create a watcher for the namespace events.
		var err error
		namespace, _, err = setupSpecNamespace(ctx, clusterNamePrefix, bootstrapClusterProxy, artifactFolder)
		// cancelWatches is the blank return arg above.
		Expect(err).NotTo(HaveOccurred())

		result = new(clusterctl.ApplyClusterTemplateAndWaitResult)
	})

	Context("Creating workload cluster [REQUIRED]", func() {
		It("With default template and calico Helm chart", func() {
			clusterName = fmt.Sprintf("%s-%s", specName, util.RandomString(6))
			clusterctl.ApplyClusterTemplateAndWait(ctx, createApplyClusterTemplateInput(
				specName,
				withNamespace(namespace.Name),
				withClusterName(clusterName),
				withControlPlaneMachineCount(1),
				withWorkerMachineCount(1),
				withControlPlaneWaiters(clusterctl.ControlPlaneWaiters{
					WaitForControlPlaneInitialized: EnsureControlPlaneInitialized,
				}),
			), result)

			// Create new Helm chart
			By("Creating new HelmChartProxy to install nginx", func() {
				hcp := &addonsv1alpha1.HelmChartProxy{
					TypeMeta: metav1.TypeMeta{
						APIVersion: addonsv1alpha1.GroupVersion.String(),
						Kind:       "HelmChartProxy",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx-ingress",
						Namespace: namespace.Name,
					},
					Spec: addonsv1alpha1.HelmChartProxySpec{
						ClusterSelector: metav1.LabelSelector{
							MatchLabels: map[string]string{
								"nginxIngress": "enabled",
							},
						},
						ReleaseName: "nginx-ingress",
						ChartName:   "nginx-ingress",
						RepoURL:     "https://helm.nginx.com/stable",
						ValuesTemplate: `controller:
            name: "{{ .ControlPlane.metadata.name }}-nginx"
              nginxStatus:
                allowCidrs: 127.0.0.1,::1,{{ index .Cluster.spec.clusterNetwork.pods.cidrBlocks 0 }}`,
						Options: &addonsv1alpha1.HelmOptions{},
					},
				}
				HelmInstallSpec(ctx, func() HelmInstallInput {
					return HelmInstallInput{
						BootstrapClusterProxy: bootstrapClusterProxy,
						Namespace:             namespace,
						ClusterName:           clusterName,
						HelmChartProxy:        hcp,
					}
				})
			})

			// Try to update Helm values and expect revision

			// Add workload cluster
		})
	})
})
