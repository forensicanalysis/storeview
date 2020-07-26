.PHONY: install build

install:
	cd frontend && yarn install
	go get -u github.com/asticode/go-astilectron-bundler/...

build:
	cd frontend && yarn build
	sed 's_=/_=_g' frontend/dist/index.html > tmp
	mv tmp frontend/dist/index.html
	rm -rf cmd/ui/resources/app
	mv frontend/dist cmd/ui/resources/app
	cp -r cmd/ui/resources/start/* cmd/ui/resources/app
	cd cmd/ui && astilectron-bundler
