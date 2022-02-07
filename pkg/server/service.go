package server

import (
	"context"
	"k8s.io/klog/v2"
	"strings"

	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

const (
	legacyDefaultDomain = "index.docker.io"
	defaultDomain       = "docker.io"
	officialRepoName    = "library"
	defaultTag          = "latest"
)

func (s *server) ListImages(ctx context.Context,
	req *api.ListImagesRequest) (*api.ListImagesResponse, error) {
	if req.Filter != nil && req.Filter.Image != nil {
		req.Filter.Image.Image = s.replaceImage(req.Filter.Image.Image, "ListImages")
	}

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
	domain, named := splitDockerDomain(image)
	if domain != "" {
		image = SealosHub + "/" + named
	}
	klog.Infof("image name is: %s , action: %s", image, action)
	return image
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
