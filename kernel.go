package opencl

/*
#ifdef __APPLE__
#include <OpenCL/opencl.h>
#else
#include <CL/cl.h>
#endif

#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Kernel struct {
	kernel C.cl_kernel
}

func createKernel(program Program, name string) (*Kernel, error) {
	clName := C.CString(name)
	defer C.free(unsafe.Pointer(clName))
	var err C.cl_int
	kernel := C.clCreateKernel(program.program, clName, &err)
	if err != C.CL_SUCCESS {
		return nil, fmt.Errorf("failed to create kernel: %d", err)
	}
	return &Kernel{kernel}, nil
}

func (k Kernel) Release() error {
	err := C.clReleaseKernel(k.kernel)
	if err != C.CL_SUCCESS {
		return fmt.Errorf("failed to release kernel: %d", err)
	}
	return nil
}

func (k Kernel) SetArg(index uint, size uintptr, value unsafe.Pointer) error {
	err := C.clSetKernelArg(k.kernel, C.cl_uint(index), C.size_t(size), value)
	if err != C.CL_SUCCESS {
		return fmt.Errorf("error setting kernel argument: %d", err)
	}
	return nil
}

func (k Kernel) SetArgBuffer(index uint, buffer Memory) error {
	return k.SetArg(index, unsafe.Sizeof(buffer.memory), unsafe.Pointer(&buffer.memory))
}
