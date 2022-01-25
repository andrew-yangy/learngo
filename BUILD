load("@bazel_gazelle//:def.bzl", "gazelle")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_docker//container:container.bzl", "container_push")

# gazelle:prefix github.com/ddvkid/learngo
gazelle(
    name = "gazelle",
)

go_image(
    name = "frontend",
    binary = "//services/frontend:frontend_linux_bin",
    importpath = "github.com/ddvkid/learngo/services/frontend",
)

container_push(
    name = "push_frontend",
    image = ":frontend",
    format = "Docker",
    registry = "270878775604.dkr.ecr.us-east-2.amazonaws.com",
    repository = "frontend",
    tag = "$(TAG)",
)

go_image(
    name = "order",
    binary = "//services/order:order_linux_bin",
    importpath = "github.com/ddvkid/learngo/services/order",
)

container_push(
    name = "push_order",
    image = ":order",
    format = "Docker",
    registry = "270878775604.dkr.ecr.us-east-2.amazonaws.com",
    repository = "order",
    tag = "$(TAG)",
)
