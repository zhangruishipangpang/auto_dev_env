package util

import (
	"github.com/schollz/progressbar/v3"
)

func GetProgressBar(prefix string, total int) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(total,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowCount(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetDescription("[cyan]"+prefix+"[reset]"),
		progressbar.OptionSetTheme(progressbar.ThemeDefault),
	)

	return bar
}
