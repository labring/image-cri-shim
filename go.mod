module github.com/labring/image-cri-shim

go 1.16

require (
	github.com/containers/image/v5 v5.21.0
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sealyun/endpoints-operator/library v0.0.0-20220415050637-2f3e971b4b3d
	github.com/spf13/cobra v1.4.0
	google.golang.org/grpc v1.44.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/apimachinery v0.22.5
	k8s.io/cri-api v0.23.1
	k8s.io/klog/v2 v2.40.1

)
