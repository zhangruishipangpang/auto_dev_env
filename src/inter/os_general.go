package inter

type GenOsGeneral interface {
	PathGeneral(path, newPath string) string

	PathMapping(path string) string
}
