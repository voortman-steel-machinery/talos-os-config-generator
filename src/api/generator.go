package main

import (
	"log"

	"github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/configpatcher"
	v1alpha1 "github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1/generate"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1/machine"
	"github.com/siderolabs/talos/pkg/machinery/constants"
)

type configBundle struct {
	ControlplaneConfig *v1alpha1.Config
	WorkerConfig       *v1alpha1.Config
	TalosConfig        []byte
}

func GenerateConfig(clusterName string, controlPlaneEndpoint string, ipAddress string) (configBundle, error) {
	configbundle := configBundle{}

	// * Kubernetes version to install, using the latest here
	kubernetesVersion := constants.DefaultKubernetesVersion

	// * version contract defines the version of the Talos cluster configuration is generated for
	//   generate package can generate machine configuration compatible with current and previous versions of Talos
	targetVersion := "v1.0"

	// parse the version contract
	var (
		versionContract = config.TalosVersionCurrent //nolint:wastedassign,ineffassign // version of the Talos machinery package
		err             error
	)

	versionContract, err = config.ParseContractFromVersion(targetVersion)
	if err != nil {
		log.Println("failed to parse version contract: ", err)
		return configBundle{}, err
	}

	// generate the cluster-wide secrets once and use it for every node machine configuration
	// secrets can be stashed for future use by marshaling the structure to YAML or JSON
	secrets, err := generate.NewSecretsBundle(generate.NewClock(), generate.WithVersionContract(versionContract))
	if err != nil {
		log.Println("failed to generate secrets bundle: ", err)
		return configBundle{}, err
	}

	input, err := generate.NewInput(clusterName, controlPlaneEndpoint, kubernetesVersion, secrets,
		generate.WithVersionContract(versionContract),
		// there are many more generate options available which allow to tweak generated config programmatically
	)
	if err != nil {
		log.Println("failed to generate input: ", err)
		return configBundle{}, err
	}

	// Generate the controlplane config
	configbundle.ControlplaneConfig, err = generate.Config(machine.TypeControlPlane, input)
	if err != nil {
		log.Println("failed to generate controlplane config: ", err)
		return configBundle{}, err
	}

	// Generate the worker config
	configbundle.WorkerConfig, err = generate.Config(machine.TypeWorker, input)
	if err != nil {
		log.Println("failed to generate worker config: ", err)
		return configBundle{}, err
	}

	// generate the client Talos configuration (for API access, e.g. talosctl)
	clientCfg, err := generate.Talosconfig(input, generate.WithEndpointList([]string{ipAddress}))
	if err != nil {
		log.Println("failed to generate client config: ", err)
		return configBundle{}, err
	}
	configbundle.TalosConfig, err = clientCfg.Bytes()
	if err != nil {
		log.Println("failed to generate talos config ", err)
		return configBundle{}, err
	}

	return configbundle, nil
}

func ApplyPatch(configbundle configBundle, configPatch []byte) ([]byte, []byte, error) {
	patch, err := configpatcher.LoadPatch(configPatch)
	if err != nil {
		log.Println("Cannot create patch: ", err)
		return nil, nil, err
	}

	controlplane := _patchAndMarshal(configbundle.ControlplaneConfig, []configpatcher.Patch{patch})
	worker := _patchAndMarshal(configbundle.WorkerConfig, []configpatcher.Patch{patch})

	return controlplane, worker, nil
	// config.
	// patch, err := config.configPatcher.Patch.LoadPatch(configPatch)
}

func _patchAndMarshal(config config.Provider, patches []configpatcher.Patch) []byte {
	cpcfg := configpatcher.WithConfig(config)
	cpPatched, err := configpatcher.Apply(cpcfg, patches)
	if err != nil {
		log.Println("Cannot apply patch: ", err)
	}
	marshal, err := cpPatched.Bytes()
	if err != nil {
		log.Println("failed to marshall config ", err)
	}
	return marshal
}
