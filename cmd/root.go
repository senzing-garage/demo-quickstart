/*
 */
package cmd

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/senzing-garage/demo-quickstart/httpserver"
	"github.com/senzing-garage/go-cmdhelping/cmdhelper"
	"github.com/senzing-garage/go-cmdhelping/engineconfiguration"
	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/senzing-garage/go-grpcing/grpcurl"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-rest-api-service/senzingrestservice"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	Short string = "HTTP server supporting various services"
	Use   string = "demo-quickstart"
	Long  string = `
An HTTP server supporting the following services:
    - Senzing API server
    - Swagger UI
    - Xterm
    `
)

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var ContextVariablesForMultiPlatform = []option.ContextVariable{
	option.Configuration,
	option.DatabaseUrl,
	option.EnableAll,
	option.EnableSenzingRestApi,
	option.EnableSwaggerUi,
	option.EnableXterm,
	option.EngineConfigurationJson,
	option.EngineLogLevel,
	option.EngineModuleName,
	option.GrpcUrl,
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

	senzingEngineConfigurationJson, err := engineconfiguration.BuildAndVerifySenzingEngineConfigurationJson(ctx, viper.GetViper())
	if err != nil {
		return err
	}

	// Determine if gRPC is being used.

	grpcUrl := viper.GetString(option.GrpcUrl.Arg)
	grpcTarget := ""
	grpcDialOptions := []grpc.DialOption{}
	if len(grpcUrl) > 0 {
		grpcTarget, grpcDialOptions, err = grpcurl.Parse(ctx, grpcUrl)
		if err != nil {
			return err
		}
	}

	// Build observers.

	observers := []observer.Observer{}

	// Create object and Serve.

	httpServer := &httpserver.HttpServerImpl{
		ApiUrlRoutePrefix:              "api",
		EnableAll:                      viper.GetBool(option.EnableAll.Arg),
		EnableSenzingRestAPI:           viper.GetBool(option.EnableSenzingRestApi.Arg),
		EnableSwaggerUI:                viper.GetBool(option.EnableSwaggerUi.Arg),
		EnableXterm:                    viper.GetBool(option.EnableXterm.Arg),
		GrpcDialOptions:                grpcDialOptions,
		GrpcTarget:                     grpcTarget,
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
	return httpServer.Serve(ctx)
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
