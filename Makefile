# TODO change packagename
# TODO add decent run
# TODO fix make shell

version     ?= latest
drunkardimg = cadicallegari/drunkard:$(version)
devimg       = drunkarddev
# GOPATH      ?= $(HOME)/go
# packagename  = $(shell pwd | sed "s:"$(GOPATH)"/src/::")
packagename  = cadicallegari/drunkard
workdir      = /go/src/$(packagename)
runargs      = --rm -v `pwd`:$(workdir) --workdir $(workdir) $(devimg)
runcmd       = docker run -ti $(runargs)
runcompose   = docker-compose run $(runargs)
gitversion   = $(version)

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Variable '$*' not set"; \
		exit 1; \
	fi

release: guard-version publish
	git tag -a $(version) -m "Generated release "$(version)
	git push origin $(version)

publish: image
	docker push $(drunkardimg)

image: build
	docker build -t $(drunkardimg) .

imagedev:
	docker build -t $(devimg) -f ./hack/Dockerfile.dev .

vendor: imagedev
	$(runcmd) ./hack/vendor.sh
	sudo chown -R $(USER):$(id -g -n) ./vendor
	sudo chown -R $(USER):$(id -g -n) ./Godeps

build: imagedev
	$(runcmd) go build -v -ldflags "-X main.Version=$(gitversion)" -o ./cmd/drunkard/drunkard ./cmd/drunkard/main.go

analyze:
	$(runcmd) ./hack/analyze.sh

check: imagedev
	$(runcompose) ./hack/check.sh $(pkg) $(test) $(args)

check-integration: imagedev
	$(runcompose) ./hack/check-integration.sh $(pkg) $(test) $(args)

run: image
	docker-compose run --service-ports --entrypoint "/app/drunkard" --rm drunkard

coverage: imagedev
	$(runcmd) ./hack/coverage.sh

coverage-show: coverage
	xdg-open fullcov.html

shell: imagedev
	$(runcmd)

cleanup:
	docker-compose down
