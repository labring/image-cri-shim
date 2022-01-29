package server

const (
	// SealosShimSock is the CRI socket the shim listens on.
	SealosShimSock = "/var/run/cri-resmgr/cri-resmgr.sock"
	// DirPermissions is the permissions to create the directory for sockets with.
	DirPermissions = 0711
)

var SealosHub = "sealos.hub:5000"
