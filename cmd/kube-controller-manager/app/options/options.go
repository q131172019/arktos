/*
Copyright 2014 The Kubernetes Authors.
Copyright 2020 Authors of Arktos - file modified.

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

// Package options provides the flags used for the controller manager.
//
package options

import (
	"fmt"
	"net"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	apiserveroptions "k8s.io/apiserver/pkg/server/options"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	clientset "k8s.io/client-go/kubernetes"
	clientgokubescheme "k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/clientutil"
	cliflag "k8s.io/component-base/cli/flag"
	kubectrlmgrconfigv1alpha1 "k8s.io/kube-controller-manager/config/v1alpha1"
	cmoptions "k8s.io/kubernetes/cmd/controller-manager/app/options"
	"k8s.io/kubernetes/cmd/genutils"
	kubecontrollerconfig "k8s.io/kubernetes/cmd/kube-controller-manager/app/config"
	kubectrlmgrconfig "k8s.io/kubernetes/pkg/controller/apis/config"
	kubectrlmgrconfigscheme "k8s.io/kubernetes/pkg/controller/apis/config/scheme"
	"k8s.io/kubernetes/pkg/controller/garbagecollector"
	garbagecollectorconfig "k8s.io/kubernetes/pkg/controller/garbagecollector/config"
	"k8s.io/kubernetes/pkg/master/ports"

	// add the kubernetes feature gates
	_ "k8s.io/kubernetes/pkg/features"

	"k8s.io/klog"
)

const (
	// KubeControllerManagerUserAgent is the userAgent name when starting kube-controller managers.
	KubeControllerManagerUserAgent = "kube-controller-manager"
)

// KubeControllerManagerOptions is the main context object for the kube-controller manager.
type KubeControllerManagerOptions struct {
	Generic           *cmoptions.GenericControllerManagerConfigurationOptions
	KubeCloudShared   *cmoptions.KubeCloudSharedOptions
	ServiceController *cmoptions.ServiceControllerOptions

	AttachDetachController           *AttachDetachControllerOptions
	CSRSigningController             *CSRSigningControllerOptions
	DaemonSetController              *DaemonSetControllerOptions
	DeploymentController             *DeploymentControllerOptions
	DeprecatedFlags                  *DeprecatedControllerOptions
	EndpointController               *EndpointControllerOptions
	GarbageCollectorController       *GarbageCollectorControllerOptions
	HPAController                    *HPAControllerOptions
	JobController                    *JobControllerOptions
	NamespaceController              *NamespaceControllerOptions
	NodeIPAMController               *NodeIPAMControllerOptions
	NodeLifecycleController          *NodeLifecycleControllerOptions
	PersistentVolumeBinderController *PersistentVolumeBinderControllerOptions
	PodGCController                  *PodGCControllerOptions
	ReplicaSetController             *ReplicaSetControllerOptions
	ReplicationController            *ReplicationControllerOptions
	ResourceQuotaController          *ResourceQuotaControllerOptions
	SAController                     *SAControllerOptions
	TenantController                 *TenantControllerOptions
	TTLAfterFinishedController       *TTLAfterFinishedControllerOptions

	SecureServing *apiserveroptions.SecureServingOptionsWithLoopback
	// TODO: remove insecure serving mode
	InsecureServing *apiserveroptions.DeprecatedInsecureServingOptionsWithLoopback
	Authentication  *apiserveroptions.DelegatingAuthenticationOptions
	Authorization   *apiserveroptions.DelegatingAuthorizationOptions

	Master     string
	Kubeconfig string

	// optional resource provider kubeconfig
	ResourceProviderKubeConfig string
}

// NewKubeControllerManagerOptions creates a new KubeControllerManagerOptions with a default config.
func NewKubeControllerManagerOptions() (*KubeControllerManagerOptions, error) {
	componentConfig, err := NewDefaultComponentConfig(ports.InsecureKubeControllerManagerPort)
	if err != nil {
		return nil, err
	}

	s := KubeControllerManagerOptions{
		Generic:         cmoptions.NewGenericControllerManagerConfigurationOptions(&componentConfig.Generic),
		KubeCloudShared: cmoptions.NewKubeCloudSharedOptions(&componentConfig.KubeCloudShared),
		ServiceController: &cmoptions.ServiceControllerOptions{
			ServiceControllerConfiguration: &componentConfig.ServiceController,
		},
		AttachDetachController: &AttachDetachControllerOptions{
			&componentConfig.AttachDetachController,
		},
		CSRSigningController: &CSRSigningControllerOptions{
			&componentConfig.CSRSigningController,
		},
		DaemonSetController: &DaemonSetControllerOptions{
			&componentConfig.DaemonSetController,
		},
		DeploymentController: &DeploymentControllerOptions{
			&componentConfig.DeploymentController,
		},
		DeprecatedFlags: &DeprecatedControllerOptions{
			&componentConfig.DeprecatedController,
		},
		EndpointController: &EndpointControllerOptions{
			&componentConfig.EndpointController,
		},
		GarbageCollectorController: &GarbageCollectorControllerOptions{
			&componentConfig.GarbageCollectorController,
		},
		HPAController: &HPAControllerOptions{
			&componentConfig.HPAController,
		},
		JobController: &JobControllerOptions{
			&componentConfig.JobController,
		},
		NamespaceController: &NamespaceControllerOptions{
			&componentConfig.NamespaceController,
		},
		NodeIPAMController: &NodeIPAMControllerOptions{
			&componentConfig.NodeIPAMController,
		},
		NodeLifecycleController: &NodeLifecycleControllerOptions{
			&componentConfig.NodeLifecycleController,
		},
		PersistentVolumeBinderController: &PersistentVolumeBinderControllerOptions{
			&componentConfig.PersistentVolumeBinderController,
		},
		PodGCController: &PodGCControllerOptions{
			&componentConfig.PodGCController,
		},
		ReplicaSetController: &ReplicaSetControllerOptions{
			&componentConfig.ReplicaSetController,
		},
		ReplicationController: &ReplicationControllerOptions{
			&componentConfig.ReplicationController,
		},
		ResourceQuotaController: &ResourceQuotaControllerOptions{
			&componentConfig.ResourceQuotaController,
		},
		SAController: &SAControllerOptions{
			&componentConfig.SAController,
		},
		TenantController: &TenantControllerOptions{
			&componentConfig.TenantController,
		},
		TTLAfterFinishedController: &TTLAfterFinishedControllerOptions{
			&componentConfig.TTLAfterFinishedController,
		},
		SecureServing: apiserveroptions.NewSecureServingOptions().WithLoopback(),
		InsecureServing: (&apiserveroptions.DeprecatedInsecureServingOptions{
			BindAddress: net.ParseIP(componentConfig.Generic.Address),
			BindPort:    int(componentConfig.Generic.Port),
			BindNetwork: "tcp",
		}).WithLoopback(),
		Authentication: apiserveroptions.NewDelegatingAuthenticationOptions(),
		Authorization:  apiserveroptions.NewDelegatingAuthorizationOptions(),
	}

	s.Authentication.RemoteKubeConfigFileOptional = true
	s.Authorization.RemoteKubeConfigFileOptional = true
	s.Authorization.AlwaysAllowPaths = []string{"/healthz"}

	// Set the PairName but leave certificate directory blank to generate in-memory by default
	s.SecureServing.ServerCert.CertDirectory = ""
	s.SecureServing.ServerCert.PairName = "kube-controller-manager"
	s.SecureServing.BindPort = ports.KubeControllerManagerPort

	gcIgnoredResources := make([]garbagecollectorconfig.GroupResource, 0, len(garbagecollector.DefaultIgnoredResources()))
	for r := range garbagecollector.DefaultIgnoredResources() {
		gcIgnoredResources = append(gcIgnoredResources, garbagecollectorconfig.GroupResource{Group: r.Group, Resource: r.Resource})
	}

	s.GarbageCollectorController.GCIgnoredResources = gcIgnoredResources
	s.Generic.LeaderElection.ResourceName = "kube-controller-manager"
	s.Generic.LeaderElection.ResourceNamespace = "kube-system"

	return &s, nil
}

// NewDefaultComponentConfig returns kube-controller manager configuration object.
func NewDefaultComponentConfig(insecurePort int32) (kubectrlmgrconfig.KubeControllerManagerConfiguration, error) {
	versioned := kubectrlmgrconfigv1alpha1.KubeControllerManagerConfiguration{}
	kubectrlmgrconfigscheme.Scheme.Default(&versioned)

	internal := kubectrlmgrconfig.KubeControllerManagerConfiguration{}
	if err := kubectrlmgrconfigscheme.Scheme.Convert(&versioned, &internal, nil); err != nil {
		return internal, err
	}
	internal.Generic.Port = insecurePort
	return internal, nil
}

// Flags returns flags for a specific APIServer by section name
func (s *KubeControllerManagerOptions) Flags(allControllers []string, disabledByDefaultControllers []string) cliflag.NamedFlagSets {
	fss := cliflag.NamedFlagSets{}
	s.Generic.AddFlags(&fss, allControllers, disabledByDefaultControllers)
	s.KubeCloudShared.AddFlags(fss.FlagSet("generic"))
	s.ServiceController.AddFlags(fss.FlagSet("service controller"))

	s.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	s.InsecureServing.AddUnqualifiedFlags(fss.FlagSet("insecure serving"))
	s.Authentication.AddFlags(fss.FlagSet("authentication"))
	s.Authorization.AddFlags(fss.FlagSet("authorization"))

	s.AttachDetachController.AddFlags(fss.FlagSet("attachdetach controller"))
	s.CSRSigningController.AddFlags(fss.FlagSet("csrsigning controller"))
	s.DeploymentController.AddFlags(fss.FlagSet("deployment controller"))
	s.DaemonSetController.AddFlags(fss.FlagSet("daemonset controller"))
	s.DeprecatedFlags.AddFlags(fss.FlagSet("deprecated"))
	s.EndpointController.AddFlags(fss.FlagSet("endpoint controller"))
	s.GarbageCollectorController.AddFlags(fss.FlagSet("garbagecollector controller"))
	s.HPAController.AddFlags(fss.FlagSet("horizontalpodautoscaling controller"))
	s.JobController.AddFlags(fss.FlagSet("job controller"))
	s.NamespaceController.AddFlags(fss.FlagSet("namespace controller"))
	s.NodeIPAMController.AddFlags(fss.FlagSet("nodeipam controller"))
	s.NodeLifecycleController.AddFlags(fss.FlagSet("nodelifecycle controller"))
	s.PersistentVolumeBinderController.AddFlags(fss.FlagSet("persistentvolume-binder controller"))
	s.PodGCController.AddFlags(fss.FlagSet("podgc controller"))
	s.ReplicaSetController.AddFlags(fss.FlagSet("replicaset controller"))
	s.ReplicationController.AddFlags(fss.FlagSet("replicationcontroller"))
	s.ResourceQuotaController.AddFlags(fss.FlagSet("resourcequota controller"))
	s.SAController.AddFlags(fss.FlagSet("serviceaccount controller"))
	s.TenantController.AddFlags(fss.FlagSet("tenant controller"))
	s.TTLAfterFinishedController.AddFlags(fss.FlagSet("ttl-after-finished controller"))

	fs := fss.FlagSet("misc")
	fs.StringVar(&s.Master, "master", s.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig).")
	fs.StringVar(&s.Kubeconfig, "kubeconfig", s.Kubeconfig, "Path to kubeconfig files with authorization and master location information.")
	fs.StringVar(&s.ResourceProviderKubeConfig, "resource-providers", s.ResourceProviderKubeConfig, "Path to kubeconfig files points to resource provider(s).")
	utilfeature.DefaultMutableFeatureGate.AddFlag(fss.FlagSet("generic"))

	return fss
}

// ApplyTo fills up controller manager config with options.
func (s *KubeControllerManagerOptions) ApplyTo(c *kubecontrollerconfig.Config) error {
	if err := s.Generic.ApplyTo(&c.ComponentConfig.Generic); err != nil {
		return err
	}
	if err := s.KubeCloudShared.ApplyTo(&c.ComponentConfig.KubeCloudShared); err != nil {
		return err
	}
	if err := s.AttachDetachController.ApplyTo(&c.ComponentConfig.AttachDetachController); err != nil {
		return err
	}
	if err := s.CSRSigningController.ApplyTo(&c.ComponentConfig.CSRSigningController); err != nil {
		return err
	}
	if err := s.DaemonSetController.ApplyTo(&c.ComponentConfig.DaemonSetController); err != nil {
		return err
	}
	if err := s.DeploymentController.ApplyTo(&c.ComponentConfig.DeploymentController); err != nil {
		return err
	}
	if err := s.DeprecatedFlags.ApplyTo(&c.ComponentConfig.DeprecatedController); err != nil {
		return err
	}
	if err := s.EndpointController.ApplyTo(&c.ComponentConfig.EndpointController); err != nil {
		return err
	}
	if err := s.GarbageCollectorController.ApplyTo(&c.ComponentConfig.GarbageCollectorController); err != nil {
		return err
	}
	if err := s.HPAController.ApplyTo(&c.ComponentConfig.HPAController); err != nil {
		return err
	}
	if err := s.JobController.ApplyTo(&c.ComponentConfig.JobController); err != nil {
		return err
	}
	if err := s.NamespaceController.ApplyTo(&c.ComponentConfig.NamespaceController); err != nil {
		return err
	}
	if err := s.NodeIPAMController.ApplyTo(&c.ComponentConfig.NodeIPAMController); err != nil {
		return err
	}
	if err := s.NodeLifecycleController.ApplyTo(&c.ComponentConfig.NodeLifecycleController); err != nil {
		return err
	}
	if err := s.PersistentVolumeBinderController.ApplyTo(&c.ComponentConfig.PersistentVolumeBinderController); err != nil {
		return err
	}
	if err := s.PodGCController.ApplyTo(&c.ComponentConfig.PodGCController); err != nil {
		return err
	}
	if err := s.ReplicaSetController.ApplyTo(&c.ComponentConfig.ReplicaSetController); err != nil {
		return err
	}
	if err := s.ReplicationController.ApplyTo(&c.ComponentConfig.ReplicationController); err != nil {
		return err
	}
	if err := s.ResourceQuotaController.ApplyTo(&c.ComponentConfig.ResourceQuotaController); err != nil {
		return err
	}
	if err := s.SAController.ApplyTo(&c.ComponentConfig.SAController); err != nil {
		return err
	}
	if err := s.ServiceController.ApplyTo(&c.ComponentConfig.ServiceController); err != nil {
		return err
	}
	if err := s.TenantController.ApplyTo(&c.ComponentConfig.TenantController); err != nil {
		return err
	}
	if err := s.TTLAfterFinishedController.ApplyTo(&c.ComponentConfig.TTLAfterFinishedController); err != nil {
		return err
	}
	if err := s.InsecureServing.ApplyTo(&c.InsecureServing, &c.LoopbackClientConfig); err != nil {
		return err
	}
	if err := s.SecureServing.ApplyTo(&c.SecureServing, &c.LoopbackClientConfig); err != nil {
		return err
	}
	if s.SecureServing.BindPort != 0 || s.SecureServing.Listener != nil {
		if err := s.Authentication.ApplyTo(&c.Authentication, c.SecureServing, nil); err != nil {
			return err
		}
		if err := s.Authorization.ApplyTo(&c.Authorization); err != nil {
			return err
		}
	}

	// sync back to component config
	// TODO: find more elegant way than syncing back the values.
	c.ComponentConfig.Generic.Port = int32(s.InsecureServing.BindPort)
	c.ComponentConfig.Generic.Address = s.InsecureServing.BindAddress.String()

	return nil
}

// Validate is used to validate the options and config before launching the controller manager
func (s *KubeControllerManagerOptions) Validate(allControllers []string, disabledByDefaultControllers []string) error {
	var errs []error

	errs = append(errs, s.Generic.Validate(allControllers, disabledByDefaultControllers)...)
	errs = append(errs, s.KubeCloudShared.Validate()...)
	errs = append(errs, s.AttachDetachController.Validate()...)
	errs = append(errs, s.CSRSigningController.Validate()...)
	errs = append(errs, s.DaemonSetController.Validate()...)
	errs = append(errs, s.DeploymentController.Validate()...)
	errs = append(errs, s.DeprecatedFlags.Validate()...)
	errs = append(errs, s.EndpointController.Validate()...)
	errs = append(errs, s.GarbageCollectorController.Validate()...)
	errs = append(errs, s.HPAController.Validate()...)
	errs = append(errs, s.JobController.Validate()...)
	errs = append(errs, s.NamespaceController.Validate()...)
	errs = append(errs, s.NodeIPAMController.Validate()...)
	errs = append(errs, s.NodeLifecycleController.Validate()...)
	errs = append(errs, s.PersistentVolumeBinderController.Validate()...)
	errs = append(errs, s.PodGCController.Validate()...)
	errs = append(errs, s.ReplicaSetController.Validate()...)
	errs = append(errs, s.ReplicationController.Validate()...)
	errs = append(errs, s.ResourceQuotaController.Validate()...)
	errs = append(errs, s.SAController.Validate()...)
	errs = append(errs, s.ServiceController.Validate()...)
	errs = append(errs, s.TenantController.Validate()...)
	errs = append(errs, s.TTLAfterFinishedController.Validate()...)
	errs = append(errs, s.SecureServing.Validate()...)
	errs = append(errs, s.InsecureServing.Validate()...)
	errs = append(errs, s.Authentication.Validate()...)
	errs = append(errs, s.Authorization.Validate()...)

	// TODO: validate component config, master and kubeconfig

	return utilerrors.NewAggregate(errs)
}

// Config return a controller manager config objective
func (s KubeControllerManagerOptions) Config(allControllers []string, disabledByDefaultControllers []string) (*kubecontrollerconfig.Config, error) {
	if err := s.Validate(allControllers, disabledByDefaultControllers); err != nil {
		return nil, err
	}

	if err := s.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	kubeconfigs, err := clientcmd.BuildConfigFromFlags(s.Master, s.Kubeconfig)
	if err != nil {
		return nil, err
	}

	for _, kubeConfig := range kubeconfigs.GetAllConfigs() {
		kubeConfig.ContentConfig.ContentType = s.Generic.ClientConnection.ContentType
		kubeConfig.QPS = s.Generic.ClientConnection.QPS
		kubeConfig.Burst = int(s.Generic.ClientConnection.Burst)
	}

	client, err := clientset.NewForConfig(restclient.AddUserAgent(kubeconfigs, KubeControllerManagerUserAgent))
	if err != nil {
		return nil, err
	}

	// shallow copy, do not modify the kubeconfig.Timeout.
	config := *kubeconfigs.GetConfig()
	config.Timeout = s.Generic.LeaderElection.RenewDeadline.Duration
	leaderElectionClient := clientset.NewForConfigOrDie(restclient.AddUserAgent(restclient.NewAggregatedConfig(&config), "leader-election"))

	eventRecorder := createRecorder(client, KubeControllerManagerUserAgent)

	// get resource provider kube configs
	var resourceProviderKubeconfigs []*restclient.Config
	if len(s.ResourceProviderKubeConfig) > 0 {
		resourceProviderKubeConfigFiles, existed := genutils.ParseKubeConfigFiles(s.ResourceProviderKubeConfig)
		// TODO: rewrite the IF block when perf test env no longer requires sequential setup of TP/RP clusters
		if !existed {
			klog.Warningf("--resource-providers points to non existed file(s), default to local cluster kubeconfig file")
			resourceProviderKubeconfigs = make([]*restclient.Config, 1)
			resourceProviderKubeconfigs[0] = restclient.CopyConfigs(kubeconfigs)
		} else {
			resourceProviderKubeconfigs = make([]*restclient.Config, len(resourceProviderKubeConfigFiles))
			for i, kubeConfigFile := range resourceProviderKubeConfigFiles {
				resourceProviderKubeconfig, err := clientutil.CreateClientConfigFromKubeconfigFileAndSetQpsNoUserAgent(kubeConfigFile, 0, 0, "")
				if err != nil {
					return nil, fmt.Errorf("failed to create resource provider rest client from kubeconfig [%s], error [%v]", kubeConfigFile, err)
				}
				resourceProviderKubeconfigs[i] = resourceProviderKubeconfig
				klog.V(3).Infof("Created resource provider kubeconfig %d %p", i, resourceProviderKubeconfig)
			}
		}
	}

	c := &kubecontrollerconfig.Config{
		Client:                      client,
		Kubeconfig:                  kubeconfigs,
		EventRecorder:               eventRecorder,
		LeaderElectionClient:        leaderElectionClient,
		ResourceProviderKubeconfigs: resourceProviderKubeconfigs,
	}
	if err := s.ApplyTo(c); err != nil {
		return nil, err
	}

	return c, nil
}

func createRecorder(kubeClient clientset.Interface, userAgent string) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: kubeClient.CoreV1().EventsWithMultiTenancy(metav1.NamespaceAll, metav1.TenantAll)})
	return eventBroadcaster.NewRecorder(clientgokubescheme.Scheme, v1.EventSource{Component: userAgent})
}
