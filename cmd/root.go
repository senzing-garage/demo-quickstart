/*
 */
package cmd

import (
	"context"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/senzing-garage/demo-quickstart/httpserver"
	"github.com/senzing-garage/go-cmdhelping/cmdhelper"
	"github.com/senzing-garage/go-cmdhelping/engineconfiguration"
	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-rest-api-service/senzingrestservice"
	"github.com/senzing-garage/serve-grpc/grpcserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Short string = "HTTP/gRPC server supporting various services"
	Use   string = "demo-quickstart"
	Long  string = `
A server supporting the following services:
    - HTTP: Senzing API server
    - HTTP: Swagger UI
    - HTTP: Xterm
	- gRPC:
    `
)

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var ContextVariablesForMultiPlatform = []option.ContextVariable{
	option.Configuration,
	option.DatabaseUrl,
	option.EngineConfigurationJson,
	option.EngineLogLevel,
	option.EngineModuleName,
	option.GrpcPort,
	option.HttpPort,
	option.LogLevel,
	option.ObserverOrigin,
	option.ObserverUrl,
	option.ServerAddress,
	option.TtyOnly,
	option.XtermAllowedHostnames.SetDefault(getDefaultAllowedHostnames()),
	option.XtermArguments,
	option.XtermCommand,
	option.XtermConnectionErrorLimit,
	option.XtermKeepalivePingTimeout,
	option.XtermMaxBufferSizeBytes,
}

var ContextVariables = append(ContextVariablesForMultiPlatform, ContextVariablesForOsArch...)

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, ContextVariables)
}

// --- Networking -------------------------------------------------------------

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func getDefaultAllowedHostnames() []string {
	result := []string{"localhost"}
	outboundIpAddress := getOutboundIP().String()
	if len(outboundIpAddress) > 0 {
		result = append(result, outboundIpAddress)
	}
	return result
}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Used in construction of cobra.Command
func PreRun(cobraCommand *cobra.Command, args []string) {
	cmdhelper.PreRun(cobraCommand, args, Use, ContextVariables)
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	var err error = nil
	ctx := context.Background()

	// Set default value for SENZING_TOOLS_DATABASE_URL.

	_, isSet := os.LookupEnv("SENZING_TOOLS_DATABASE_URL")
	if !isSet {
		err = os.Setenv("SENZING_TOOLS_DATABASE_URL", SENZING_TOOLS_DATABASE_URL)
		if err != nil {
			return err
		}
	}

	// Build configuration for Senzing engine.

	senzingEngineConfigurationJson, err := engineconfiguration.BuildAndVerifySenzingEngineConfigurationJson(ctx, viper.GetViper())
	if err != nil {
		return err
	}

	// Build observers.

	observers := []observer.Observer{}

	// Setup gRPC server

	grpcserver := &grpcserver.GrpcServerImpl{
		EnableAll:                      true,
		LogLevelName:                   viper.GetString(option.LogLevel.Arg),
		ObserverOrigin:                 viper.GetString(option.ObserverOrigin.Arg),
		ObserverUrl:                    viper.GetString(option.ObserverUrl.Arg),
		Port:                           viper.GetInt(option.GrpcPort.Arg),
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
		SenzingModuleName:              viper.GetString(option.EngineModuleName.Arg),
		SenzingVerboseLogging:          viper.GetInt64(option.EngineLogLevel.Arg),
	}

	// Create object and Serve.

	httpServer := &httpserver.HttpServerImpl{
		ApiUrlRoutePrefix:              "api",
		EnableAll:                      true,
		EntitySearchRoutePrefix:        "entity-search",
		LogLevelName:                   viper.GetString(option.LogLevel.Arg),
		ObserverOrigin:                 viper.GetString(option.ObserverOrigin.Arg),
		Observers:                      observers,
		OpenApiSpecificationRest:       senzingrestservice.OpenApiSpecificationJson,
		ReadHeaderTimeout:              60 * time.Second,
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
		SenzingModuleName:              viper.GetString(option.EngineModuleName.Arg),
		SenzingVerboseLogging:          viper.GetInt64(option.EngineLogLevel.Arg),
		ServerAddress:                  viper.GetString(option.ServerAddress.Arg),
		ServerPort:                     viper.GetInt(option.HttpPort.Arg),
		SwaggerUrlRoutePrefix:          "swagger",
		TtyOnly:                        viper.GetBool(option.TtyOnly.Arg),
		XtermAllowedHostnames:          viper.GetStringSlice(option.XtermAllowedHostnames.Arg),
		XtermArguments:                 viper.GetStringSlice(option.XtermArguments.Arg),
		XtermCommand:                   viper.GetString(option.XtermCommand.Arg),
		XtermConnectionErrorLimit:      viper.GetInt(option.XtermConnectionErrorLimit.Arg),
		XtermKeepalivePingTimeout:      viper.GetInt(option.XtermKeepalivePingTimeout.Arg),
		XtermMaxBufferSizeBytes:        viper.GetInt(option.XtermMaxBufferSizeBytes.Arg),
		XtermUrlRoutePrefix:            "xterm",
	}

	// Start servers.

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()
		err = httpServer.Serve(ctx)
	}()

	go func() {
		defer waitGroup.Done()
		err = grpcserver.Serve(ctx)
	}()

	waitGroup.Wait()

	return nil
}

// Used in construction of cobra.Command
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Command
// ----------------------------------------------------------------------------

// RootCmd represents the command.
var RootCmd = &cobra.Command{
	Use:     Use,
	Short:   Short,
	Long:    Long,
	PreRun:  PreRun,
	RunE:    RunE,
	Version: Version(),
}
