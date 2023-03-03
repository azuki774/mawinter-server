// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
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
