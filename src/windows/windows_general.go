package windows

type WindowsGeneral struct {
}

func (w WindowsGeneral) PathGeneral(path, newPath string) string {
	return path + ";" + newPath
}

func (w WindowsGeneral) PathMapping(path string) string {
	return "%" + path + "%"
}
