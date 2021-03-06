package ioctl

import (
	"errors"
	"fmt"
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

	VIDIOC_QUERYCAP        = (IOC_READ << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (0 << IOC_NR_SHIFT) | ((unsafe.Sizeof(v4l2.V4l2Capability{})) << IOC_SIZE_SHIFT)
	VIDIOC_ENUM_FMT        = ((IOC_READ | IOC_WRITE) << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (2 << IOC_NR_SHIFT) | ((unsafe.Sizeof(v4l2.V4l2Fmtdesc{})) << IOC_SIZE_SHIFT)
	VIDIOC_ENUM_FRAMESIZES = ((IOC_READ | IOC_WRITE) << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (74 << IOC_NR_SHIFT) | ((unsafe.Sizeof(v4l2.V4l2Frmsizeenum{})) << IOC_SIZE_SHIFT)
	VIDIOC_S_FMT           = ((IOC_READ | IOC_WRITE) << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (5 << IOC_NR_SHIFT) | (unsafe.Sizeof(v4l2.V4l2Format{}) << IOC_SIZE_SHIFT)
	VIDIOC_REQBUFS         = ((IOC_READ | IOC_WRITE) << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (8 << IOC_NR_SHIFT) | ((unsafe.Sizeof(v4l2.V4l2RequestBuffers{})) << IOC_SIZE_SHIFT)
	VIDIOC_QUERYBUF        = ((IOC_READ | IOC_WRITE) << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (9 << IOC_NR_SHIFT) | ((unsafe.Sizeof(v4l2.V4l2Buffer{})) << IOC_SIZE_SHIFT)
	VIDIOC_STREAMON        = (IOC_WRITE << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (18 << IOC_NR_SHIFT) | (unsafe.Sizeof(uint32(0)) << IOC_SIZE_SHIFT)
	VIDIOC_STREAMOFF       = (IOC_WRITE << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (19 << IOC_NR_SHIFT) | (unsafe.Sizeof(uint32(0)) << IOC_SIZE_SHIFT)
	VIDIOC_DQBUF           = ((IOC_READ | IOC_WRITE) << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (17 << IOC_NR_SHIFT) | (unsafe.Sizeof(v4l2.V4l2Buffer{}) << IOC_SIZE_SHIFT)
	VIDIOC_QBUF            = ((IOC_READ | IOC_WRITE) << IOC_DIR_SHIFT) | (uintptr('V') << IOC_TYPE_SHIFT) | (15 << IOC_NR_SHIFT) | ((unsafe.Sizeof(v4l2.V4l2Buffer{})) << IOC_SIZE_SHIFT)
)

func QueryCapability(fd uintptr) (v4l2.V4l2Capability, error) {
	capability := v4l2.V4l2Capability{}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, VIDIOC_QUERYCAP, uintptr(unsafe.Pointer(&capability)))

	if err != 0 {
		return capability, err
	}

	return capability, nil
}

func QueryFormat(fd uintptr, desc *v4l2.V4l2Fmtdesc) (bool, error) {

	r1, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, VIDIOC_ENUM_FMT, uintptr(unsafe.Pointer(desc)))

	if err != 0 {
		return false, err
	}

	if r1 == 0 {
		return false, nil
	}

	return true, nil
}

func QueryFrameSize(fd uintptr, str *v4l2.V4l2Frmsizeenum) (bool, error) {

	r1, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, VIDIOC_ENUM_FRAMESIZES, uintptr(unsafe.Pointer(str)))

	if r1 > 0 {
		return false, nil
	}

	if err != 0 {
		return false, err
	}

	return true, nil
}

func SetFrameSize(fd uintptr, str *v4l2.V4l2Format) error {

	r1, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, VIDIOC_S_FMT, uintptr(unsafe.Pointer(str)))

	if r1 > 0 {
		return errors.New(fmt.Sprintf("Cannot set frame size, ioctl system call returned status %v", r1))
	}

	if err != 0 {
		return err
	}

	return nil
}

func RequestBuffer(fd uintptr, str *v4l2.V4l2RequestBuffers) error {

	r1, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, VIDIOC_REQBUFS, uintptr(unsafe.Pointer(str)))

	if err != 0 {
		return err
	}

	if r1 != 0 {
		return errors.New(fmt.Sprintf("Cannot request buffer, ioctl system call returned status %v", r1))
	}

	return nil
}

func QueryBuffer(fd uintptr, buffer *v4l2.V4l2Buffer) error {

	r1, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, VIDIOC_QUERYBUF, uintptr(unsafe.Pointer(buffer)))

	if err != 0 {
		return err
	}

	if r1 != 0 {
		return errors.New(fmt.Sprintf("Cannot query buffer, ioctl system call returned status %v", r1))
	}

	return nil
}

func ActivateStreaming(fd uintptr, bufType uint32) error {

	r1, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, VIDIOC_STREAMON, uintptr(unsafe.Pointer(&bufType)))

	if err != 0 {
		return err
	}

	if r1 != 0 {
		return errors.New(fmt.Sprintf("Cannot activate streaming, ioctl system call returned with status %d\n", r1))
	}

	return nil
}

func DeactivateStreaming(fd uintptr, bufType uint32) error {

	r1, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, VIDIOC_STREAMOFF, uintptr(unsafe.Pointer(&bufType)))

	if err != 0 {
		return err
	}

	if r1 != 0 {
		return errors.New(fmt.Sprintf("Cannot deactivate streaming, ioctl system call returned with status %d\n", r1))
	}

	return nil
}

func QueueBuffer(fd uintptr, buffer *v4l2.V4l2Buffer) error {

	r1, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, VIDIOC_QBUF, uintptr(unsafe.Pointer(buffer)))

	if r1 != 0 {
		return errors.New(fmt.Sprintf("Cannot queue buffer, ioctl system call returned with status %d\n", r1))
	}

	if err != 0 {
		return err
	}

	return nil
}

func DequeueBuffer(fd uintptr, buffer *v4l2.V4l2Buffer) error {

	r1, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, VIDIOC_DQBUF, uintptr(unsafe.Pointer(buffer)))

	if r1 != 0 {
		return errors.New(fmt.Sprintf("Cannot dequeue buffer, ioctl system call returned with status %d\n", r1))
	}

	if err != 0 {
		return err
	}

	return nil
}
