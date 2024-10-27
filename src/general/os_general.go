package general

type OsGeneral interface {
	PathGeneral(path string) string
}

type WindowsGeneral struct {
}

func (w WindowsGeneral) PathGeneral(path string) string {
	return ";" + path
}
