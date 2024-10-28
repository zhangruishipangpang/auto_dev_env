package general

type OsGeneral interface {
	PathGeneral(path, newPath string) string
}

type WindowsGeneral struct {
}

func (w WindowsGeneral) PathGeneral(path, newPath string) string {
	return path + ";" + "%" + newPath + "%"
}

type LinuxCentOSGeneral struct {
}

func (l LinuxCentOSGeneral) PathGeneral(path, newPath string) string {
	return path + ":" + "%" + newPath + "%"
}
