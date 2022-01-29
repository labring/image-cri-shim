/*
Copyright © 2022 cuisongliu@qq.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sealyun-market/sealos-cri-shim/pkg/cri"
	"github.com/sealyun-market/sealos-cri-shim/pkg/server"
	"github.com/sealyun-market/sealos-cri-shim/pkg/shim"
	"github.com/sealyun-market/sealos-cri-shim/pkg/version"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"os"
	"os/signal"
	"syscall"
)

var shimSocket, criSocket string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sealos-cri-shim",
	Short: "cri for kubelet endpoint-image-service",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Version: version.Get(),
	Run: func(cmd *cobra.Command, args []string) {
		run(shimSocket, criSocket)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if shimSocket == "" {
			return errors.New("socket path is empty")
		}
		if server.SealosHub == "" {
			return errors.New("registry addr is empty")
		}
		if criSocket == "" {
			socket, err := cri.DetectCRISocket()
			if err != nil {
				return err
			}
			criSocket = socket
		}
		if !isExist(criSocket) {
			return errors.New("cri is running?")
		}
		return nil
	},
}

func isExist(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		return os.IsExist(err)
	}
	return true
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&shimSocket, "shim-socket", server.SealosShimSock, "The endpoint of local image socket path.")
	rootCmd.Flags().StringVar(&criSocket, "cri-socket", "", "The endpoint of remote image socket path.")
	rootCmd.Flags().StringVar(&server.SealosHub, "registry-address", server.SealosHub, "The registry address.")
}

func run(socket string, criSocket string) {
	options := shim.Options{
		ShimSocket:  socket,
		ImageSocket: criSocket,
	}
	klog.Infof("socket info shim: %v ,image: %v, registry: %v", socket, criSocket, server.SealosHub)
	_shim, err := shim.NewShim(options)
	if err != nil {
		klog.Fatalf("failed to new _shim, %s", err)
	}

	err = _shim.Setup()
	if err != nil {
		klog.Fatalf("failed to setup sealos _shim, %s", err)
	}

	err = _shim.Start()
	if err != nil {
		klog.Fatalf(fmt.Sprintf("failed to start sealos _shim, %s", err))
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	stopCh := make(chan struct{}, 1)
	select {
	case <-signalCh:
		close(stopCh)
	case <-stopCh:

	}
	_ = os.Remove(socket)
	klog.Infof("Shutting down the sealos _shim")
}
