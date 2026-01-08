package docs

import _ "embed"

//go:embed swagger.json
var JsonSwagger string

//go:embed swagger.yaml
var YamlSwagger string
