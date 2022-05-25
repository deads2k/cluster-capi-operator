package cluster

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/openshift/cluster-capi-operator/pkg/controllers"
	"github.com/openshift/cluster-capi-operator/pkg/operatorstatus"
	"github.com/openshift/cluster-capi-operator/pkg/test"
)

var _ = Describe("Reconcile Core cluster", func() {
	var r *CoreClusterReconciler
	var coreCluster *clusterv1.Cluster

	BeforeEach(func() {
		r = &CoreClusterReconciler{
			ClusterOperatorStatusClient: operatorstatus.ClusterOperatorStatusClient{
				Client: cl,
			},
			Cluster: &clusterv1.Cluster{},
		}

		coreCluster = &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-name",
				Namespace: controllers.DefaultManagedNamespace,
			},
		}

		Expect(cl.Create(ctx, coreCluster)).To(Succeed())
	})

	AfterEach(func() {
		Expect(test.CleanupAndWait(ctx, cl, coreCluster)).To(Succeed())
	})

	It("should update core cluster status", func() {
		_, err := r.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: coreCluster.Namespace,
				Name:      coreCluster.Name,
			},
		})
		Expect(err).ToNot(HaveOccurred())

		Expect(cl.Get(ctx, client.ObjectKey{
			Name:      coreCluster.Name,
			Namespace: coreCluster.Namespace,
		}, coreCluster)).To(Succeed())

		Expect(coreCluster.Status.Conditions).ToNot(BeEmpty())
		Expect(coreCluster.Status.Conditions[0].Type).To(Equal(clusterv1.ControlPlaneInitializedCondition))
		Expect(coreCluster.Status.Conditions[0].Status).To(Equal(corev1.ConditionTrue))
	})
})
