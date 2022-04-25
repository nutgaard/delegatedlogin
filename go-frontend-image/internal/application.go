package internal

import (
	. "frontend-image/internal/config"
	copy_process "frontend-image/internal/copy-process"
	"frontend-image/internal/router"
	"frontend-image/internal/tmpl"
	"github.com/rs/zerolog/log"
)

func StartApplication() {
	appData := LoadAndCreateServices()
	log.Printf("Moving static resource to tmp folder in order to change them")
	copy_process.CopyAndProcessFiles("./www", "/tmp/www", func(content []byte) []byte {
		return []byte(tmpl.ReplaceVariableReferences(string(content), nil))
	})

	log.Printf("Starting application: %s (%s)", appData.Config.AppName, appData.Config.AppVersion)

	r := router.New(appData)

	router.StartServer(appData.Config.Port, r)
}
