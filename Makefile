clean:
	bazel clean --expunge

dep-ensure:
	go mod tidy

gazelle-repos:
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies

gazelle:
	gazelle-repos
	bazel run //:gazelle

build:
	gazelle
	bazel build //...

go-setup: dep-ensure gazelle

image-push:
	bazel run push_frontend --define=TAG=latest
	bazel run push_order --define=TAG=latest

k8s-deploy:
	helm upgrade frontend k8s \
		--values k8s/values.frontend.yaml \
		--namespace learngo \
		--install --atomic --cleanup-on-fail
	helm upgrade order k8s \
		--values k8s/values.order.yaml \
		--namespace learngo \
		--install --atomic --cleanup-on-fail