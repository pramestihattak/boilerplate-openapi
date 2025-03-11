all: validate

clean:
	rm -fv api.yaml

SRC_FILES = $(shell find src -name '*.yaml')
build: $(SRC_FILES)
	@echo "Generating api.yaml"
	npx swagger-cli bundle -t yaml src/main.yaml -o api.yaml

validate: clean build
	npx swagger-cli validate api.yaml

generate-gen:
	rm -rf gen
	mkdir -p ../backend/api
	mkdir -p gen/ts
	oapi-codegen -generate types,chi-server -package api -o ../backend/api/api.gen.go api.yaml
	openapi-typescript api.yaml --output gen/ts/types.ts