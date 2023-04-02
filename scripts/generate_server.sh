#!/bin/bash

PACKAGE="httpapi"
CODEGEN_DIR="./internal/controller/codegen/${PACKAGE}"

mkdir -p ${CODEGEN_DIR}

${HOME}/go/bin/oapi-codegen -package ${PACKAGE} -generate chi-server openapi/openapi.yaml > ${CODEGEN_DIR}/task_server.gen.go
${HOME}/go/bin/oapi-codegen -package ${PACKAGE} -generate types openapi/openapi.yaml > ${CODEGEN_DIR}/task_models.gen.go

