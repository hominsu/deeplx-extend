SRC_MAKEFILES := $(foreach dir, app, $(wildcard $(dir)/*/Makefile))

.PHONY: init dep api conf ent wire openapi build clean run test

# init env
init:
	go install entgo.io/ent/cmd/ent@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	cd api && buf dep update

# download dependencies of module
dep:
	$(foreach dir, $(dir $(realpath $(SRC_MAKEFILES))),\
      cd $(dir);\
      make dep;\
    )

# generate protobuf api go code
api:
	cd $(lastword $(dir $(realpath $(SRC_MAKEFILES))));\
	make api;

# generate config define code
conf:
	$(foreach dir, $(dir $(realpath $(SRC_MAKEFILES))),\
      cd $(dir);\
      make conf;\
    )

# generate ent code
ent:
	$(foreach dir, $(dir $(realpath $(SRC_MAKEFILES))),\
      cd $(dir);\
      make ent;\
    )

# generate wire code
wire:
	$(foreach dir, $(dir $(realpath $(SRC_MAKEFILES))),\
      cd $(dir);\
      make wire;\
    )

# generate OpenAPI v3 docs.
openapi:
	cd $(lastword $(dir $(realpath $(SRC_MAKEFILES))));\
	make openapi;

# build all service applications
build:
	$(foreach dir, $(dir $(realpath $(SRC_MAKEFILES))),\
      cd $(dir);\
      make build;\
    )

# clean build files
clean:
	$(foreach dir, $(dir $(realpath $(SRC_MAKEFILES))),\
	  cd $(dir);\
	  make clean;\
	)

# run tests
test:
	$(foreach dir, $(dir $(realpath $(SRC_MAKEFILES))),\
	  cd $(dir);\
	  make test;\
	)

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help