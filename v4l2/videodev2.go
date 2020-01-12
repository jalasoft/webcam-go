package v4l2

type V4l2Capability struct {
	Driver       [16]uint8
	Card         [32]uint8
	BusInfo      [32]uint8
	Version      uint32
	Capabilities uint32
	DeviceCaps   uint32
	Reserved     [3]uint32
}

/* Values for 'capabilities' field */
const (
	V4L2_CAP_VIDEO_CAPTURE        = 0x00000001 /* Is a video capture device */
	V4L2_CAP_VIDEO_OUTPUT         = 0x00000002 /* Is a video output device */
	V4L2_CAP_VIDEO_OVERLAY        = 0x00000004 /* Can do video overlay */
	V4L2_CAP_VBI_CAPTURE          = 0x00000010 /* Is a raw VBI capture device */
	V4L2_CAP_VBI_OUTPUT           = 0x00000020 /* Is a raw VBI output device */
	V4L2_CAP_SLICED_VBI_CAPTURE   = 0x00000040 /* Is a sliced VBI capture device */
	V4L2_CAP_SLICED_VBI_OUTPUT    = 0x00000080 /* Is a sliced VBI output device */
	V4L2_CAP_RDS_CAPTURE          = 0x00000100 /* RDS data capture */
	V4L2_CAP_VIDEO_OUTPUT_OVERLAY = 0x00000200 /* Can do video output overlay */
	V4L2_CAP_HW_FREQ_SEEK         = 0x00000400 /* Can do hardware frequency seek  */
	V4L2_CAP_RDS_OUTPUT           = 0x00000800 /* Is an RDS encoder */

	/* Is a video capture device that supports multiplanar formats */
	V4L2_CAP_VIDEO_CAPTURE_MPLANE = 0x00001000
	/* Is a video output device that supports multiplanar formats */
	V4L2_CAP_VIDEO_OUTPUT_MPLANE = 0x00002000
	/* Is a video mem-to-mem device that supports multiplanar formats */
	V4L2_CAP_VIDEO_M2M_MPLANE = 0x00004000
	/* Is a video mem-to-mem device */
	V4L2_CAP_VIDEO_M2M = 0x00008000

	V4L2_CAP_TUNER     = 0x00010000 /* has a tuner */
	V4L2_CAP_AUDIO     = 0x00020000 /* has audio support */
	V4L2_CAP_RADIO     = 0x00040000 /* is a radio device */
	V4L2_CAP_MODULATOR = 0x00080000 /* has a modulator */

	V4L2_CAP_SDR_CAPTURE    = 0x00100000 /* Is a SDR capture device */
	V4L2_CAP_EXT_PIX_FORMAT = 0x00200000 /* Supports the extended pixel format */
	V4L2_CAP_SDR_OUTPUT     = 0x00400000 /* Is a SDR output device */

	V4L2_CAP_READWRITE = 0x01000000 /* read/write systemcalls */
	V4L2_CAP_ASYNCIO   = 0x02000000 /* async I/O */
	V4L2_CAP_STREAMING = 0x04000000 /* streaming I/O ioctls */

	V4L2_CAP_TOUCH = 0x10000000 /* Is a touch device */

	V4L2_CAP_DEVICE_CAPS = 0x80000000 /* sets device capabilities field */
)

/*
 *	F O R M A T   E N U M E R A T I O N
 */
type V4l2Fmtdesc struct {
	Index       uint32 /* Format number      */
	Typ         uint32 /* enum v4l2_buf_type */
	Flags       uint32
	Description [32]uint8 /* Description string */
	Pixelformat uint32    /* Format fourcc      */
	Reserved    [4]uint32
}

const (
	V4L2_BUF_TYPE_VIDEO_CAPTURE        = 1
	V4L2_BUF_TYPE_VIDEO_OUTPUT         = 2
	V4L2_BUF_TYPE_VIDEO_OVERLAY        = 3
	V4L2_BUF_TYPE_VBI_CAPTURE          = 4
	V4L2_BUF_TYPE_VBI_OUTPUT           = 5
	V4L2_BUF_TYPE_SLICED_VBI_CAPTURE   = 6
	V4L2_BUF_TYPE_SLICED_VBI_OUTPUT    = 7
	V4L2_BUF_TYPE_VIDEO_OUTPUT_OVERLAY = 8
	V4L2_BUF_TYPE_VIDEO_CAPTURE_MPLANE = 9
	V4L2_BUF_TYPE_VIDEO_OUTPUT_MPLANE  = 10
	V4L2_BUF_TYPE_SDR_CAPTURE          = 11
	V4L2_BUF_TYPE_SDR_OUTPUT           = 12
	/* Deprecated, do not use */
	V4L2_BUF_TYPE_PRIVATE = 0x80
)

func v4l2_fourcc(a uint32, b uint32, c uint32, d uint32) uint32 {
	return (a | (b << 8) | (c << 16) | (d << 24))
}

/*      Pixel format         FOURCC                          depth  Description  */

/* RGB formats */
var V4L2_PIX_FMT_RGB332 uint32 = v4l2_fourcc('R', 'G', 'B', '1') /*  8  RGB-3-3-2     */
