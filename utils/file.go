package utils

import (
	"archive/zip"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// ApiRoot 当前项目根目录
var ApiRoot string

// GetPath 获取项目路径
func GetPath() string {
	if ApiRoot != "" {
		return ApiRoot
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		print(err.Error())
	}
	ApiRoot = strings.Replace(dir, "\\", "/", -1)
	return ApiRoot
}

// CheckDir 判断文件目录否存在
func CheckDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

// MkdirDir 创建文件夹,支持x/a/a  多层级
func MkdirDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// RemoveDir 删除文件
func RemoveDir(filePath string) error {
	return os.RemoveAll(filePath)
}

// CheckFile 判断文件是否存在  存在返回 true 不存在返回false
func CheckFile(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetFileSize 获取文件大小
func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	size := fileInfo.Size()
	return size, nil
}

// UploadFile 创建文件并生产目录
func UploadFile(file *multipart.FileHeader, path string) (string, error) {
	if reflect.ValueOf(file).IsNil() || !reflect.ValueOf(file).IsValid() {
		return "", errors.New("invalid memory address or nil pointer dereference")
	}
	src, err := file.Open()
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			log.Println(err)
		}
	}(src)
	if err != nil {
		return "", err
	}

	err = MkdirDir(path)
	if err != nil {
		return "", err
	}
	filename := strings.Replace(file.Filename, " ", "", -1)
	filename = strings.Replace(filename, "\n", "", -1)
	dst, err := os.Create(path + filename)
	if err != nil {
		return "", err
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			log.Println(err)
		}
	}(dst)

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	return filename, nil
}

// Zip 文件压缩
func Zip(srcFile string, destZip string) error {
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+"/")
		// header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return err
	}

	return err
}

// Unzip 文件解压
func Unzip(zipFile, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
