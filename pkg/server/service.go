package server

import (
	"context"
	"github.com/sealyun-market/image-cri-shim/pkg/utils"
	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/klog/v2"
	"strings"
)

const (
	legacyDefaultDomain = "index.docker.io"
	defaultDomain       = "docker.io"
	officialRepoName    = "library"
	defaultTag          = "latest"
)

func (s *server) ListImages(ctx context.Context,
	req *api.ListImagesRequest) (*api.ListImagesResponse, error) {

	rsp, err := (*s.imageService).ListImages(ctx, req)

	if err != nil {
		return nil, err
	}

	return rsp, err
}

func (s *server) ImageStatus(ctx context.Context,
	req *api.ImageStatusRequest) (*api.ImageStatusResponse, error) {
	if req.Image != nil {
		req.Image.Image = s.replaceImage(req.Image.Image, "ImageStatus")
	}
	rsp, err := (*s.imageService).ImageStatus(ctx, req)

	if err != nil {
		return nil, err
	}

	return rsp, err
}

func (s *server) PullImage(ctx context.Context,
	req *api.PullImageRequest) (*api.PullImageResponse, error) {
	if req.Image != nil {
		req.Image.Image = s.replaceImage(req.Image.Image, "PullImage")
	}
	rsp, err := (*s.imageService).PullImage(ctx, req)

	if err != nil {
		return nil, err
	}

	return rsp, err
}

func (s *server) RemoveImage(ctx context.Context,
	req *api.RemoveImageRequest) (*api.RemoveImageResponse, error) {
	if req.Image != nil {
		req.Image.Image = s.replaceImage(req.Image.Image, "RemoveImage")
	}
	rsp, err := (*s.imageService).RemoveImage(ctx, req)

	if err != nil {
		return nil, err
	}

	return rsp, err
}

func (s *server) ImageFsInfo(ctx context.Context,
	req *api.ImageFsInfoRequest) (*api.ImageFsInfoResponse, error) {
	rsp, err := (*s.imageService).ImageFsInfo(ctx, req)

	if err != nil {
		return nil, err
	}

	return rsp, err
}
func (s *server) replaceImage(image, action string) string {
	// TODO we can change the image name of req, and make the cri pull the image we need.
	// for example:
	// req.Image.Image = "sealer.hub/library/nginx:1.1.1"
	// and the cri will pull "sealer.hub/library/nginx:1.1.1", and save it as "sealer.hub/library/nginx:1.1.1"
	// note:
	// but kubelet sometimes will invoke imageService.RemoveImage() or something else. The req.Image.Image will the original name.
	// so we'd better tag "sealer.hub/library/nginx:1.1.1" with original name "req.Image.Image" After "rsp, err := (*s.imageService).PullImage(ctx, req)".
	//for image id
	images, err := utils.RunBashCmd("crictl images -q")
	if err != nil {
		klog.Warning("exec crictl images -q error: %s", err.Error())
		return image
	}
	if utils.IsImageId(images, image) {
		klog.Infof("image: %s is imageID,skip replace")
		return image
	}
	//for image name
	domain, named := splitDockerDomain(image)
	klog.Infof("domain: %s,named: %s,action: %s", domain, named, action)
	if len(ShimImages) == 0 || (len(ShimImages) != 0 && utils.NotIn(image, ShimImages)) {
		if utils.RegistryHasImage(SealosHub, Base64Auth, named) {
			newImage := getRegistrDomain() + "/" + named
			klog.Infof("begin image: %s ,after image: %s", image, newImage, action)
			return newImage
		}
		klog.Infof("skip replace images %s", image)
		return image
	}

	fixImageName := image
	if SealosHub != "" {
		if domain != "" {
			fixImageName = getRegistrDomain() + "/" + named
		}
	}

	if Debug {
		klog.Infof("begin image: %s ,after image: %s , action: %s", image, fixImageName, action)
	}
	return fixImageName
}
func splitDockerDomain(name string) (domain, remainder string) {
	i := strings.IndexRune(name, '/')
	if i == -1 || (!strings.ContainsAny(name[:i], ".:") && name[:i] != "localhost" && strings.ToLower(name[:i]) == name[:i]) {
		domain, remainder = "", name
	} else {
		domain, remainder = name[:i], name[i+1:]
	}
	if domain == legacyDefaultDomain || domain == "" {
		domain = defaultDomain
	}
	if domain == defaultDomain && !strings.ContainsRune(remainder, '/') {
		remainder = officialRepoName + "/" + remainder
	}
	return
}
