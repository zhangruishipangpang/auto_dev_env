package linux

type LinuxGeneral struct {
}

func (l LinuxGeneral) PathGeneral(path, newPath string) string {
	return path + ":" + newPath
}

func (l LinuxGeneral) PathMapping(path string) string {
	return "$" + path
}
