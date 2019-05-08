.RECIPEPREFIX =
.DEFAULT_GOAL = run

define pprint
 @echo -e "\033[32m$(1)\033[0m"
endef

.PHONY: run
run:
	$(call Starting server ...")
	@MYCNC_REST_API_CONFIG_FILE_PATH=${PWD}/config.json GOPATH=${PWD}/.gopath go run main/main.go