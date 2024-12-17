package api

//go:generate go run -modfile=../../tools/go.mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=server.config.yaml ../../api/specs.yaml
//go:generate go run -modfile=../../tools/go.mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=models.config.yaml ../../api/specs.yaml
