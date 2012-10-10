package client

import (
	"fmt"
	""
)

var (
	MASTER_SERVER = ""
)

func SetServer(master string) {
	MASTER_SERVER = string
}

type GoPHSFile struct {
	name string

}

func Create(name string) (file *GoPHSFile, err error) {

}

func Open(name String) (file *GoPHSFile, err error) {

}

func (f *GoPHSFile) Close() error {

}

func (f *GoPHSFile) Name() string {

}

func (f *GoPHSFile) Read(b []byte) (n int, err error) {

}

func (f *GoPHSFile) ReadAt(b []byte, off int64) (n int, err error) {

}

func (f *GoPHSFile) Seek(offset int64, whence int) (ret int64, err error) {

}

func (f *GoPHSFile) Write(b []byte) (n int, err error) {

}

func (f *GoPHSFile) WriteAt(b []byte, off int64) {

}

func (f *GoPHSFile) WriteString(s string) (ret int, err error) {

}