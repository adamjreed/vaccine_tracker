.PHONY: build clean deploy deploy-prod deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/check cmd/check/main.go

clean:
	rm -rf ./bin ./vendor

deploy: clean build
	sls deploy --verbose

deploy-prod: clean build
	sls deploy --verbose --stage prod

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
