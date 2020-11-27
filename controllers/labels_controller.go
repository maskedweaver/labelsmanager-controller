/*


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

package controllers

import (
	"context"
	"reflect"
	"strings"
	"sync"

	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	labelsv1beta1 "github.com/nyambati/labelsmanager-controller/api/v1beta1"
)

// LabelsReconciler reconciles a Labels object
type LabelsReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=labels.rackbrains.io,resources=labels,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=labels.rackbrains.io,resources=labels/status,verbs=get;update;patch

func (r *LabelsReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("labels", req.NamespacedName)

	log.V(1).Info("Getting all deployment resources from the cluster")
	deployments := &apps.DeploymentList{}
	options := []client.ListOption{
		// client.MatchingFieldsSelector{Selector: fields.OneTermEqualSelector("spec.nodeName", nodeName)},
	}

	if err := r.List(ctx, deployments, options...); err != nil {
		log.Error(err, "Failed to list deployments")
		return ctrl.Result{}, err
	}

	instance := &labelsv1beta1.Labels{}

	log.V(1).Info("Getting resource from the cluster")

	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if len(instance.Spec.Labels) <= 0 {
		log.V(1).Info("resource is missing labels")
		return ctrl.Result{}, nil
	}

	numberOfDeployments := len(deployments.Items)

	var wg sync.WaitGroup
	wg.Add(numberOfDeployments)

	for i, deployment := range deployments.Items {
		go func(i int) {
			defer wg.Done()
			labels := mergeLabels(deployment.Labels, instance.Spec.Labels)
			deployment.SetLabels(labels)
			if err := r.Update(ctx, &deployment); err != nil {
				log.Error(err, err.Error())
			}
		}(i)
	}

	wg.Wait()
	return r.UpdateStatus(ctx, instance, numberOfDeployments)
}

func (r *LabelsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&labelsv1beta1.Labels{}).
		Watches(&source.Kind{Type: &apps.Deployment{}}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}

func (r *LabelsReconciler) UpdateStatus(ctx context.Context, instance *labelsv1beta1.Labels, deployments int) (ctrl.Result, error) {

	// instance.Status.Applied = "yes"
	instance.Status.ManagedDeployments = deployments
	instance.Status.Labels = getKeysFromMap(instance.Spec.Labels)

	if err := r.Status().Update(ctx, instance); err != nil {
		log.Error(err, "Unable to update label status")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func getKeysFromMap(m map[string]string) string {
	keys := reflect.ValueOf(m).MapKeys()
	strkeys := make([]string, 0, len(keys))
	for key := range m {
		strkeys = append(strkeys, key)
	}
	return strings.Join(strkeys, ",")
}

func mergeLabels(deploymentLabels, managedLabels map[string]string) map[string]string {

	if len(deploymentLabels) == 0 {
		return managedLabels
	}

	for key, value := range managedLabels {
		deploymentLabels[key] = value
	}

	return deploymentLabels
}
