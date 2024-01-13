BIN_DIR = bin
PROTO_DIR = proto
SERVER_DIR = cmd/orders
CLIENT_DIR = client

ifeq ($(OS), Windows_NT)
	SERVER_BIN = ${SERVER_DIR}.exe
	CLIENT_BIN = ${CLIENT_DIR}.exe
else
	SERVER_BIN = ${SERVER_DIR}
	CLIENT_BIN = ${CLIENT_DIR}
endif

PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')

.DEFAULT_GOAL := help
.PHONY: help proto clean clean_orders build about



build:
	@${CHECK_DIR_CMD}
	protoc -I${PROTO_DIR} --go_opt=module=${PACKAGE} --go_out=. --go-grpc_opt=module=${PACKAGE} --go-grpc_out=. ${PROTO_DIR}/*.proto
	go build -o ${BIN_DIR}/cmd/orders ./${SERVER_DIR}
	go build -o ${BIN_DIR}/${CLIENT_DIR}/${CLIENT_BIN} ./${CLIENT_DIR}



clean: clean_orders ## Clean generated files
	${RM_F_CMD} ssl/*.crt
	${RM_F_CMD} ssl/*.csr
	${RM_F_CMD} ssl/*.key
	${RM_F_CMD} ssl/*.pem
	${RM_RF_CMD} ${BIN_DIR}

clean_orders: ## Clean generated files for orders
	${RM_F_CMD} ${PROTO_DIR}/*.pb.go


about: ## Display info related to the build
	@echo "OS: ${OS}"
	@echo "Shell: ${SHELL} ${SHELL_VERSION}"
	@echo "Protoc version: $(shell protoc --version)"
	@echo "Go version: $(shell go version)"
	@echo "Go package: ${PACKAGE}"
	@echo "Openssl version: $(shell openssl version)"

help: ## Show this help
	@${HELP_CMD}
