package plugins

import (
	"path/filepath"
	"plugin"
)

type GeneralPlugin struct {
	Location string
	Filename string
}

func (p *GeneralPlugin) Load(lookup string) (plugin.Symbol, error) {
	mod := filepath.Join(p.Location, p.Filename+PLUGIN_EXTENSION)
	plug, err := plugin.Open(mod)

	if err != nil {
		return nil, err
	}

	symReqEntryPlugin, err := plug.Lookup(lookup)
	if err != nil {
		return nil, err
	}

	return symReqEntryPlugin, nil
}
