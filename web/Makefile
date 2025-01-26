# Install command variables
INSTALL=install -C
INSTALL_PROGRAM=$(INSTALL)
INSTALL_DATA=${INSTALL} -m 644

# Directories
JS_SRC_DIR=src
DIST_DIR=dist
STATIC_DIR=static
DEPLOY_DIR=var/www/filler

ESBUILD_SCRIPT_DIR=esbuild
ESBUILD_SCRIPT=$(ESBUILD_SCRIPT_DIR)/esbuild


# Esbuild is fast enough that we can run it every time
all: $(ESBUILD_SCRIPT)
	$(ESBUILD_SCRIPT) -src=$(JS_SRC_DIR) -out=$(DIST_DIR)


$(ESBUILD_SCRIPT):
	go build -C $(ESBUILD_SCRIPT_DIR)

install: all 
	$(PRE_INSTALL)     # Pre-install commands follow.
	$(POST_INSTALL)    # Post-install commands follow.
	$(NORMAL_INSTALL)  # Normal commands follow.
	$(INSTALL_PROGRAM) -d $(DEST_DIR)/$(DEPLOY_DIR)
	$(INSTALL_DATA) $(STATIC_DIR)/* $(DEST_DIR)/$(DEPLOY_DIR)
	$(INSTALL_DATA) $(DIST_DIR)/* $(DEST_DIR)/$(DEPLOY_DIR)

clean: 
	rm ESBUILD_SCRIPT
	rm -Rf $(DIST_DIR)/*

uninstall:
	$(PRE_UNINSTALL)     # Pre-install commands follow.
	$(POST_UNINSTALL)    # Post-install commands follow.
	$(NORMAL_UNINSTALL)  # Normal commands follow.
	rm -Rf $(DEST_DIR)/$(DEPLOY_DIR)

.PHONY: all install uninstall clean
