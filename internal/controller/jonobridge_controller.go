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

// JonobridgeReconciler reconciles a Jonobridge object
type JonobridgeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=jonobridge.dwim.mx,resources=jonobridges,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=jonobridge.dwim.mx,resources=jonobridges/status,verbs=get;update;patch

// Reconcile manages the state of the Jonobridge custom resource
func (r *JonobridgeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Fetch Jonobridge instance
	var instance jonobridgev1alpha1.Jonobridge
	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		if errors.IsNotFound(err) {
			// Resource not found, could have been deleted after reconcile request
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Update status with a simulated client count
	instance.Status.ConnectedClients = r.simulateConnectedClients()
	if err := r.Status().Update(ctx, &instance); err != nil {
		log.Printf("Failed to update Jonobridge status: %v", err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// simulateConnectedClients simulates the number of connected clients
func (r *JonobridgeReconciler) simulateConnectedClients() int {
	// This function should simulate tracking connected clients for Jonobridge.
	return 10 // Replace with actual logic if available
}

// SetupWithManager registers this controller with the manager
func (r *JonobridgeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jonobridgev1alpha1.Jonobridge{}).
		Complete(r)
}
