package server

const (
	// SealosShimSock is the CRI socket the shim listens on.
	SealosShimSock = "/var/run/image-cri-shim.sock"
	// DirPermissions is the permissions to create the directory for sockets with.
	DirPermissions = 0711
)

var SealosHub = ""
var IgnoreHub []string
var Debug = false
