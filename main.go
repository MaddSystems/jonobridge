package main

import (
	"log"
	"os"

	"bitbucket.org/maddisontest/jonobridge/internal/controller"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {
	// Set up the manager
	mgr, err := ctrl.NewManager(config.GetConfigOrDie(), ctrl.Options{
		Scheme:             manager.NewScheme(),
		LeaderElection:     false,
		LeaderElectionID:   "jonobridge-forwarder-controller",
		Port:               9443,
	})
	if err != nil {
		log.Fatalf("Unable to start manager: %v", err)
		os.Exit(1)
	}

	// Initialize Jonobridge controller
	if err = (&controller.JonobridgeReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		log.Fatalf("Unable to create Jonobridge controller: %v", err)
		os.Exit(1)
	}

	// Initialize Forwarder controller
	if err = (&controller.ForwarderReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		log.Fatalf("Unable to create Forwarder controller: %v", err)
		os.Exit(1)
	}

	// Start manager
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Fatalf("Manager stopped: %v", err)
		os.Exit(1)
	}
}
