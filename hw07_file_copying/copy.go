package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()
	srcinfo, err := src.Stat()
	if err != nil {
		return err
	}
	// fmt.Printf("%T\n%+v\n", srcinfo.Size(), srcinfo.Size())
	lenSrcFile := srcinfo.Size()
	if offset > 0 {
		if offset > lenSrcFile {
			return ErrOffsetExceedsFileSize
		}
		src.Seek(offset, io.SeekStart)
		lenSrcFile -= offset
	}
	switch {
	case lenSrcFile == 0:
		return ErrUnsupportedFile
	case limit == 0 || limit > lenSrcFile:
		limit = lenSrcFile
	}
	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	if err != nil {
		return err
	}
	bar := pb.New(int(limit)).SetUnits(pb.U_BYTES)
	bar.Start()
	_, err = io.CopyN(dst, bar.NewProxyReader(src), limit)
	bar.Finish()
	// fmt.Printf("written %+v bytes\n", written)

	return err
}
