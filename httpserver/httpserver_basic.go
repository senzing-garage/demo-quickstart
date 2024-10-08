package httpserver

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/docktermj/cloudshell/xtermservice"
	"github.com/flowchartsman/swaggerui"
	"github.com/pkg/browser"
	"github.com/senzing-garage/demo-entity-search/entitysearchservice"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-rest-api-service-legacy/restapiservicelegacy"
	"github.com/senzing-garage/go-rest-api-service/senzingrestapi"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicHTTPServer is the default implementation of the HttpServer interface.
type BasicHTTPServer struct {
	APIUrlRoutePrefix         string // FIXME: Only works with "api"
	AvoidServing              bool
	EnableAll                 bool
	EnableEntitySearch        bool
	EnableJupyterLab          bool
	EnableSenzingRestAPI      bool
	EnableSwaggerUI           bool
	EnableXterm               bool
	EntitySearchRoutePrefix   string // FIXME: Only works with "entity-search"
	GrpcDialOptions           []grpc.DialOption
	GrpcTarget                string
	JupyterLabRoutePrefix     string
	LogLevelName              string
	ObserverOrigin            string
	Observers                 []observer.Observer
	OpenAPISpecificationRest  []byte
	ReadHeaderTimeout         time.Duration
	SenzingInstanceName       string
	SenzingSettings           string
	SenzingVerboseLogging     int64
	ServerAddress             string
	ServerOptions             []senzingrestapi.ServerOption
	ServerPort                int
	SwaggerURLRoutePrefix     string // FIXME: Only works with "swagger"
	TtyOnly                   bool
	XtermAllowedHostnames     []string
	XtermArguments            []string
	XtermCommand              string
	XtermConnectionErrorLimit int
	XtermKeepalivePingTimeout int
	XtermMaxBufferSizeBytes   int
	XtermURLRoutePrefix       string // FIXME: Only works with "xterm"
}

type TemplateVariables struct {
	APIServerStatus string
	APIServerURL    string
	BasicHTTPServer
	EntitySearchStatus string
	EntitySearchURL    string
	HTMLTitle          string
	JupyterLabStatus   string
	JupyterLabURL      string
	RequestHost        string
	SwaggerStatus      string
	SwaggerURL         string
	XtermStatus        string
	XtermURL           string
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

//go:embed static/*
var static embed.FS

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

func (httpServer *BasicHTTPServer) Serve(ctx context.Context) error {
	rootMux := http.NewServeMux()
	var userMessage string

	// Enable Senzing HTTP REST API.

	if httpServer.EnableAll || httpServer.EnableSenzingRestAPI {
		senzingAPIMux := httpServer.getSenzingRestAPIMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.APIUrlRoutePrefix), http.StripPrefix("/api", senzingAPIMux))
		userMessage = fmt.Sprintf("%sServing Senzing REST API at http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.APIUrlRoutePrefix)
	}

	// Enable Senzing HTTP REST API as reverse proxy.

	if httpServer.EnableAll || httpServer.EnableSenzingRestAPI || httpServer.EnableEntitySearch {
		senzingAPIProxyMux := httpServer.getSenzingRestAPIProxyMux(ctx)
		rootMux.Handle("/entity-search/api/", http.StripPrefix("/entity-search/api", senzingAPIProxyMux))
		userMessage = fmt.Sprintf("%sServing Senzing REST API Reverse Proxy at http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, "entity-search/api")
	}

	// Enable Senzing Entity Search.

	if httpServer.EnableAll || httpServer.EnableEntitySearch {
		entitySearchMux := httpServer.getEntitySearchMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.EntitySearchRoutePrefix), http.StripPrefix("/entity-search", entitySearchMux))
		userMessage = fmt.Sprintf("%sServing Entity Search at    http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.EntitySearchRoutePrefix)
	}

	// Enable SwaggerUI.

	if httpServer.EnableAll || httpServer.EnableSwaggerUI {
		swaggerUIMux := httpServer.getSwaggerUIMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.SwaggerURLRoutePrefix), http.StripPrefix("/swagger", swaggerUIMux))
		userMessage = fmt.Sprintf("%sServing SwaggerUI at        http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.SwaggerURLRoutePrefix)
	}

	// Enable JupyterLab.

	if httpServer.EnableAll || httpServer.EnableJupyterLab {
		proxy, err := newReverseProxy("http://localhost:8888")
		if err != nil {
			panic(err)
		}
		rootMux.HandleFunc(fmt.Sprintf("/%s/", httpServer.JupyterLabRoutePrefix), reverseProxyRequestHandler(proxy))
		userMessage = fmt.Sprintf("%sServing JupyterLab at       http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.JupyterLabRoutePrefix)
	}

	// Enable Xterm.

	if httpServer.EnableAll || httpServer.EnableXterm {
		err := os.Setenv("SENZING_ENGINE_CONFIGURATION_JSON", httpServer.SenzingSettings)
		if err != nil {
			panic(err)
		}
		xtermMux := httpServer.getXtermMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.XtermURLRoutePrefix), http.StripPrefix("/xterm", xtermMux))
		userMessage = fmt.Sprintf("%sServing XTerm at            http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.XtermURLRoutePrefix)
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

	if !httpServer.AvoidServing {
		err = server.ListenAndServe()
	}
	return err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

func (httpServer *BasicHTTPServer) getServerStatus(up bool) string {
	result := "red"
	if httpServer.EnableAll {
		result = "green"
	}
	if up {
		result = "green"
	}
	return result
}

func (httpServer *BasicHTTPServer) getServerURL(up bool, url string) string {
	result := ""
	if httpServer.EnableAll {
		result = url
	}
	if up {
		result = url
	}
	return result
}

func (httpServer *BasicHTTPServer) openAPIFunc(ctx context.Context, openAPISpecification []byte) http.HandlerFunc {
	_ = ctx
	_ = openAPISpecification
	return func(w http.ResponseWriter, r *http.Request) {
		var bytesBuffer bytes.Buffer
		bufioWriter := bufio.NewWriter(&bytesBuffer)
		openAPISpecificationTemplate, err := template.New("OpenApiTemplate").Parse(string(httpServer.OpenAPISpecificationRest))
		if err != nil {
			panic(err)
		}
		templateVariables := TemplateVariables{
			RequestHost: r.Host,
		}
		err = openAPISpecificationTemplate.Execute(bufioWriter, templateVariables)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(bytesBuffer.Bytes())
		if err != nil {
			panic(err)
		}
	}
}
func (httpServer *BasicHTTPServer) populateStaticTemplate(responseWriter http.ResponseWriter, request *http.Request, filepath string, templateVariables TemplateVariables) {
	_ = request
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

// ----------------------------------------------------------------------------
// Methods for Go-based API server - in development
// ----------------------------------------------------------------------------

// func (httpServer *HttpServerImpl) getSenzingRestApiGenericMux(ctx context.Context, urlRoutePrefix string) *senzingrestapi.Server {
// 	service := &senzingrestservice.SenzingRestServiceImpl{
// 		GrpcDialOptions:                httpServer.GrpcDialOptions,
// 		GrpcTarget:                     httpServer.GrpcTarget,
// 		LogLevelName:                   httpServer.LogLevelName,
// 		ObserverOrigin:                 httpServer.ObserverOrigin,
// 		Observers:                      httpServer.Observers,
// 		SenzingEngineConfigurationJson: httpServer.SenzingEngineConfigurationJson,
// 		SenzingModuleName:              httpServer.SenzingModuleName,
// 		SenzingVerboseLogging:          httpServer.SenzingVerboseLogging,
// 		UrlRoutePrefix:                 urlRoutePrefix,
// 		OpenApiSpecificationSpec:       httpServer.OpenApiSpecificationRest,
// 	}
// 	srv, err := senzingrestapi.NewServer(service, httpServer.ServerOptions...)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return srv
// }

// func (httpServer *HttpServerImpl) getSenzingRestApiMux(ctx context.Context) *senzingrestapi.Server {
// 	return httpServer.getSenzingRestApiGenericMux(ctx, "/api")
// }

// func (httpServer *HttpServerImpl) getSenzingApiProxyMux(ctx context.Context) *senzingrestapi.Server {
// 	return httpServer.getSenzingApiGenericMux(ctx, "/entity-search/api")
// }

// --- http.ServeMux ----------------------------------------------------------

func (httpServer *BasicHTTPServer) getEntitySearchMux(ctx context.Context) *http.ServeMux {
	service := &entitysearchservice.BasicHTTPService{}
	return service.Handler(ctx)
}

func (httpServer *BasicHTTPServer) getSenzingRestAPIMux(ctx context.Context) *http.ServeMux {
	service := &restapiservicelegacy.RestApiServiceLegacyImpl{
		JarFile:         "/app/senzing-poc-server.jar",
		ProxyTemplate:   "http://localhost:8250%s",
		CustomTransport: http.DefaultTransport,
	}
	return service.Handler(ctx)
}

func (httpServer *BasicHTTPServer) getSenzingRestAPIProxyMux(ctx context.Context) *http.ServeMux {
	service := &restapiservicelegacy.RestApiServiceLegacyImpl{
		JarFile:         "/app/senzing-poc-server.jar",
		ProxyTemplate:   "http://localhost:8250%s",
		CustomTransport: http.DefaultTransport,
	}
	return service.Handler(ctx)
}

func (httpServer *BasicHTTPServer) getSwaggerUIMux(ctx context.Context) *http.ServeMux {
	swaggerMux := swaggerui.Handler([]byte{}) // OpenAPI specification handled by openApiFunc()
	swaggerFunc := swaggerMux.ServeHTTP
	submux := http.NewServeMux()
	submux.HandleFunc("/", swaggerFunc)
	submux.HandleFunc("/swagger_spec", httpServer.openAPIFunc(ctx, httpServer.OpenAPISpecificationRest))
	return submux
}

func (httpServer *BasicHTTPServer) getXtermMux(ctx context.Context) *http.ServeMux {
	xtermService := &xtermservice.XtermServiceImpl{
		AllowedHostnames:     httpServer.XtermAllowedHostnames,
		Arguments:            httpServer.XtermArguments,
		Command:              httpServer.XtermCommand,
		ConnectionErrorLimit: httpServer.XtermConnectionErrorLimit,
		KeepalivePingTimeout: httpServer.XtermKeepalivePingTimeout,
		MaxBufferSizeBytes:   httpServer.XtermMaxBufferSizeBytes,
		UrlRoutePrefix:       httpServer.XtermURLRoutePrefix,
	}
	return xtermService.Handler(ctx)
}

// --- Http Funcs -------------------------------------------------------------

func (httpServer *BasicHTTPServer) siteFunc(w http.ResponseWriter, r *http.Request) {
	templateVariables := TemplateVariables{
		BasicHTTPServer:    *httpServer,
		HTMLTitle:          "Senzing Tools",
		APIServerStatus:    httpServer.getServerStatus(httpServer.EnableSenzingRestAPI),
		APIServerURL:       httpServer.getServerURL(httpServer.EnableSenzingRestAPI, fmt.Sprintf("http://%s/api", r.Host)),
		EntitySearchStatus: httpServer.getServerStatus(httpServer.EnableEntitySearch),
		EntitySearchURL:    httpServer.getServerURL(httpServer.EnableEntitySearch, fmt.Sprintf("http://%s/entity-search", r.Host)),
		SwaggerStatus:      httpServer.getServerStatus(httpServer.EnableSwaggerUI),
		SwaggerURL:         httpServer.getServerURL(httpServer.EnableSwaggerUI, fmt.Sprintf("http://%s/swagger", r.Host)),
		XtermStatus:        httpServer.getServerStatus(httpServer.EnableXterm),
		XtermURL:           httpServer.getServerURL(httpServer.EnableXterm, fmt.Sprintf("http://%s/xterm", r.Host)),
	}
	w.Header().Set("Content-Type", "text/html")
	filePath := fmt.Sprintf("static/templates%s", r.RequestURI)
	httpServer.populateStaticTemplate(w, r, filePath, templateVariables)
}

// newReverseProxy takes target host and creates a reverse proxy
func newReverseProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(url), nil
}

// reverseProxyRequestHandler handles the http request using proxy
func reverseProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}
