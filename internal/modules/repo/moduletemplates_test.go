package repo

import (
	"context"
	"testing"

	"github.com/kyma-project/cli.v3/internal/kube/fake"
	"github.com/kyma-project/cli.v3/internal/kube/kyma"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestModuleTemplatesRepo_All(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		fakeKymaClient := fake.KymaClient{
			ReturnErr: nil,
			ReturnModuleTemplateList: kyma.ModuleTemplateList{
				Items: []kyma.ModuleTemplate{},
			},
		}
		fakeKubeClient := fake.KubeClient{
			TestKymaInterface: &fakeKymaClient,
		}

		repo := NewModuleTemplatesRepo(&fakeKubeClient)

		all, err := repo.All(context.Background())
		require.NoError(t, err)
		require.Len(t, all, 0)
	})

	t.Run("found", func(t *testing.T) {
		fakeKymaClient := fake.KymaClient{
			ReturnErr: nil,
			ReturnModuleTemplateList: kyma.ModuleTemplateList{
				Items: []kyma.ModuleTemplate{
					{
						Spec: kyma.ModuleTemplateSpec{
							ModuleName: "test-module",
						},
					},
				},
			},
		}
		fakeKubeClient := fake.KubeClient{
			TestKymaInterface: &fakeKymaClient,
		}

		repo := NewModuleTemplatesRepo(&fakeKubeClient)

		all, err := repo.All(context.Background())
		require.NoError(t, err)
		require.Len(t, all, 1)
		require.Equal(t, "test-module", all[0].Spec.ModuleName)
	})
}

func TestModuleTemplatesRepo_CommunityInstalledByName(t *testing.T) {
	fakeKymaClient := fake.KymaClient{
		ReturnErr: nil,
		ReturnModuleTemplateList: kyma.ModuleTemplateList{
			Items: []kyma.ModuleTemplate{{Spec: kyma.ModuleTemplateSpec{ModuleName: "foo"}}},
		},
	}
	fakeKubeClient := fake.KubeClient{
		TestKymaInterface: &fakeKymaClient,
	}

	repo := NewModuleTemplatesRepo(&fakeKubeClient)

	mods, err := repo.CommunityInstalledByName(context.Background(), "foo")
	require.NoError(t, err)
	require.Len(t, mods, 1)
	require.Equal(t, "foo", mods[0].Spec.ModuleName)
}

func TestModuleTemplatesRepo_CommunityByNameAndVersion_NotFound(t *testing.T) {
	fakeKymaClient := fake.KymaClient{
		ReturnErr: nil,
		ReturnModuleTemplateList: kyma.ModuleTemplateList{
			Items: []kyma.ModuleTemplate{
				{
					Spec: kyma.ModuleTemplateSpec{
						ModuleName: "foo",
						Version:    "v1",
					},
				},
			},
		},
	}
	fakeKubeClient := fake.KubeClient{
		TestKymaInterface: &fakeKymaClient,
	}

	repo := NewModuleTemplatesRepo(&fakeKubeClient)

	mod, err := repo.CommunityByNameAndVersion(context.Background(), "foo", "v2")
	require.NoError(t, err)
	require.Nil(t, mod)
}

func TestModuleTemplatesRepo_RunningAssociatedResourcesOfModule(t *testing.T) {
	fakeRootlessDynamicClient := fake.RootlessDynamicClient{
		ReturnListObjs: &unstructured.UnstructuredList{
			Items: []unstructured.Unstructured{
				{
					Object: map[string]any{
						"metadata": map[string]any{
							"name": "res1",
						},
					},
				},
			},
		},
	}

	fakeKymaClient := fake.KymaClient{}
	fakeKubeClient := fake.KubeClient{
		TestKymaInterface:            &fakeKymaClient,
		TestRootlessDynamicInterface: &fakeRootlessDynamicClient,
	}
	repo := NewModuleTemplatesRepo(&fakeKubeClient)

	mod := kyma.ModuleTemplate{
		Spec: kyma.ModuleTemplateSpec{
			AssociatedResources: []metav1.GroupVersionKind{{Group: "g", Version: "v1", Kind: "Kind"}},
		},
	}
	resources, err := repo.RunningAssociatedResourcesOfModule(context.Background(), mod)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	require.Equal(t, "res1", resources[0].GetName())
}
