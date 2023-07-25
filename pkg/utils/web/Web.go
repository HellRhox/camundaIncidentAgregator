package web

import (
	"camundaIncidentAggregator/pkg/utils/constants"
	"fmt"
	"github.com/charmbracelet/log"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) {
	var err error
	url += constants.AppPath
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.With(err).Error("Error Opening Link")
	}
}
