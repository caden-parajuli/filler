SERVER_DIR=server
WEB_DIR=web
DIST_DIR=$(WEB_DIR)/dist
DATABASE=server/database/games.db

all:
	esbuild $(WEB_DIR)/index.js --outdir=$(DIST_DIR) --bundle --minify --sourcemap 

deploy_web: all 
ifneq ($(shell id -u), 0)
	($error You must be root to deploy)
	exit 1
else
	cp -r web/. /var/www/
endif

run: deploy_web
	nginx -p $$PWD -e stderr -c nginx.conf & cd $(SERVER_DIR) && go run *.go; pkill nginx

drop_database:
	sqlite3 $(DATABASE) < server/database/create-tables.sql

.PHONY: all deploy_web run drop_database
