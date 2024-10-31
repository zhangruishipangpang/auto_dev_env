package mac

type MacGeneral struct {
}

func (w MacGeneral) PathGeneral(path, newPath string) string {
	return path + ":" + newPath
}

func (w MacGeneral) PathMapping(path string) string {
	return "%" + path + "%"
}
