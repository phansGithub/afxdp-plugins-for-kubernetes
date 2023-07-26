#Compiling the C Library

>NOTE: You may have to run ```go get github.com/intel/afxdp-plugins-for-kubernetes/pkg/goclient@93d43de```

##Dynamically compiling the library into a test-app:


```go build -o lib_udsclient.so -buildmode=c-shared cclient.go```
	- this creates .h and .so files from the cclient.go file.

```cp lib_udsclient.so /lib```
	- copy the .so to the lib directory.
  >NOTE: we could also skip this step and explicitly point gcc to the shared library when compiling the test app like so:
  
	- ```gcc main.c ../lib_udsclient.so -W```

```cp lib_udsclient.h /< Dynamic folder name >```
	- copy the .h file to the test app directory

```cd < Dynamic folder name >```

```g++ -fPIC -Llib/ -Wall -o testApp main.c -l_udsclient```
	- build the test app
	- results in testApp binary

```./testApp```


##Statically compiling the library into a test-app:

```go build -o lib_udsclient.a -buildmode=c-archive ./cclient.go```
	- this creates .h and .a files

```cp lib_udsclient.a <Static folder name>```
```cp lib_udsclient.h <Static folder name>```
	- copy those files to the test app

```gcc -c main.c -o testApp.o```
	- compile test app
	- results in a .o file

```gcc -pthread -o testApp testApp.o -L. -l_udsclient```
	- link compiled test app to static library
	- results in testApp binary

```./testApp```


## Next steps

After creating the .h and/or the .a and the .o files as described in the previous steps,
import the generated cgo functions into your c file like so:

``` 
#include <stdio.h>
#include "lib_udsclient.h"

extern char* GetClientVersion();
extern struct ServerVersion_return ServerVersion();
extern struct XskMapFd_return XskMapFd(char* device);
extern int RequestBusyPoll(GoInt busyTimeout, GoInt busyBudget, GoInt fd);
extern void CleanUpConnection(char* function);

int main(void) {
	// Call functions here...
}
```