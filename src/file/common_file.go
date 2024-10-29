package file

import (
	"archive/zip"
	"auto_dev_env/src/util"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type CommonFileProcessor struct {
}

func (c CommonFileProcessor) ReadFile(path string) ([]byte, error) {

	fs, err := os.Stat(path)

	if os.IsNotExist(err) {
		return nil, errors.New("path not found --> " + path)
	}

	if err != nil {
		return nil, err
	}

	if fs.IsDir() {
		return nil, errors.New("path is not a file --> " + path)
	}

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func (c CommonFileProcessor) Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (c CommonFileProcessor) UnZip(src, target string) error {

	log.Println("\n开始解压文件：" + src + "\n")

	// 打开 zip 文件
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	total := 0
	for _, f := range r.File {
		// 检查目标路径是否为目录
		if f.FileInfo().IsDir() {
			continue
		}
		total++
	}

	bar := util.GetProgressBar("unZip", total)

	// 创建目标目录
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	// 遍历 zip 文件中的每个文件
	for _, f := range r.File {
		// 构建目标文件的完整路径
		fpath := filepath.Join(target, f.Name)

		// 检查目标路径是否为目录
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, f.Mode()); err != nil {
				return err
			}
			continue
		}

		// 创建目标文件
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}
		dstFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// 打开 zip 文件中的文件
		srcFile, err := f.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// 复制文件内容
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}

		bar.Add(1)
	}

	bar.Finish()
	fmt.Println(" finish.")

	return nil
}

func (c CommonFileProcessor) Copy(sourcePath string, targetPath string, del bool) (bool, error) {

	log.Println("\n 开始复制文件 : " + sourcePath + "\n")

	bar := util.GetProgressBar("Copy", 3)

	stat, err := os.Stat(sourcePath)
	if os.IsNotExist(err) {
		return false, errors.New("sourcePath not found")
	}

	isDir := stat.IsDir()

	if isDir {
		err := c.copyDir(sourcePath, targetPath)
		if err != nil {
			return false, err
		}
	} else {
		err := c.copyFile(sourcePath, targetPath)
		if err != nil {
			return false, err
		}
	}
	bar.Add(1)

	if del {
		if isDir {
			err := os.RemoveAll(sourcePath)
			if err != nil {
				return false, err
			}
		} else {
			err := os.Remove(sourcePath)
			if err != nil {
				return false, err
			}
		}
		bar.Add(1)
	}

	bar.Add(1)
	bar.Finish()
	fmt.Println(" finish.")

	return true, nil
}

// copyFile 复制单个文件
func (c CommonFileProcessor) copyFile(src, dst string) error {

	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// 复制文件权限
	return os.Chmod(dst, 0644)
}

// copyDir 复制整个目录
func (c CommonFileProcessor) copyDir(src, dst string) error {
	// 读取源目录
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 创建目标目录
	err = os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	// 递归复制每个条目
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = c.copyDir(srcPath, dstPath)
		} else {
			err = c.copyFile(srcPath, dstPath)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
