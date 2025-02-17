package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "builtins.pb.go",
        "grpc_adaptor.go",
        "grpc_adaptor_mock.go",
        "mizar-arktos-network-controller.go",
        "mizar-endpoints-controller.go",
        "mizar-namespace-controller.go",
        "mizar-network-policy-controller.go",
        "mizar-node-controller.go",
        "mizar-pod-controller.go",
        "mizar-service-controller.go",
        "mizar-starter-controller.go",
        "testutil.go",
        "util.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/controller/mizar",
    deps = [
        "//pkg/controller:go_default_library",
        "//pkg/util/node:go_default_library",
        "//staging/src/k8s.io/api/core/v1:go_default_library",
        "//staging/src/k8s.io/api/networking/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/api/errors:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/util/intstr:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/util/runtime:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//staging/src/k8s.io/arktos-ext/pkg/apis/arktosextensions/v1:go_default_library",
        "//staging/src/k8s.io/arktos-ext/pkg/generated/clientset/versioned:go_default_library",
        "//staging/src/k8s.io/arktos-ext/pkg/generated/clientset/versioned/scheme:go_default_library",
        "//staging/src/k8s.io/arktos-ext/pkg/generated/informers/externalversions/arktosextensions/v1:go_default_library",
        "//staging/src/k8s.io/arktos-ext/pkg/generated/listers/arktosextensions/v1:go_default_library",
        "//staging/src/k8s.io/client-go/informers/core/v1:go_default_library",
        "//staging/src/k8s.io/client-go/informers/networking/v1:go_default_library",
        "//staging/src/k8s.io/client-go/kubernetes:go_default_library",
        "//staging/src/k8s.io/client-go/kubernetes/scheme:go_default_library",
        "//staging/src/k8s.io/client-go/kubernetes/typed/core/v1:go_default_library",
        "//staging/src/k8s.io/client-go/listers/core/v1:go_default_library",
        "//staging/src/k8s.io/client-go/listers/networking/v1:go_default_library",
        "//staging/src/k8s.io/client-go/tools/cache:go_default_library",
        "//staging/src/k8s.io/client-go/tools/record:go_default_library",
        "//staging/src/k8s.io/client-go/util/workqueue:go_default_library",
        "//vendor/github.com/gogo/protobuf/proto:go_default_library",
        "//vendor/google.golang.org/grpc:go_default_library",
        "//vendor/google.golang.org/grpc/codes:go_default_library",
        "//vendor/google.golang.org/grpc/status:go_default_library",
        "//vendor/k8s.io/klog:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "mizar_node_controller_test.go",
        "util_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//pkg/controller:go_default_library",
        "//pkg/controller/testutil:go_default_library",
        "//staging/src/k8s.io/api/core/v1:go_default_library",
        "//staging/src/k8s.io/api/networking/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/util/intstr:go_default_library",
        "//staging/src/k8s.io/client-go/informers:go_default_library",
        "//staging/src/k8s.io/client-go/informers/core/v1:go_default_library",
        "//staging/src/k8s.io/client-go/kubernetes/fake:go_default_library",
        "//staging/src/k8s.io/client-go/tools/cache:go_default_library",
    ],
)
