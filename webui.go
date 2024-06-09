package main

import (
	"encoding/json"
	"fmt"

	"github.com/reeveci/reeve-lib/schema"
)

type EnvBundle struct {
	BundleID string                `json:"bundleID"`
	Env      map[string]schema.Env `json:"env"`
	Prompts  []Prompt              `json:"prompts"`
}

type Prompt struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	NameOption  string         `json:"nameOption"`
	ValueOption string         `json:"valueOption"`
	Secret      bool           `json:"secret"`
	Message     schema.Message `json:"message"`
}

func SendEnvBundleMessage(plugin *LocalPlugin, bundle EnvBundle) bool {
	data, err := json.Marshal(bundle)
	if err != nil {
		plugin.Log.Error(fmt.Sprintf("error building env bundle for WebUI - %s", err))
		return false
	}

	err = plugin.API.NotifyMessages([]schema.Message{{
		Target:  "webui",
		Options: map[string]string{"webui": "env"},
		Data:    data,
	}})
	if err != nil {
		plugin.Log.Error(fmt.Sprintf("error sending env to WebUI - %s", err))
		return false
	}

	return true
}
