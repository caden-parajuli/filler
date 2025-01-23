SERVER_DIR=server
DATABASE=$(SERVER_DIR)/database/games.db

WEB_DIR=web
JS_SRC_DIR=$(WEB_DIR)/src
DIST_DIR=$(WEB_DIR)/dist
STATIC_DIR=$(WEB_DIR)/static
DEPLOY_DIR=/var/www

ESBUILD_SCRIPT_DIR=$(WEB_DIR)/esbuild
ESBUILD_SCRIPT=$(ESBUILD_SCRIPT_DIR)/esbuild


all: $(ESBUILD_SCRIPT)
	$(ESBUILD_SCRIPT) -src=$(JS_SRC_DIR) -out=$(DIST_DIR)


$(ESBUILD_SCRIPT):
	go build -C $(ESBUILD_SCRIPT_DIR)


deploy_web: all 
ifneq ($(shell id -u), 0)
	$(error You must be root to deploy)
	exit 1
else
	cp -r $(STATIC_DIR)/. $(DEPLOY_DIR)/
	cp -r $(DIST_DIR)/. $(DEPLOY_DIR)/
endif


test_run: deploy_web
	nginx -p $$PWD -e stderr -c nginx.conf & cd $(SERVER_DIR) && go run *.go; pkill nginx


drop_database:
	sqlite3 $(DATABASE) < server/database/create-tables.sql


.PHONY: all deploy_web test_run drop_database
