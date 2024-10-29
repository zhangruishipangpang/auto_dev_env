package general

type OsGeneral interface {
	PathGeneral(path, newPath string) string

	PathMapping(path string) string
}

type WindowsGeneral struct {
}

func (w WindowsGeneral) PathGeneral(path, newPath string) string {
	return path + ";" + newPath
}

func (w WindowsGeneral) PathMapping(path string) string {
	return "%" + path + "%"
}
