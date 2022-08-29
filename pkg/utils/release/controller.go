package release

import (
	"context"
	"github.com/redhat-appstudio/release-service/kcp"

	kubeCl "github.com/redhat-appstudio/e2e-tests/pkg/apis/kubernetes"
	gitopsv1alpha1 "github.com/redhat-appstudio/managed-gitops/appstudio-shared/apis/appstudio.redhat.com/v1alpha1"
	"github.com/redhat-appstudio/release-service/api/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SuiteController struct {
	*kubeCl.K8sClient
}

func NewSuiteController(kube *kubeCl.K8sClient) (*SuiteController, error) {
	return &SuiteController{
		kube,
	}, nil
}

// CreateApplicationSnapshot creates a new ApplicationSnapshot using the given parameters.
func (s *SuiteController) CreateApplicationSnapshot(name string, namespace string, applicationName string, snapshotComponents []gitopsv1alpha1.ApplicationSnapshotComponent) (*gitopsv1alpha1.ApplicationSnapshot, error) {
	applicationSnapshot := &gitopsv1alpha1.ApplicationSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: gitopsv1alpha1.ApplicationSnapshotSpec{
			Application: applicationName,
			Components:  snapshotComponents,
		},
	}

	return applicationSnapshot, s.KubeRest().Create(context.TODO(), applicationSnapshot)
}

// CreateRelease creates a new Release using the given parameters.
func (s *SuiteController) CreateRelease(name, namespace, snapshot, releasePlan string) (*v1alpha1.Release, error) {
	release := &v1alpha1.Release{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.ReleaseSpec{
			ApplicationSnapshot: snapshot,
			ReleasePlan:         releasePlan,
		},
	}

	return release, s.KubeRest().Create(context.TODO(), release)
}

// CreateReleasePlan creates a new ReleasePlan using the given parameters.
func (s *SuiteController) CreateReleasePlan(name, namespace, application, targetNamespace string) (*v1alpha1.ReleasePlan, error) {
	releasePlan := &v1alpha1.ReleasePlan{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.ReleasePlanSpec{
			DisplayName: name,
			Application: application,
			Target: kcp.NamespaceReference{
				Namespace: targetNamespace,
			},
		},
	}

	return releasePlan, s.KubeRest().Create(context.TODO(), releasePlan)
}

// CreateReleasePlanAdmission creates a new ReleasePlan using the given parameters.
func (s *SuiteController) CreateReleasePlanAdmission(name, namespace, application, originNamespace, environment, releaseStrategy string) (*v1alpha1.ReleasePlanAdmission, error) {
	releasePlanAdmission := &v1alpha1.ReleasePlanAdmission{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.ReleasePlanAdmissionSpec{
			DisplayName: name,
			Application: application,
			Origin: kcp.NamespaceReference{
				Namespace: originNamespace,
			},
			Environment:     environment,
			ReleaseStrategy: releaseStrategy,
		},
	}

	return releasePlanAdmission, s.KubeRest().Create(context.TODO(), releasePlanAdmission)
}

// CreateReleaseStrategy creates a new ReleaseStrategy using the given parameters.
func (s *SuiteController) CreateReleaseStrategy(name, namespace, pipelineName, bundle string, policy string) (*v1alpha1.ReleaseStrategy, error) {
	releaseStrategy := &v1alpha1.ReleaseStrategy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.ReleaseStrategySpec{
			Pipeline: pipelineName,
			Bundle:   bundle,
			Policy:   policy,
		},
	}

	return releaseStrategy, s.KubeRest().Create(context.TODO(), releaseStrategy)
}

// GetPipelineRunInNamespace returns the Release PipelineRun referencing the given release.
func (s *SuiteController) GetPipelineRunInNamespace(namespace, releaseName, releaseNamespace string) (*v1beta1.PipelineRun, error) {
	pipelineRuns := &v1beta1.PipelineRunList{}
	opts := []client.ListOption{
		client.MatchingLabels{
			"release.appstudio.openshift.io/name":      releaseName,
			"release.appstudio.openshift.io/namespace": releaseNamespace,
		},
		client.InNamespace(namespace),
	}

	err := s.KubeRest().List(context.TODO(), pipelineRuns, opts...)

	if err == nil && len(pipelineRuns.Items) > 0 {
		return &pipelineRuns.Items[0], nil
	}

	return nil, err
}

// GetRelease returns the release with the given name in the given namespace.
func (s *SuiteController) GetRelease(releaseName, releaseNamespace string) (*v1alpha1.Release, error) {
	release := &v1alpha1.Release{}

	err := s.KubeRest().Get(context.TODO(), types.NamespacedName{
		Name:      releaseName,
		Namespace: releaseNamespace,
	}, release)

	return release, err
}

// GetReleasePlan returns a ReleasePlan from a given namespace.
func (s *SuiteController) GetReleasePlan(name string, namespace string) (*v1alpha1.ReleasePlan, error) {
	namespacedName := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}

	releasePlan := &v1alpha1.ReleasePlan{}
	err := s.KubeRest().Get(context.TODO(), namespacedName, releasePlan)
	if err != nil {
		return nil, err
	}
	return releasePlan, nil
}

// GetReleasePlanAdmission returns a ReleasePlan from a given namespace.
func (s *SuiteController) GetReleasePlanAdmission(name string, namespace string) (*v1alpha1.ReleasePlanAdmission, error) {
	namespacedName := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}

	releasePlanAdmission := &v1alpha1.ReleasePlanAdmission{}
	err := s.KubeRest().Get(context.TODO(), namespacedName, releasePlanAdmission)
	if err != nil {
		return nil, err
	}
	return releasePlanAdmission, nil
}
