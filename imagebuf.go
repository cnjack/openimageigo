package oiio

/*
#include "stdlib.h"

#include "cpp/oiio.h"

*/
import "C"

import (
	"runtime"
	"unsafe"
	"errors"
)

// Description of where the pixels live for this ImageBuf
type IBStorage int

const (
	IBStorageLocalBuffer IBStorage = C.IBSTORAGE_LOCALBUFFER
	IBStorageAppBuffer   IBStorage = C.IBSTORAGE_APPBUFFER
	IBStorageImageCache  IBStorage = C.IBSTORAGE_IMAGECACHE
	IBStorageUninitialized IBStorage = C.IBSTORAGE_UNINITIALIZED
)

// An ImageBuf is a simple in-memory representation of a 2D image. 
// It uses ImageInput and ImageOutput underneath for its file I/O, and has simple 
// routines for setting and getting individual pixels, that hides most of the details 
// of memory layout and data representation (translating to/from float automatically).
type ImageBuf struct {
	ptr unsafe.Pointer
}

func newImageBuf(i unsafe.Pointer) *ImageBuf {
	spec := new(ImageBuf)
	spec.ptr = i
	runtime.SetFinalizer(spec, deleteImageBuf)
	return spec
}

func deleteImageBuf(i *ImageBuf) {
	if i.ptr != nil {
		C.free(i.ptr)
		i.ptr = nil
	}
}

// Return the last error generated by API calls.
// An nil error will be returned if no error has occured.
func (i *ImageBuf) LastError() error {
	err := C.GoString(C.ImageBuf_geterror(i.ptr))
	if err == "" {
		return nil
	}
	return errors.New(err)
}

// Construct an empty/uninitialized ImageBuf. This is relatively useless until you call reset().
func NewImageBuf() *ImageBuf {
	buf := C.ImageBuf_New()
	return newImageBuf(buf)
}

// Construct an ImageBuf to read the named image – but don't actually read it yet! 
// The image will actually be read when other methods need to access the spec and/or pixels, 
// or when an explicit call to init_spec() or read() is made, whichever comes first. 
// Uses the global/shared ImageCache.
func NewImageBufPath(path string) (*ImageBuf, error) {
	c_str := C.CString(path)	
	defer C.free(unsafe.Pointer(c_str))

	buf := newImageBuf(C.ImageBuf_New_WithCache(c_str, nil))
	err := buf.LastError()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// Is this ImageBuf object initialized?
func (i *ImageBuf) Initialized() bool {
	return bool(C.ImageBuf_initialized(i.ptr))
}

// Restore the ImageBuf to an uninitialized state.
func (i *ImageBuf) Clear() {
	C.ImageBuf_clear(i.ptr)
}

func (i *ImageBuf) Storage() IBStorage {
	return IBStorage(C.ImageBuf_storage(i.ptr))
}


// Read the file from disk. Generally will skip the read if we've already got a current 
// version of the image in memory, unless force==true. 
// This uses ImageInput underneath, so will read any file format for which an appropriate 
// imageio plugin can be found. 
func (i *ImageBuf) Read(force bool) error {
	return i.ReadFormatCallback(force, TypeUnknown, nil)
}

// Read the file from disk. Generally will skip the read if we've already got a current 
// version of the image in memory, unless force==true. 
// This uses ImageInput underneath, so will read any file format for which an appropriate 
// imageio plugin can be found. 
// 
// This call supports passing a callback pointer to both track the progress,
// and to optionally abort the processing. The callback function will receive
// a float32 value indicating the percentage done of the processing, and should
// return true if the process should abort, and false if it should continue.
// 
func (i *ImageBuf) ReadCallback(force bool, progress *ProgressCallback) error {
	return i.ReadFormatCallback(force, TypeUnknown, progress)
}

// Read the file from disk. Generally will skip the read if we've already got a current 
// version of the image in memory, unless force==true. 
// This uses ImageInput underneath, so will read any file format for which an appropriate 
// imageio plugin can be found. 
// 
// Specify a specific conversion format or TypeUnknown for automatic handling.
// 
// This call supports passing a callback pointer to both track the progress,
// and to optionally abort the processing. The callback function will receive
// a float32 value indicating the percentage done of the processing, and should
// return true if the process should abort, and false if it should continue.
// 
func (i *ImageBuf) ReadFormatCallback(force bool, convert TypeDesc, progress *ProgressCallback) error {
	var cbk unsafe.Pointer
	if progress != nil {
		cbk = unsafe.Pointer(progress)
	}

	ok := C.ImageBuf_read(i.ptr, 0, 0, C.bool(force), C.TypeDesc(convert), cbk)
	if !bool(ok) {
		return i.LastError()
	}

	return nil
}

// Return the name of this image.
func (i *ImageBuf) Name() string {
	return C.GoString(C.ImageBuf_name(i.ptr))
}

// Return the name of the image file format of the disk file we read into this image. 
// Returns an empty string if this image was not the result of a Read().
func (i *ImageBuf) FileFormatName() string {
	return C.GoString(C.ImageBuf_file_format_name(i.ptr))
}

// Return the index of the subimage are we currently viewing
func (i *ImageBuf) SubImage() int {
	return int(C.ImageBuf_subimage(i.ptr))
}

// Return the number of subimages in the file.
func (i *ImageBuf) NumSubImages() int {
	return int(C.ImageBuf_nsubimages(i.ptr))
}

// Return the index of the miplevel are we currently viewing
func (i *ImageBuf) MipLevel() int {
	return int(C.ImageBuf_miplevel(i.ptr))
}

// Return the number of miplevels of the current subimage.
func (i *ImageBuf) NumMipLevels() int {
	return int(C.ImageBuf_nmiplevels(i.ptr))
}

// Return the number of color channels in the image.
func (i *ImageBuf) NumChannels() int {
	return int(C.ImageBuf_nchannels(i.ptr))
}

func (i *ImageBuf) Orientation() int {
	return int(C.ImageBuf_orientation(i.ptr))
}

func (i *ImageBuf) OrientedWidth() int {
	return int(C.ImageBuf_oriented_width(i.ptr))
}

func (i *ImageBuf) OrientedHeight() int {
	return int(C.ImageBuf_oriented_height(i.ptr))
}

func (i *ImageBuf) OrientedX() int {
	return int(C.ImageBuf_oriented_x(i.ptr))
}

func (i *ImageBuf) OrientedY() int {
	return int(C.ImageBuf_oriented_y(i.ptr))
}

func (i *ImageBuf) OrientedFullWidth() int {
	return int(C.ImageBuf_oriented_full_width(i.ptr))
}

func (i *ImageBuf) OrientedFullHeight() int {
	return int(C.ImageBuf_oriented_full_height(i.ptr))
}

func (i *ImageBuf) OrientedFullX() int {
	return int(C.ImageBuf_oriented_full_x(i.ptr))
}

func (i *ImageBuf) OrientedFullY() int {
	return int(C.ImageBuf_oriented_full_y(i.ptr))
}

// Return the beginning (minimum) x coordinate of the defined image.
func (i *ImageBuf) XBegin() int {
	return int(C.ImageBuf_xbegin(i.ptr))
}

// Return the end (one past maximum) x coordinate of the defined image.
func (i *ImageBuf) XEnd() int {
	return int(C.ImageBuf_xend(i.ptr))
}

// Return the beginning (minimum) y coordinate of the defined image
func (i *ImageBuf) YBegin() int {
	return int(C.ImageBuf_ybegin(i.ptr))
}

// Return the end (one past maximum) y coordinate of the defined image. 
func (i *ImageBuf) YEnd() int {
	return int(C.ImageBuf_yend(i.ptr))
}

// Return the beginning (minimum) z coordinate of the defined image.
func (i *ImageBuf) ZBegin() int {
	return int(C.ImageBuf_zbegin(i.ptr))
}

// Return the end (one past maximum) z coordinate of the defined image.
func (i *ImageBuf) ZEnd() int {
	return int(C.ImageBuf_zend(i.ptr))
}

// Return the end (one past maximum) z coordinate of the defined image.
func (i *ImageBuf) XMin() int {
	return int(C.ImageBuf_xmin(i.ptr))
}

// Return the maximum x coordinate of the defined image.
func (i *ImageBuf) XMax() int {
	return int(C.ImageBuf_xmax(i.ptr))
}

// Return the minimum y coordinate of the defined image.
func (i *ImageBuf) YMin() int {
	return int(C.ImageBuf_ymin(i.ptr))
}

// Return the maximum y coordinate of the defined image.
func (i *ImageBuf) YMax() int {
	return int(C.ImageBuf_ymax(i.ptr))
}

// Return the minimum z coordinate of the defined image.
func (i *ImageBuf) ZMin() int {
	return int(C.ImageBuf_zmin(i.ptr))
}

// Return the maximum z coordinate of the defined image. 
func (i *ImageBuf) ZMax() int {
	return int(C.ImageBuf_zmax(i.ptr))
}




