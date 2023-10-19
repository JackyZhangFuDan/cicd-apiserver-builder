/*
Copyright 2023.

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

package cicd

import (
	"context"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cicdv1 "github.com/cicd-apiserver-builder/pkg/apis/cicd/v1"
)

// JenkinsServiceReconciler reconciles a JenkinsService object
type JenkinsServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cicd,resources=jenkinsservices,verbs=create;delete
//+kubebuilder:rbac:groups=cicd,resources=jenkinsservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cicd,resources=jenkinsservices/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the JenkinsService object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *JenkinsServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("reconcile is triggered. ")

	js := &cicdv1.JenkinsService{}
	err := r.Get(ctx, req.NamespacedName, js)
	// 失败了直接返回，垃圾收集会清理js
	if errors.IsNotFound(err) {
		logger.Info("The Jenkins Service has been deleted")
		return ctrl.Result{}, nil
	}
	// 有错误会使得js再入工作队列，不用我们关心
	if err != nil {
		logger.Error(err, "Get jenkins service failed")
		return ctrl.Result{}, err
	}

	var replicas int32
	replicas = int32(js.Spec.InstanceAmount)
	selector := map[string]string{}
	selector["type"] = "jenkinsservice"
	selector["jsname"] = js.Name

	// 创建一个Deployment
	d := apps.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			UID:         uuid.NewUUID(),
			Name:        js.Name,
			Namespace:   js.Namespace,
			Annotations: make(map[string]string),
		},
		Spec: apps.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: selector},
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: selector,
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            "jenkinsforjs",
							Image:           "nginx:latest",
							ImagePullPolicy: "IfNotPresent",
						},
					},
				},
			},
		},
	}
	err = r.Get(ctx, types.NamespacedName{Namespace: js.Namespace, Name: js.Name}, &apps.Deployment{})
	if err == nil {
		logger.Info("Deployment with the same name is existed already")
		return ctrl.Result{}, nil
	}

	err = r.Create(ctx, &d, &client.CreateOptions{})
	if err != nil {
		logger.Error(err, "create deployment failed")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JenkinsServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cicdv1.JenkinsService{}).
		Complete(r)
}
