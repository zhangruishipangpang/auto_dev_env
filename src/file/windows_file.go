package file

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type CommonFileProcessor struct {
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
	// 打开 zip 文件
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

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
	}

	return nil
}

func (c CommonFileProcessor) Copy(sourcePath string, targetPath string, del bool) (bool, error) {

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
	}

	return false, nil
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
