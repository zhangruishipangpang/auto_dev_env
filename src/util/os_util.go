package util

import "runtime"

func GetCurrentOs() string {

	return runtime.GOOS
}
