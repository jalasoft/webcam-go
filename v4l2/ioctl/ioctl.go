package ioctl

import (
	"syscall"
	"unsafe"
	"v4l2"
)

const (
	IOC_NR_BITS   = 8
	IOC_TYPE_BITS = 8
	IOC_SIZE_BITS = 14

	IOC_NR_SHIFT = 0

	IOC_READ  = 2
	IOC_WRITE = 1

	IOC_TYPE_SHIFT = IOC_NR_SHIFT + IOC_NR_BITS
	IOC_SIZE_SHIFT = IOC_TYPE_SHIFT + IOC_TYPE_BITS
	IOC_DIR_SHIFT  = IOC_SIZE_SHIFT + IOC_SIZE_BITS

	VIDIOC_QUERYCAP = (IOC_READ << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (0 << IOC_NR_SHIFT) | ((unsafe.Sizeof(v4l2.V4l2Capability{})) << IOC_SIZE_SHIFT)
	VIDIOC_ENUM_FMT        = ((IOC_READ | IOC_WRITE) << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (2 << IOC_NR_SHIFT) | ((unsafe.Sizeof(v4l2_fmtdesc{})) << IOC_SIZE_SHIFT)
	//VIDIOC_ENUM_FRAMESIZES = ((IOC_READ | IOC_WRITE) << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (74 << IOC_NR_SHIFT) | ((unsafe.Sizeof(v4l2_frmsizeenum{})) << IOC_SIZE_SHIFT)
	//VIDIOC_S_FMT           = ((IOC_READ | IOC_WRITE) << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (5 << IOC_NR_SHIFT) | (204 << IOC_SIZE_SHIFT)
	//VIDIOC_REQBUFS         = ((IOC_READ | IOC_WRITE) << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (8 << IOC_NR_SHIFT) | ((unsafe.Sizeof(V4l2RequestBuffers{})) << IOC_SIZE_SHIFT)
)

func QueryCapability(fd uintptr) (v4l2.V4l2Capability, error) {
	capability := v4l2.V4l2Capability{}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, VIDIOC_QUERYCAP, uintptr(unsafe.Pointer(&capability)))

	if err != 0 {
		return capability, err
	}

	return capability, nil
}

func QueryFormat(fd uintptr, desc* V4l2Fmtdesc) bool, error {

	r1, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, VIDIOC_ENUM_FMT, uintptr(unsafe.Pointer(desc)))

	if err != 0 {
		return false, err
	}

	if r1 == 0 {
		return false, nil
	}

	return true, nil
}
