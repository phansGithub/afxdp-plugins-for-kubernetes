package goclient

import (
	"fmt"

	"github.com/intel/afxdp-plugins-for-kubernetes/constants"

	"github.com/intel/afxdp-plugins-for-kubernetes/internal/uds"
)

// This is just to make things compile
func compile() {
	fmt.Println(constants.Uds.Protocol)
	uds.NewHandler()
}
