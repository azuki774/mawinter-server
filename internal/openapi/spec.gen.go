// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
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
	res := make(map[string]func() ([]byte, error))
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
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
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
