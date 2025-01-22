all:

deploy: all
	cp -r web/. /var/www/

run: deploy
ifneq ($(shell id -u), 0)
	@echo "You must be root to run the server"
else
	nginx -p $$PWD -e stderr -c nginx.conf & cd server && go run *.go; pkill nginx
endif

drop_database:
	sqlite3 server/database/games.db < server/database/create-tables.sql

.PHONY: all deploy run
