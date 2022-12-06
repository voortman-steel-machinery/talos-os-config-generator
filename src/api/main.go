package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	e.GET("/generate-config", generateConfig)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

// Handler
func generateConfig(c echo.Context) error {
	configRequest := new(ConfigRequest)
	if err := c.Bind(configRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var controlplaneConfig []byte
	var workerConfig []byte

	var err error
	configBundle, err := GenerateConfig(configRequest.ClusterName, configRequest.ControlEndpoint, configRequest.IpAddress)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Cannot create config.")
	}

	if configRequest.ConfigPatch != nil {
		controlplaneConfig, workerConfig, err = ApplyPatch(configBundle, configRequest.ConfigPatch)
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
