TAG = latest
NAMESPACE = learngo
functions := $(shell find functions -name \*main.go | awk -F'/' '{print $$2}')

clean:
	bazel clean --expunge

dep-ensure:
	go mod tidy

gazelle-repos:
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies

gazelle: gazelle-repos
	bazel run //:gazelle

build-services: gazelle
	bazel build //services/...

build-functions:
	@for function in $(functions) ; do \
		env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$$function functions/$$function/main.go ; \
	done

.PHONY: build
build: build-services build-functions

go-setup: dep-ensure gazelle

image-push:
	bazel run push_frontend --define=TAG=$(TAG)
	bazel run push_order --define=TAG=$(TAG)

k8s-deploy:
	helm upgrade frontend k8s \
		--values k8s/values.frontend.yaml \
		--namespace $(NAMESPACE) \
		--install --atomic --cleanup-on-fail
	helm upgrade order k8s \
		--values k8s/values.order.yaml \
		--namespace $(NAMESPACE) \
		--install --atomic --cleanup-on-fail

lambda-deploy:
	serverless deploy