.PHONY: install build

install:
	yarn install
	go get -u github.com/asticode/go-astilectron-bundler/...

build:
	yarn build
	sed 's_=/_=_g' dist/index.html > tmp
	mv tmp dist/index.html
	rm -rf app/resources/app
	mv dist app/resources/app
	cp -r app/resources/open/* app/resources/app
	cd app && astilectron-bundler
