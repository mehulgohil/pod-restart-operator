/*
Copyright 2025.

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

package controller

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplev1 "github.com/mehulgohil/pod-restart-operator/api/v1"
)

// PodRestartReconciler reconciles a PodRestart object
type PodRestartReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.example.com,resources=podrestarts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.example.com,resources=podrestarts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.example.com,resources=podrestarts/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodRestart object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *PodRestartReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling PodRestart")

	// Fetch the PodRestart resource
	podRestart := &examplev1.PodRestart{}
	err := r.Get(ctx, req.NamespacedName, podRestart)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Parse restart interval
	restartInterval, err := time.ParseDuration(podRestart.Spec.RestartInterval)
	if err != nil {
		log.Error(err, "Invalid restart interval format")
		return reconcile.Result{}, nil
	}

	// List all pods matching the LabelSelector
	podList := &v1.PodList{}
	labelSelector := labels.SelectorFromSet(podRestart.Spec.LabelSelector)
	listOpts := &client.ListOptions{LabelSelector: labelSelector}
	err = r.List(ctx, podList, listOpts)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info("Found matching pods", "count", len(podList.Items))

	// Restart each matching pod
	for _, pod := range podList.Items {
		log.Info("Restarting pod:", "pod", pod.Name)
		err := r.Delete(ctx, &pod)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	// Update status
	podRestart.Status.LastRestart = metav1.Now()
	err = r.Status().Update(ctx, podRestart)
	if err != nil {
		log.Error(err, "Failed to update PodRestart status")
	}

	// Requeue to check again after the interval
	return reconcile.Result{RequeueAfter: restartInterval}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodRestartReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.PodRestart{}).
		Complete(r)
}
