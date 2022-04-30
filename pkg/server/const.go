package server

import (
	"github.com/labring/image-cri-shim/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"strings"
	"time"
)

const (
	// SealosShimSock is the CRI socket the shim listens on.
	SealosShimSock = "/var/run/image-cri-shim.sock"
	// DirPermissions is the permissions to create the directory for sockets with.
	DirPermissions = 0711
)

var ShimImages []string
var Debug = false
var (
	Base64Auth string
	Auth       string
	ConfigFile string
	SealosHub  string
)

func getData() map[string]interface{} {
	data, err := utils.Unmarshal(ConfigFile)
	if err != nil {
		klog.Warning("load config from image shim: %v", err)
		return nil
	}
	return data
}

func getRegistrDomain() string {
	domain := SealosHub
	domain = strings.ReplaceAll(domain, "http://", "")
	domain = strings.ReplaceAll(domain, "https://", "")
	return domain
}

func RunLoad() {
	data := getData()
	imageDir, _, _ := unstructured.NestedString(data, "image")
	sync, _, _ := unstructured.NestedInt64(data, "sync")
	if sync != 0 {
		go wait.Forever(func() {
			images, err := utils.LoadImages(imageDir)
			if err != nil {
				klog.Warning("load images from image dir: %v", err)
			}
			ShimImages = images
			klog.Infof("sync image list for image dir,sync second is %d,data is %+v", sync, images)
		}, time.Duration(sync*int64(time.Second)))
	}
}
