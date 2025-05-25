package indexer

import (
	"context"
	tritonappsv1alpha1 "github.com/triton-io/triton/apis/apps/v1alpha1"
	"github.com/triton-io/triton/pkg/log"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func AddDefaultIndexers(mgr manager.Manager) error {
	err := mgr.GetFieldIndexer().IndexField(
		context.TODO(),
		&tritonappsv1alpha1.DeployFlow{},
		"spec.Application.AppName",
		func(rawObj runtime.Object) []string {
			flow := rawObj.(*tritonappsv1alpha1.DeployFlow)
			return []string{flow.Spec.Application.AppName}
		},
	)
	if err != nil {
		return err
	}
	log.Info("Added indexer for DeployFlow")
	return nil
}
