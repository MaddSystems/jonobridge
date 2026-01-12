package controller

import (
	"context"
	"fmt"
	"log"

	jonobridgev1alpha1 "bitbucket.org/maddisontest/jonobridge/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ForwarderReconciler reconciles a Forwarder object
type ForwarderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=jonobridge.dwim.mx,resources=forwarders,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=jonobridge.dwim.mx,resources=forwarders/status,verbs=get;update;patch

// Reconcile manages the state of the Forwarder custom resource
func (r *ForwarderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Fetch Forwarder instance
	var instance jonobridgev1alpha1.Forwarder
	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		if errors.IsNotFound(err) {
			// Resource not found, could have been deleted after reconcile request
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Update status with a simulated message count
	instance.Status.MessageCount = r.simulateMessageCount()
	if err := r.Status().Update(ctx, &instance); err != nil {
		log.Printf("Failed to update Forwarder status: %v", err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// simulateMessageCount simulates the message count for Forwarder
func (r *ForwarderReconciler) simulateMessageCount() int {
	// This function should simulate counting messages for Forwarder.
	return 20 // Replace with actual logic if available
}

// SetupWithManager registers this controller with the manager
func (r *ForwarderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jonobridgev1alpha1.Forwarder{}).
		Complete(r)
}
