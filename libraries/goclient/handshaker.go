package goclient

import (
	"fmt"
	"time"

	"github.com/intel/afxdp-plugins-for-kubernetes/constants"

	"github.com/intel/afxdp-plugins-for-kubernetes/internal/host"
	"github.com/intel/afxdp-plugins-for-kubernetes/internal/uds"
	"github.com/intel/afxdp-plugins-for-kubernetes/internal/udsserver"
	logging "github.com/sirupsen/logrus"
)

var (
	srv  udsserver.Server
	hWR  uds.Handler
	hPod host.Handler
)

/*
Get XSK map FD
Request Busy Poll
Version, do this first?
*/

/*Returns the version of our Handshake as a string*/
func GetVersionStr() string {
	return constants.Uds.Handshake.Version
}

/*Returns the version of our Handshake*/
func GetVersion() {
	makeRequest("/version")
}

/*Gets the XSK map FD, may be broken down into sub-methods*/
func RequestXSKmapFD(devName string) {
	cleaner, err := hWR.Dial()
	if err != nil {
		logError("Failed to dial server", err)
		cleaner()
	}
	defer cleaner()
	hostName, err := hPod.Hostname()
	if err != nil {
		logError("Failed to authenticate hostname", err)
	}
	makeRequest(fmt.Sprintf("connect, %s", hostName))
	makeRequest(fmt.Sprintf("/xsk_map_fd, %s", devName))
}

func CreateSession() {
	hWR = uds.NewHandler()
	err := hWR.Init(constants.Uds.SockDir,
		constants.Uds.Protocol,
		constants.Uds.MsgBufSize,
		constants.Uds.CtlBufSize,
		time.Duration(constants.Uds.MinTimeout),
		"")
	if err != nil {
		logError("Failed to initialize UDS socket", err)
	}
	srv.Start()
}

func logError(message string, e error) {
	logging.Errorf("%s : %v", message, e)
}

func makeRequest(request string) {
	fmt.Println("Request: " + request)
	if err := hWR.Write(request, -1); err != nil {
		logError("ERROR: %v failed to write to socket", err)
	}
	response, fd, err := hWR.Read()
	if err != nil {
		logError("Error Reading", err)
	}
	fmt.Printf("Response: %s, FD: %d", response, fd)
	fmt.Println()
}
