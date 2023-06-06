.PHONY: clean build run run.local swag.init swag swag.local swag.hook
APP_NAME = ekira-backend
BUILD_DIR = $(PWD)/build

clean:
	rm -rf ./build

build: clean
	go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: swag build
	$(BUILD_DIR)/$(APP_NAME)

run.local: swag.local build
	$(BUILD_DIR)/$(APP_NAME)

swag.init:
	swag init

swag: swag.init swag.hook

swag.local: swag.init swag.hook
	sed -i 's/api.e-kira.tk/localhost:5000/g' docs/docs.go

swag.hook:
	sed -i 's/"x-example"/"example"/g' docs/docs.go
	sed -i 's/"x-nullable"/"nullable"/g' docs/docs.go
	sed -i 's/"null"/null/g' docs/docs.go
	sed -i 's/models.//g' docs/docs.go
	sed -i 's/controllers.//g' docs/docs.go
