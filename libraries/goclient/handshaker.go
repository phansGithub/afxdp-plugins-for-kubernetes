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
func RequestXSKmapFD(hostname string) {

func CreateSession() {
	hWR = uds.NewHandler()
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
