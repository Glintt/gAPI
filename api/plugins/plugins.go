package plugins

import (
	"os"
	"path/filepath"
	"strings"
)

func ListAll() (map[string][]string, error) {
	var filesMap map[string][]string
	filesMap = make(map[string][]string)

	configuration := Configurations()

	for _, pluginType := range PLUGINS_TYPES {
		var filesList []string
		typeLocation := configuration.Location + string(filepath.Separator) + pluginType

		filepath.Walk(typeLocation, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(info.Name()) != PLUGIN_EXTENSION {
				return nil
			}

			filesList = append(filesList, strings.Replace(info.Name(), ".so", "", -1))
			return nil
		})

		filesMap[pluginType] = filesList
	}

	return filesMap, nil
}
