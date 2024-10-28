package env

var pathNews []string = make([]string, 0)

func addPathStore(pathNew string) {

	pathNews = append(pathNews, pathNew)
}

func getNeedAddPaths() []string {
	if len(pathNews) == 0 {
		return nil
	}
	return pathNews
}
