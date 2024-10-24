package file

// Processor FP(File Processor) 处理程序
type Processor interface {
	// Exist 检测文件路径是否存在
	Exist(path string) (bool, error)
	// UnZip 解压缩文件
	// file: 需要解压的文件
	// target: 目标文件的全路径与名称
	UnZip(file, target string) error
	// Copy 复制文件或文件夹
	// del: 是否同时删除sourcePath文件
	Copy(sourcePath string, targetPath string, del bool) (bool, error)
}
