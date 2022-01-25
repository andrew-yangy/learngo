TAG = latest
NAMESPACE = learngo

clean:
	bazel clean --expunge

dep-ensure:
	go mod tidy

gazelle-repos:
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies

gazelle: gazelle-repos
	bazel run //:gazelle

.PHONY: build
build: gazelle
	bazel build //...

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