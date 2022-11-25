package generator

import (
	"log"

	"github.com/siderolabs/talos/pkg/machinery/config"
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

func GenerateConfig(clusterName string, controlPlaneEndpoint string, ipAddress string) configBundle {
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
		log.Fatalf("failed to parse version contract: %s", err)
	}

	// generate the cluster-wide secrets once and use it for every node machine configuration
	// secrets can be stashed for future use by marshaling the structure to YAML or JSON
	secrets, err := generate.NewSecretsBundle(generate.NewClock(), generate.WithVersionContract(versionContract))
	if err != nil {
		log.Fatalf("failed to generate secrets bundle: %s", err)
	}

	input, err := generate.NewInput(clusterName, controlPlaneEndpoint, kubernetesVersion, secrets,
		generate.WithVersionContract(versionContract),
		// there are many more generate options available which allow to tweak generated config programmatically
	)
	if err != nil {
		log.Fatalf("failed to generate input: %s", err)
	}

	// Generate the controlplane config
	configbundle.ControlplaneConfig, err = generate.Config(machine.TypeControlPlane, input)
	if err != nil {
		log.Fatalf("failed to generate controlplane config: %s", err)
	}

	// Generate the worker config
	configbundle.WorkerConfig, err = generate.Config(machine.TypeWorker, input)
	if err != nil {
		log.Fatalf("failed to generate worker config: %s", err)
	}

	// generate the client Talos configuration (for API access, e.g. talosctl)
	clientCfg, err := generate.Talosconfig(input, generate.WithEndpointList([]string{ipAddress}))
	if err != nil {
		log.Fatalf("failed to generate client config: %s", err)
	}
	configbundle.TalosConfig, err = clientCfg.Bytes()
	if err != nil {
		log.Fatalf("failed to generate talos config %s", err)
	}

	return configbundle
}

func ApplyPatch(configbundle configBundle, configPatch []byte) ([]byte, []byte) {
	a := []byte("a")
	b := []byte("b")
	return a, b
	// config.
	// patch, err := config.configPatcher.Patch.LoadPatch(configPatch)
}
