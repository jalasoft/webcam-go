#include "camera.h"

#include <sstream>
#include <fstream>
#include <vector>
#include <cstdlib>
#include <tuple>
#include <cstring>
#include <sys/ioctl.h>
#include <linux/videodev2.h>
#include <sys/mman.h>
#include <errno.h>

namespace camera_handler {

std::ostream& operator<<(std::ostream& stream, const resolution& res) {
  stream << res.to_string();
  return stream;
}

std::string camera::to_string() const {
  std::stringstream s;
  s << "Camera[" << file_path() << "]";
  return s.str();
}


void camera::check_video_device(int& fd) {

  logger_.info() << "Checking camera capabilities." << std::endl;

  v4l2_capability cap;
  int result = ioctl(fd, VIDIOC_QUERYCAP, &cap);

  if (result < 0) {
    throw camera_error(path_, "Cannot query capabilities.");
  }

  if (!cap.capabilities & V4L2_CAP_VIDEO_CAPTURE) {
    throw camera_error(path_, "Device cannot capture video");
  }

  if (!cap.capabilities & V4L2_CAP_STREAMING) {
    throw camera_error(path_, "Device cannot stream frames");
  }
}

void camera::check_supported_format(int& fd) {
  logger_.info() << "Checking camera supported format." << std::endl;

  v4l2_fmtdesc desc;
  desc.index = 0;
  desc.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;

  while(ioctl(fd, VIDIOC_ENUM_FMT, &desc) >= 0) {
    if (desc.pixelformat == V4L2_PIX_FMT_MJPEG) {
      return;
    }
    desc.index++;
  }
  throw camera_error(path_, "Only V4L2_PIX_FMT_MJPEG is currently suported type of format.");
}

std::vector<resolution*>* camera::load_supported_resolutions(int& fd) {
  
  logger_.info() << "Loading supported resolutions. " << std::endl;

  v4l2_frmsizeenum e;
  e.index = 0;
  e.pixel_format = V4L2_PIX_FMT_MJPEG;

  std::vector<resolution*>* result = new std::vector<resolution*>();

  while(ioctl(fd, VIDIOC_ENUM_FRAMESIZES, &e) >= 0) {
    resolution* r = new resolution(e.discrete.width, e.discrete.height);
    result->push_back(r);
    e.index++;
    logger_.info() << "Suported resolution " << *r << " loaded." << std::endl;
  }
  
  return result; 
}

frame* camera::take_frame_repeatedly(resolution* res) {
  
  while(1) {
    try {
      return take_frame(res);
    } catch(camera_error& err) {
      logger_.info() << "Error occurred during take frame." << std::endl;
    }
    std::this_thread::sleep_for(std::chrono::milliseconds(2000));	
  }
}

frame* camera::take_frame(resolution* res) {
  
  logger_.info() << "Taking frame with " << *res << std::endl;

  file_descriptor file(path_, &logger_);
  int fd = file.descriptor();
  
  
  set_frame_format(fd, res);
  
  request_buffer(fd, 1);
  
  std::tuple<v4l2_buffer, void*> buff_and_ptr = query_buffer(fd, 0);
  
  v4l2_buffer buffer = std::get<0>(buff_and_ptr);
  void* ptr = std::get<1>(buff_and_ptr);

  //---new buffer---
  v4l2_buffer b;
  memset(&b, 0, sizeof(b));
  b.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
  b.memory = V4L2_MEMORY_MMAP;
  b.index = 0;
  //----------------
  
  activate_streaming(fd, b);
  queue_buffer(fd, b);
  dequeue_buffer(fd, b);
 
  char* data = (char*) malloc(buffer.length);
  memcpy(data, (char*) ptr, buffer.length);

  munmap(ptr, buffer.length);

  deactivate_streaming(fd, b);

  logger_.info() << "Frame taken successfully. " << std::endl;

  return new frame(data, buffer.length, res);
}

void camera::set_frame_format(int fd, resolution* res) {
  v4l2_format f;
  f.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
  f.fmt.pix.width = res->width();
  f.fmt.pix.height = res->height();
  f.fmt.pix.pixelformat = V4L2_PIX_FMT_MJPEG;
  f.fmt.pix.field = V4L2_FIELD_NONE;

  if (ioctl(fd, VIDIOC_S_FMT, &f) != 0) {
    throw camera_error(path_, "Cannot specify format of frame.");
  }
}

void camera::request_buffer(int fd, int count) {
  v4l2_requestbuffers b = {0};
  b.count = count;
  b.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
  b.memory = V4L2_MEMORY_MMAP;

  if (ioctl(fd, VIDIOC_REQBUFS, &b) != 0) {
    throw camera_error(path_, "Cannot requst buffer for frame.");
  }
}

std::tuple<v4l2_buffer, void*> camera::query_buffer(int fd, int index) {
  v4l2_buffer b = {0};
  b.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
  b.memory = V4L2_MEMORY_MMAP;
  b.index = index;

  if (ioctl(fd, VIDIOC_QUERYBUF, &b) < 0) {
    throw camera_error(path_, "Cannot allocate buffer for a frame.");
  }

  void* ptr = mmap(NULL, b.length, PROT_READ | PROT_WRITE, MAP_SHARED, fd, b.m.offset);

  if (ptr == MAP_FAILED) {
    throw camera_error(path_, "Cannot map allocated buffer to memory.");
  }

  memset(ptr, 0, b.length);

  return std::make_tuple(b, ptr);
}

void camera::activate_streaming(const int fd, v4l2_buffer buffer) {
  if (ioctl(fd, VIDIOC_STREAMON, &buffer.type) < 0) {
    throw camera_error(path_, "Cannot turn on streaming.");
  }
}

void camera::queue_buffer(const int fd, v4l2_buffer buffer) {

  if (ioctl(fd, VIDIOC_QBUF, &buffer) < 0) {
    throw camera_error(path_, "Cannot queue buffer");
  }
}

void camera::dequeue_buffer(const int fd, v4l2_buffer buffer) {
  if (ioctl(fd, VIDIOC_DQBUF, &buffer) < 0) {
    throw camera_error(path_, "Cannot dequeue buffer.");
  }
}

void camera::deactivate_streaming(const int fd, v4l2_buffer buffer) {

  if (ioctl(fd, VIDIOC_STREAMOFF, &buffer.type) < 0) {
    throw camera_error(path_, "Cannot turn off streaming.");
  }
}
}

