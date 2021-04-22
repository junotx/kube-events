/*
Copyright 2020 The KubeSphere Authors.

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

package main

import (
	"flag"
	"os"

	"github.com/kubesphere/kube-events/pkg/config"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	eventsv1alpha1 "github.com/kubesphere/kube-events/pkg/apis/v1alpha1"
	"github.com/kubesphere/kube-events/pkg/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = eventsv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var cfg config.OperatorConfig
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&cfg.ConfigReloaderImage, "config-reloader-image", "jimmidyson/configmap-reload:v0.3.0", "Reload Image")
	flag.StringVar(&cfg.ConfigReloaderCPU, "config-reloader-cpu", "100m", "Config Reloader CPU. Value \"0\" disables it and causes no limit to be configured.")
	flag.StringVar(&cfg.ConfigReloaderMemory, "config-reloader-memory", "25Mi", "Config Reloader Memory. Value \"0\" disables it and causes no limit to be configured.")
	flag.Parse()

	{
		if ns, rsc := os.Getenv("POD_NAMESPACE"), eventsv1alpha1.DefaultRuleScopeConfig(); ns != "" && ns != rsc.NamespaceScopeCluster {
			rsc.NamespaceScopeCluster = ns
			rsc.NamespaceScopeWorkspace = ns
			eventsv1alpha1.SetRuleScopeConfig(rsc)
		}
	}

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "operator.events.kubesphere.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.RulerReconciler{
		Conf:   &cfg,
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Ruler"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Ruler")
		os.Exit(1)
	}
	if err = (&controllers.ExporterReconciler{
		Conf:   &cfg,
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Exporter"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Exporter")
		os.Exit(1)
	}
	if err = (&controllers.RuleReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Rule"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Rule")
		os.Exit(1)
	}
	if err = (&eventsv1alpha1.Rule{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "Rule")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
