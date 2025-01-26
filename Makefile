SERVER_DIR=server
SERVER_BINARY=$(SERVER_DIR)/server
DATABASE=$(SERVER_DIR)/database/games.db

WEB_DIR=web
JS_SRC_DIR=$(WEB_DIR)/src
DIST_DIR=$(WEB_DIR)/dist
STATIC_DIR=$(WEB_DIR)/static
DEPLOY_DIR=/var/www/filler

ESBUILD_SCRIPT_DIR=$(WEB_DIR)/esbuild
ESBUILD_SCRIPT=$(ESBUILD_SCRIPT_DIR)/esbuild


# Esbuild is fast enough that we can run it every time
all: $(ESBUILD_SCRIPT)
	$(ESBUILD_SCRIPT) -src=$(JS_SRC_DIR) -out=$(DIST_DIR)


$(ESBUILD_SCRIPT):
	go build -C $(ESBUILD_SCRIPT_DIR)


$(SERVER_BINARY): $(wildcard $(SERVER_DIR)/*.go)
	go build -C $(SERVER_DIR)


deploy_web: all 
ifneq ($(shell id -u), 0)
	$(error You must be root to install)
	exit 1
else
	mkdir -p $(DEPLOY_DIR)
	cp -r $(STATIC_DIR)/. $(DEPLOY_DIR)/
	cp -r $(DIST_DIR)/. $(DEPLOY_DIR)/
endif


test_run: deploy_web $(SERVER_BINARY)
	nginx -p $$PWD -e stderr -c nginx.conf & cd $(SERVER_DIR) && ./$(SERVER_BINARY); pkill nginx


drop_database:
	sqlite3 $(DATABASE) < server/database/create-tables.sql

shallow_clean:
	rm -Rf $(DIST_DIR)

deep_clean: shallow_clean drop_database
	rm $(ESBUILD_SCRIPT)
	rm -Rf $(DEPLOY_DIR)/*


.PHONY: all deploy_web test_run drop_database shallow_clean deep_clean
