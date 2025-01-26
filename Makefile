SERVER_DIR=server
SERVER_BINARY_NAME=server
SERVER_BINARY=$(SERVER_DIR)/$(SERVER_BINARY_NAME)
DATABASE=$(SERVER_DIR)/database/games.db

WEB_DIR=web

ESBUILD_SCRIPT_DIR=$(WEB_DIR)/esbuild
ESBUILD_SCRIPT=$(ESBUILD_SCRIPT_DIR)/esbuild


# Esbuild is fast enough that we can run it every time
all: $(SERVER_BINARY) web

web:
	cd $(WEB_DIR) && make
deploy_web:
	cd $(WEB_DIR) && make install


$(SERVER_BINARY): $(wildcard $(SERVER_DIR)/*.go)
	go build -C $(SERVER_DIR)


test_run: deploy_web $(SERVER_BINARY)
	nginx -p $$PWD -e stderr -c nginx.conf & cd $(SERVER_DIR) && ./$(SERVER_BINARY_NAME); pkill nginx


drop_database:
	sqlite3 $(DATABASE) < $(SERVER_DIR)/database/create-tables.sql

shallow_clean:
	cd $(WEB_DIR) && make clean

deep_clean: shallow_clean drop_database


.PHONY: all web deploy_web test_run drop_database shallow_clean deep_clean
