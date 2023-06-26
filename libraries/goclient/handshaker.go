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

const (
	requestDelay = 100 * time.Millisecond
)

/*
Get XSK map FD
Request Busy Poll
Version, do this first?
*/

/*
Returns the version of our Handshake as a string
*/
func GetVersionStr() string {
	return constants.Uds.Handshake.Version
}

/*
Returns the version of our Handshake
*/
func GetVersion() {
	makeRequest("/version")
}

/*
Gets the XSK map FD, may be broken down into sub-methods
*/
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
	authString := fmt.Sprintf("connect, %s", hostName)
	makeRequest(authString)
	time.Sleep(requestDelay)
	mapString := fmt.Sprintf("/xsk_map_fd, %s", devName)
	makeRequest(mapString)
	time.Sleep(requestDelay)
}

/*
Call this first to initialize the UDS socket
*/
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

/*
Logs an error with a message
*/
func logError(message string, e error) {
	logging.Errorf("%s : %v", message, e)
}

/*
Makes a request to the server
*/
func makeRequest(request string) {
	fmt.Println("Request: " + request)
	if err := hWR.Write(request, -1); err != nil {
		logError("ERROR: %v failed to write to socket", err)
	}
	response, fd, err := hWR.Read()
	if err != nil {
		logError("Error Reading", err)
	}
	checkXSKReq(request, response)
	fmt.Printf("Response: %s, FD: %d", response, fd)
	fmt.Println()
}

/*
To check if the request for the XSK map FD is either ack or nak
*/
func checkXSKReq(request, response string) {
	fmt.Println(request)
	fmt.Println(response)
}
