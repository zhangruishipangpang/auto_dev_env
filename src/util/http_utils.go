package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func GetInternetResourceToLocal(resourceUrl, localPath, downloadFileName string) error {

	if exist, _ := checkDirExist(localPath); !exist {
		if err := os.Mkdir(localPath, os.ModeDir); err != nil {
			panic(err)
		}
	}

	response, err := http.Get(resourceUrl)
	if err != nil {
		return err
	}

	newDirPath := filepath.Join(localPath, downloadFileName)

	dirFile, err := os.Create(newDirPath)
	if err != nil {
		return err
	}
	defer dirFile.Close()

	// 获取文件大小
	contentLength := response.ContentLength
	if contentLength <= 0 {
		contentLength = -1
	}

	fmt.Println("contentLength===>", contentLength)

	// 读取响应体并写入文件
	buf := make([]byte, 32*1024)
	for {
		n, err := response.Body.Read(buf)
		if n > 0 {
			if _, err := dirFile.Write(buf[:n]); err != nil {
				return err
			}
		}
		if err == io.EOF {
			fmt.Println("eof break")
			break
		}
		if err != nil {
			return err
		}
	}

	fmt.Println("Download complete!")
	return nil
}

func checkDirExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
