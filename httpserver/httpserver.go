package httpserver

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docktermj/cloudshell/xtermservice"
	"github.com/flowchartsman/swaggerui"
	"github.com/pkg/browser"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-rest-api-service/senzingrestapi"
	"github.com/senzing-garage/go-rest-api-service/senzingrestservice"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// HttpServerImpl is the default implementation of the HttpServer interface.
type HttpServerImpl struct {
	ApiUrlRoutePrefix              string // FIXME: Only works with "api"
	EnableAll                      bool
	EnableSenzingRestAPI           bool
	EnableSwaggerUI                bool
	EnableXterm                    bool
	GrpcDialOptions                []grpc.DialOption
	GrpcTarget                     string
	LogLevelName                   string
	ObserverOrigin                 string
	Observers                      []observer.Observer
	OpenApiSpecificationRest       []byte
	ReadHeaderTimeout              time.Duration
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int64
	ServerAddress                  string
	ServerOptions                  []senzingrestapi.ServerOption
	ServerPort                     int
	SwaggerUrlRoutePrefix          string // FIXME: Only works with "swagger"
	TtyOnly                        bool
	XtermAllowedHostnames          []string
	XtermArguments                 []string
	XtermCommand                   string
	XtermConnectionErrorLimit      int
	XtermKeepalivePingTimeout      int
	XtermMaxBufferSizeBytes        int
	XtermUrlRoutePrefix            string // FIXME: Only works with "xterm"
}

type TemplateVariables struct {
	HttpServerImpl
	ApiServerStatus string
	ApiServerUrl    string
	HtmlTitle       string
	RequestHost     string
	SwaggerStatus   string
	SwaggerUrl      string
	XtermStatus     string
	XtermUrl        string
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

//go:embed static/*
var static embed.FS

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

func (httpServer *HttpServerImpl) getServerStatus(up bool) string {
	result := "red"
	if httpServer.EnableAll {
		result = "green"
	}
	if up {
		result = "green"
	}
	return result
}

func (httpServer *HttpServerImpl) getServerUrl(up bool, url string) string {
	result := ""
	if httpServer.EnableAll {
		result = url
	}
	if up {
		result = url
	}
	return result
}

func (httpServer *HttpServerImpl) openApiFunc(ctx context.Context, openApiSpecification []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bytesBuffer bytes.Buffer
		bufioWriter := bufio.NewWriter(&bytesBuffer)
		openApiSpecificationTemplate, err := template.New("OpenApiTemplate").Parse(string(httpServer.OpenApiSpecificationRest))
		if err != nil {
			panic(err)
		}
		templateVariables := TemplateVariables{
			RequestHost: r.Host,
		}
		err = openApiSpecificationTemplate.Execute(bufioWriter, templateVariables)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(bytesBuffer.Bytes())
		if err != nil {
			panic(err)
		}
	}
}
func (httpServer *HttpServerImpl) populateStaticTemplate(responseWriter http.ResponseWriter, request *http.Request, filepath string, templateVariables TemplateVariables) {
	templateBytes, err := static.ReadFile(filepath)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		return
	}
	templateParsed, err := template.New("HtmlTemplate").Parse(string(templateBytes))
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		return
	}
	err = templateParsed.Execute(responseWriter, templateVariables)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		return
	}
}

// --- http.ServeMux ----------------------------------------------------------

func (httpServer *HttpServerImpl) getSenzingApiMux(ctx context.Context) *senzingrestapi.Server {
	service := &senzingrestservice.SenzingRestServiceImpl{
		GrpcDialOptions:                httpServer.GrpcDialOptions,
		GrpcTarget:                     httpServer.GrpcTarget,
		LogLevelName:                   httpServer.LogLevelName,
		ObserverOrigin:                 httpServer.ObserverOrigin,
		Observers:                      httpServer.Observers,
		SenzingEngineConfigurationJson: httpServer.SenzingEngineConfigurationJson,
		SenzingModuleName:              httpServer.SenzingModuleName,
		SenzingVerboseLogging:          httpServer.SenzingVerboseLogging,
		UrlRoutePrefix:                 httpServer.ApiUrlRoutePrefix,
		OpenApiSpecificationSpec:       httpServer.OpenApiSpecificationRest,
	}
	srv, err := senzingrestapi.NewServer(service, httpServer.ServerOptions...)
	if err != nil {
		log.Fatal(err)
	}
	return srv
}

func (httpServer *HttpServerImpl) getSwaggerUiMux(ctx context.Context) *http.ServeMux {
	swaggerMux := swaggerui.Handler([]byte{}) // OpenAPI specification handled by openApiFunc()
	swaggerFunc := swaggerMux.ServeHTTP
	submux := http.NewServeMux()
	submux.HandleFunc("/", swaggerFunc)
	submux.HandleFunc("/swagger_spec", httpServer.openApiFunc(ctx, httpServer.OpenApiSpecificationRest))
	return submux
}

func (httpServer *HttpServerImpl) getXtermMux(ctx context.Context) *http.ServeMux {
	xtermService := &xtermservice.XtermServiceImpl{
		AllowedHostnames:     httpServer.XtermAllowedHostnames,
		Arguments:            httpServer.XtermArguments,
		Command:              httpServer.XtermCommand,
		ConnectionErrorLimit: httpServer.XtermConnectionErrorLimit,
		KeepalivePingTimeout: httpServer.XtermKeepalivePingTimeout,
		MaxBufferSizeBytes:   httpServer.XtermMaxBufferSizeBytes,
		UrlRoutePrefix:       httpServer.XtermUrlRoutePrefix,
	}
	return xtermService.Handler(ctx)
}

// --- Http Funcs -------------------------------------------------------------

func (httpServer *HttpServerImpl) siteFunc(w http.ResponseWriter, r *http.Request) {
	templateVariables := TemplateVariables{
		HttpServerImpl:  *httpServer,
		HtmlTitle:       "Senzing Tools",
		ApiServerUrl:    httpServer.getServerUrl(httpServer.EnableSenzingRestAPI, fmt.Sprintf("http://%s/api", r.Host)),
		ApiServerStatus: httpServer.getServerStatus(httpServer.EnableSenzingRestAPI),
		SwaggerUrl:      httpServer.getServerUrl(httpServer.EnableSwaggerUI, fmt.Sprintf("http://%s/swagger", r.Host)),
		SwaggerStatus:   httpServer.getServerStatus(httpServer.EnableSwaggerUI),
		XtermUrl:        httpServer.getServerUrl(httpServer.EnableXterm, fmt.Sprintf("http://%s/xterm", r.Host)),
		XtermStatus:     httpServer.getServerStatus(httpServer.EnableXterm),
	}
	w.Header().Set("Content-Type", "text/html")
	filePath := fmt.Sprintf("static/templates%s", r.RequestURI)
	httpServer.populateStaticTemplate(w, r, filePath, templateVariables)
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Serve method simply prints the 'Something' value in the type-struct.

Input
  - ctx: A context to control lifecycle.

Output
  - Nothing is returned, except for an error.  However, something is printed.
    See the example output.
*/

func (httpServer *HttpServerImpl) Serve(ctx context.Context) error {
	rootMux := http.NewServeMux()
	var userMessage string = ""

	// Enable Senzing HTTP REST API.

	if httpServer.EnableAll || httpServer.EnableSenzingRestAPI {
		senzingApiMux := httpServer.getSenzingApiMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.ApiUrlRoutePrefix), http.StripPrefix("/api", senzingApiMux))
		userMessage = fmt.Sprintf("%sServing Senzing REST API at http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.ApiUrlRoutePrefix)
	}

	// Enable SwaggerUI.

	if httpServer.EnableAll || httpServer.EnableSwaggerUI {
		swaggerUiMux := httpServer.getSwaggerUiMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.SwaggerUrlRoutePrefix), http.StripPrefix("/swagger", swaggerUiMux))
		userMessage = fmt.Sprintf("%sServing SwaggerUI at        http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.SwaggerUrlRoutePrefix)
	}

	// Enable Xterm.

	if httpServer.EnableAll || httpServer.EnableXterm {
		err := os.Setenv("SENZING_ENGINE_CONFIGURATION_JSON", httpServer.SenzingEngineConfigurationJson)
		if err != nil {
			panic(err)
		}
		xtermMux := httpServer.getXtermMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.XtermUrlRoutePrefix), http.StripPrefix("/xterm", xtermMux))
		userMessage = fmt.Sprintf("%sServing XTerm at            http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.XtermUrlRoutePrefix)
	}

	// Add route to template pages.

	rootMux.HandleFunc("/site/", httpServer.siteFunc)
	userMessage = fmt.Sprintf("%sServing Console at          http://localhost:%d\n", userMessage, httpServer.ServerPort)

	// Add route to static files.

	rootDir, err := fs.Sub(static, "static/root")
	if err != nil {
		panic(err)
	}
	rootMux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(rootDir))))

	// Start service.

	listenOnAddress := fmt.Sprintf("%s:%v", httpServer.ServerAddress, httpServer.ServerPort)
	userMessage = fmt.Sprintf("%sStarting server on interface:port '%s'...\n", userMessage, listenOnAddress)
	fmt.Println(userMessage)
	server := http.Server{
		ReadHeaderTimeout: httpServer.ReadHeaderTimeout,
		Addr:              listenOnAddress,
		Handler:           rootMux,
	}

	// Start a web browser.  Unless disabled.

	if !httpServer.TtyOnly {
		_ = browser.OpenURL(fmt.Sprintf("http://localhost:%d", httpServer.ServerPort))
	}

	return server.ListenAndServe()
}
