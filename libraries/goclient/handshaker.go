package goclient

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/intel/afxdp-plugins-for-kubernetes/constants"

	"github.com/intel/afxdp-plugins-for-kubernetes/internal/host"
	"github.com/intel/afxdp-plugins-for-kubernetes/internal/uds"
)

var (
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
Call this function to initialize the library, returns a cleanup function
*/
func Init() uds.CleanupFunc {
	hWR = uds.NewHandler()
	if err := hWR.Init(constants.Uds.PodPath, constants.Uds.Protocol, constants.Uds.MsgBufSize, constants.Uds.CtlBufSize, 0*time.Second, ""); err != nil {
		println("Library Error: Error Initialising UDS server: ", err)
		os.Exit(1)
}

	cleanup, err := hWR.Dial()
	if err != nil {
		println("Library Error: UDS Dial error: ", err)
		cleanup()
		os.Exit(1)
	}

	return cleanup
	}

/*
Returns the version of our Handshake
*/
func GetVersion() {
	makeRequest("/connect, afxdp-e2e-test")
	makeRequest("/version")
}

	authString := fmt.Sprintf("connect, %s", hostName)
	makeRequest(authString)
	time.Sleep(requestDelay)

/*
Give it a list of device names and returns a map of the fds for each device and a cleanup function to close the connection
*/
func RequestXSKmapFD(devNames []string) map[string]int {
	// hWR = uds.NewHandler()
	fds := make(map[string]int)
	// init
	// if err := hWR.Init(constants.Uds.PodPath, constants.Uds.Protocol, constants.Uds.MsgBufSize, constants.Uds.CtlBufSize, 0*time.Second, ""); err != nil {
	// 	println("Test App Error: Error Initialising UDS server: ", err)
	// 	os.Exit(1)
	// }
	// // Execute timeoutBeforeConnect when set to true

	// cleanup, err := hWR.Dial()
	// if err != nil {
	// 	println("Test App Error: UDS Dial error:: ", err)
	// 	cleanup()
	// 	os.Exit(1)
	// }

	// connect and verify pod hostname
	makeRequest("/connect, afxdp-e2e-test")
	for _, dev := range devNames {
		fd := makeRequest("/xsk_map_fd, " + dev)
		fds[dev] = fd
	}
	makeRequest("/fin")

	return fds
}

func RequestBusyPoll() {

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
	// logging.Errorf("%s : %v", message, e)
	// panic(e)
	println(message, e)
}

/*
Makes a request to the server
*/

func makeRequest(request string) int {

	println()
	println("Test App - Request: " + request)

	if err := hWR.Write(request, -1); err != nil {
		println("Test App - Write error: ", err)
	}

	response, fd, err := hWR.Read()
	if err != nil {
		println("Test App - Read error: ", err)
	}

	println("Test App - Response: " + response)
	if fd > 0 {
		println("Test App - File Descriptor:", strconv.Itoa(fd))
	}
	println()
	return fd
}

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
