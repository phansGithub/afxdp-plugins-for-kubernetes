package main

import (
	"C"

	"github.com/intel/afxdp-plugins-for-kubernetes/pkg/goclient"
)
import (
	"fmt"

	"github.com/intel/afxdp-plugins-for-kubernetes/internal/uds"
)

func main() {
	// Needed for cgo to generate the .h
}

var closeTime int = 30
var inRoutine bool = false
var cleaners [3]uds.CleanupFunc

/*
GetClientVersion is an exported version for c of goclient's GetClientVersion()
*/
//export GetClientVersion
func GetClientVersion() *C.char {
	return C.CString(goclient.GetClientVersion())
}

/*
ServerVersion is an exported version for c of goclient's GetServerVersion()
*/
//export ServerVersion
func ServerVersion() (*C.char, C.int) {
	response, function, err := goclient.GetServerVersion()
	if err != nil {
		function()
		return C.CString(response), -1
	}

	cleaners[0] = function

	return C.CString(response), 0
}

/*
ServerVersion is an exported version for c of goclient's XskMapFd()
*/
//export XskMapFd
func XskMapFd(device *C.char) (fd, errVal C.int) {
	fdVal, function, err := goclient.RequestXSKmapFD(C.GoString(device))
	fd = C.int(fdVal)
	if err != nil {
		errVal = -1
		function()
		return 0, errVal
	}

	cleaners[1] = function

	return fd, 0
}

/*
ServerVersion is an exported version for c of goclient's RequestBusyPoll()
*/
//export RequestBusyPoll
func RequestBusyPoll(busyTimeout, busyBudget, fd int) C.int {
	function, err := goclient.RequestBusyPoll(busyTimeout, busyBudget, fd)
	if err != nil {
		function()
		return -1
	}
	cleaners[2] = function
	return 0
}

/*
CleanUpConnection an explicit exported cgo function to cleanup a connection after calling any of the other functions.
Pass in one of the available function names to clean up the connection after use.
*/
//export CleanUpConnection
func CleanUpConnection(function *C.char) {
	switch C.GoString(function) {
	case "server_version":
		{
			if cleaners[0] != nil {
				cleaners[0]()
			} else {
				fmt.Println("No available function to call")
				return
			}
		}
	case "xsk_map_fd":
		{
			if cleaners[1] != nil {
				cleaners[1]()
			} else {
				fmt.Println("No available function to call")
				return
			}
		}
	case "busy_poll":
		{
			if cleaners[2] != nil {
				cleaners[2]()
			} else {
				fmt.Println("No available function to call")
				return
			}
		}
	default:
		{
			fmt.Println("Invalid function name, available names are server_version, xsk_map_fd, busy_poll, now returning...")
			return
		}
	}

	fmt.Println(C.GoString(function), " ", "Cleaned up.")
}
