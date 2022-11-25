package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/voortman-steel-machinery/talos-os-config-generator/src/generator"
)

type GeneratedConfig struct {
	ControlplaneConfig []byte `json:"ControlplanConfig"`
	WorkerConfig       []byte `json:"WorkerConfig"`
	TalosConfig        []byte `json:"TalosConfig"`
}

type ConfigRequest struct {
	ClusterName     string `json:"ClusterName"`
	ControlEndpoint string `json:"ControlEndpoint"`
	IpAddress       string `json:"IpAddress"`
	ConfigPatch     []byte `json:"ConfigPatch"`
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/generate", generate)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

// Handler
func generate(c echo.Context) error {
	configRequest := new(ConfigRequest)
	if err := c.Bind(configRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var controlplaneConfig []byte
	var workerConfig []byte

	configBundle := generator.GenerateConfig(configRequest.ClusterName, configRequest.ControlEndpoint, configRequest.IpAddress)

	var err error
	if configRequest.ConfigPatch != nil {
		controlplaneConfig, workerConfig, err = generator.ApplyPatch(configBundle, configRequest.ConfigPatch)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Cannot apply patch.")
		}

	} else {
		controlplaneConfig, err = configBundle.ControlplaneConfig.Bytes()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Cannot generate controlplane config.")
		}
		workerConfig, err = configBundle.WorkerConfig.Bytes()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Cannot generate worker config.")
		}
	}

	return c.JSON(http.StatusOK, GeneratedConfig{
		controlplaneConfig,
		workerConfig,
		configBundle.TalosConfig,
	})
}

/*

marshaledCfg, err = cfg.Bytes()
		if err != nil {
			log.Fatalf("failed to generate config for node %q: %s", node, err)
		}
*/
