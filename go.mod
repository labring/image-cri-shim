module github.com/sealyun-market/image-cri-shim

go 1.16

require (
	github.com/pkg/errors v0.9.1
	github.com/sealyun/endpoints-operator/library v0.0.0-20220415050637-2f3e971b4b3d
	github.com/spf13/cobra v1.3.0
	google.golang.org/grpc v1.42.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/apimachinery v0.21.1
	k8s.io/cri-api v0.23.1
	k8s.io/klog/v2 v2.40.1

)
