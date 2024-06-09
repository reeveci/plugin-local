package main

import (
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/reeveci/reeve-lib/plugin"
	"github.com/reeveci/reeve-lib/schema"
)

const PLUGIN_NAME = "local"

func main() {
	log := hclog.New(&hclog.LoggerOptions{})

	plugin.Serve(&plugin.PluginConfig{
		Plugin: &LocalPlugin{
			Log: log,
		},

		Logger: log,
	})
}

type LocalPlugin struct {
	ConfigPath string
	SecretKey  string
	Priority   uint32

	Log hclog.Logger
	API plugin.ReeveAPI

	WebUIPresent bool
	sync.Mutex

	Store *LocalStore
}

func (p *LocalPlugin) Name() (string, error) {
	return PLUGIN_NAME, nil
}

func (p *LocalPlugin) Register(settings map[string]string, api plugin.ReeveAPI) (capabilities plugin.Capabilities, err error) {
	p.API = api

	var enabled bool
	if enabled, err = boolSetting(settings, "ENABLED"); !enabled || err != nil {
		return
	}
	if p.ConfigPath, err = requireSetting(settings, "CONFIG_PATH"); err != nil {
		return
	}
	if p.SecretKey, err = requireSetting(settings, "SECRET_KEY"); err != nil {
		return
	}
	var priority int
	if priority, err = intSetting(settings, "PRIORITY", 1); err != nil {
		return
	} else {
		p.Priority = uint32(priority)
	}

	if p.Store, err = NewLocalStore(p); err != nil {
		return
	}

	capabilities.Message = true
	capabilities.Resolve = true
	capabilities.CLIMethods = CLIMethods
	return
}

func (p *LocalPlugin) Unregister() error {
	p.API.Close()

	return nil
}

func (p *LocalPlugin) Discover(trigger schema.Trigger) ([]schema.Pipeline, error) {
	return nil, nil
}

func (p *LocalPlugin) Notify(status schema.PipelineStatus) error {
	return nil
}
