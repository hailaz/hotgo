// Package file
package file

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/os/gfile"

	"hotgo/utility/format"
)

// 文件信息
type fileInfo struct {
	name string
	size int64
}

// WalkDir 递归获取目录下文件的名称和大小
func WalkDir(dirname string) ([]fileInfo, error) {
	op, err := filepath.Abs(dirname) // 获取目录的绝对路径
	if nil != err {
		return nil, err
	}
	files, err := os.ReadDir(op) // 获取目录下所有文件的信息，包括文件和文件夹
	if nil != err {
		return nil, err
	}

	var fileInfos []fileInfo // 返回值，存储读取的文件信息
	for _, f := range files {
		if f.IsDir() { // 如果是目录，那么就递归调用
			fs, err := WalkDir(op + `/` + f.Name()) // 路径分隔符，linux 和 windows 不同
			if nil != err {
				return nil, err
			}
			fileInfos = append(fileInfos, fs...) // 将 slice 添加到 slice
		} else {
			info, err := f.Info()
			if nil != err {
				return nil, err
			}
			fi := fileInfo{op + `/` + f.Name(), info.Size()}
			fileInfos = append(fileInfos, fi) // slice 中添加成员
		}
	}
	return fileInfos, nil
}

// DirSize 获取目录下所有文件大小
func DirSize(dirname string) string {
	var (
		ss       int64
		files, _ = WalkDir(dirname)
	)
	for _, n := range files {
		ss += n.size
	}
	return format.FileSize(ss)
}

func MergeAbs(path string, fileName ...string) string {
	var paths = []string{gfile.RealPath(path)}
	paths = append(paths, fileName...)
	return gfile.Join(paths...)
}

// NewMultipartFileHeader 创建一个MultipartFileHeader
func NewMultipartFileHeader(filename string, b []byte) (*multipart.FileHeader, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, bytes.NewReader(b)); err != nil {
		return nil, err
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(int64(len(b)) + 1024)
	if err != nil {
		return nil, err
	}
	defer func() { _ = form.RemoveAll() }()

	if files, ok := form.File["file"]; ok && len(files) > 0 {
		return files[0], nil
	}
	return nil, io.ErrUnexpectedEOF
}
