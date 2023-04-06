// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
)

// Defines values for ControllerActionAction.
const (
	Off ControllerActionAction = "off"
	On  ControllerActionAction = "on"
)

// Defines values for ControllerStatusStatus.
const (
	Active   ControllerStatusStatus = "active"
	Inactive ControllerStatusStatus = "inactive"
)

// ControllerAction Action on the specified controller
type ControllerAction struct {
	// Action Type of the action to be performed
	Action *ControllerActionAction `json:"action,omitempty"`

	// Duration Activate controller for a given time in MilliSeconds
	Duration int `json:"duration"`
}

// ControllerActionAction Type of the action to be performed
type ControllerActionAction string

// ControllerStatus Status of the specified controller
type ControllerStatus struct {
	// Description Description of the controller
	Description *string                `json:"description,omitempty"`
	Identifier  string                 `json:"identifier"`
	Status      ControllerStatusStatus `json:"status"`

	// StatusTime For how long controller is in specified status
	StatusTime *int `json:"status-time,omitempty"`
}

// ControllerStatusStatus defines model for ControllerStatus.Status.
type ControllerStatusStatus string

// Controllers List of configured Controllers
type Controllers struct {
	// Description Description of the controller
	Description *string `json:"description,omitempty"`
	Identifier  string  `json:"identifier"`
}

// Health health check data model
type Health struct {
	ApiVersion         *string    `json:"api_version,omitempty"`
	ApplicationName    *string    `json:"application_name,omitempty"`
	ApplicationVersion *string    `json:"application_version,omitempty"`
	Status             string     `json:"status"`
	Timestamp          *time.Time `json:"timestamp,omitempty"`
}

// ControllerActionJSONRequestBody defines body for ControllerAction for application/json ContentType.
type ControllerActionJSONRequestBody = ControllerAction

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /service/health)
	GetServiceHealth(w http.ResponseWriter, r *http.Request)

	// (GET /v1/controllers)
	GetControllers(w http.ResponseWriter, r *http.Request)

	// (GET /v1/controllers/{controllerSlug})
	GetControllerStatus(w http.ResponseWriter, r *http.Request, controllerSlug string)

	// (POST /v1/controllers/{controllerSlug})
	ControllerAction(w http.ResponseWriter, r *http.Request, controllerSlug string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// GetServiceHealth operation middleware
func (siw *ServerInterfaceWrapper) GetServiceHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetServiceHealth(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetControllers operation middleware
func (siw *ServerInterfaceWrapper) GetControllers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetControllers(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetControllerStatus operation middleware
func (siw *ServerInterfaceWrapper) GetControllerStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "controllerSlug" -------------
	var controllerSlug string

	err = runtime.BindStyledParameter("simple", false, "controllerSlug", mux.Vars(r)["controllerSlug"], &controllerSlug)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "controllerSlug", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetControllerStatus(w, r, controllerSlug)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// ControllerAction operation middleware
func (siw *ServerInterfaceWrapper) ControllerAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "controllerSlug" -------------
	var controllerSlug string

	err = runtime.BindStyledParameter("simple", false, "controllerSlug", mux.Vars(r)["controllerSlug"], &controllerSlug)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "controllerSlug", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ControllerAction(w, r, controllerSlug)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshallingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshallingParamError) Error() string {
	return fmt.Sprintf("Error unmarshalling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshallingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{})
}

type GorillaServerOptions struct {
	BaseURL          string
	BaseRouter       *mux.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r *mux.Router) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r *mux.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options GorillaServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = mux.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.HandleFunc(options.BaseURL+"/service/health", wrapper.GetServiceHealth).Methods("GET")

	r.HandleFunc(options.BaseURL+"/v1/controllers", wrapper.GetControllers).Methods("GET")

	r.HandleFunc(options.BaseURL+"/v1/controllers/{controllerSlug}", wrapper.GetControllerStatus).Methods("GET")

	r.HandleFunc(options.BaseURL+"/v1/controllers/{controllerSlug}", wrapper.ControllerAction).Methods("POST")

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xXS2/cNhD+K8S0R60luwWK6NS0QZOgLRrEzqWBEXDJWYkxRTIkJdcw9r8Xo7d25Ucf",
	"aQ71bSVqXt98/Gb2FoStnDVoYoD8FvAPXjmN7e+zLNuUyHUs6anhukb6wZ36YHiFkEOJWlt2bb2WkLQH",
	"DfqgrIEcTk+ykwwSkDzih6ja78+ys9NN9myTfXtxdppnWZ5lv0MCIfJYB8ihdrDf7xMIosSKt0n8aE30",
	"Vmv0z0VsPd+CxCC8ct0jdO+ZNSyWyIJDoXYKJROjJSTgvHXoo+oq4zNXO17rCDlYQ8kuPF/cOGR21/rt",
	"TFi0bIvMod9ZXyFVjaauIH/f2dvdDi4TiDeOyg3RK1PAPgFZe3539g2POEuX7axnnBWqQcMIOaYM+1Vp",
	"rc5RWCMDjBGUiVigBwLN46daeZSUzBhvSsZuP6KIlMwE6XkP/GFS3fuh9EdBunBw6O/F9DQ4Xbg6wktJ",
	"NJFienJ2dBzGvAf0laEGNUg07H6staGz23RsPEzyJ+tZaa+ZtqaYd0MFwn8CoQ/+YA9mNYwZ39+OlU78",
	"okIkyIQ1O1XUHiWbf/8Fu3B3tWtFvhp1ZJlTpy9MlCiumOSRs8pK1Mc3dq4tK4zgzmklWs734vTAR/c5",
	"m+h1dETUCZFXjk5JAzhpB2lcx6rkAZhG8kyyuAbX1JcNf5K9f0f2ZpiGJ9378ron/g+6Vz7p3uN1j4yU",
	"2VmKQR3iIrZcr7jSbbTG+pPAa/TfF/TuRNgKEuiX0Rd03KrOUmnevGbdPtkKjMRwNWd5QN8ogczZaySW",
	"bW/YWx7cFr2/efOa6lJRk3sy3Cx4c7jr7hOwDg13CnL4pl9/HY9lC2naR0onShQYj5nRTUpGOBDKM9L2",
	"DqAN0wndawk5vMR43h31U5agD86aMO7xA6JoYseqkRDpx9CxYb75v3Ns4sLXHneQw1fp9E8hHT5OZ/8R",
	"psX9Lqt+rU9fDQZHzfrt59ZN5EVoWdNXfEkv0+Y0PdCMVQBfYmS6lw/ecKX5Vs9vfFgDcCkufwO+x1U+",
	"D7NS/luMtTdj9jPxEwvDCaEZIddASm+nh3NdF/t7UQtr408sxfIe5M6HO+645xXGtknvDwNdkGtdF6mS",
	"x2pME+cKt3y7ETwg64WGVAHy9ipN131ZGMzlJvoa51w8lKbL/6TDPRwrbb67gwk4G1a6M3ntN7KjThz9",
	"T/4nbfg8iH+qMcQfrLz5DGD3Ra+A/XzcYPv1dViZF+UuS9kfEeTsL+W8tnZvlFzMzrpWEpLH7HgOjaRD",
	"asrGeVt4DN01q0O7i4doHY1Wa9bXv2Hve+TcXpnKq5CqwLgQ6CJKxo2kuztkd8LeBWQFRsLdY/QKGxz0",
	"ZTHaZjb3CBvhgr4ZqFx72gfKGF3I01RbwXVpQ8yfZd9lKQ1qET3sL/d/BgAA///apOMV3hIAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}