package main

import (
	"fmt"

	"github.com/reeveci/reeve-lib/schema"
)

func (p *LocalPlugin) Message(source string, message schema.Message) error {
	switch {
	case source == schema.MESSAGE_SOURCE_SERVER:
		if message.Options["event"] == schema.EVENT_STARTUP_COMPLETE {
			p.Store.Notify()
		}
		return nil

	case source == "webui" && message.Options["webui"] == "present":
		p.Lock()
		p.WebUIPresent = true
		p.Unlock()

		SendEnvBundleMessage(p, EnvBundle{
			BundleID: "operations",
			Prompts: []Prompt{
				{
					ID:          "input-variable",
					Name:        "Local storage",
					NameOption:  "name",
					ValueOption: "value",
					Message: schema.Message{
						Target: PLUGIN_NAME,
						Options: map[string]string{
							"type":      "operation",
							"operation": "set-env",
							"secret":    "false",
						},
					},
				},
				{
					ID:          "input-secret",
					Name:        "Local storage",
					Secret:      true,
					NameOption:  "name",
					ValueOption: "value",
					Message: schema.Message{
						Target: PLUGIN_NAME,
						Options: map[string]string{
							"type":      "operation",
							"operation": "set-env",
							"secret":    "true",
						},
					},
				},
			},
		})

		return nil

	default:
	}

	switch message.Options["type"] {
	case "operation":
		operation := message.Options["operation"]
		if operation == "" {
			return fmt.Errorf("missing operation")
		}

		switch operation {
		case "set-env":
			var secret bool
			switch message.Options["secret"] {
			case "true":
				secret = true
			case "false", "":
				secret = false
			default:
				return fmt.Errorf("invalid boolean %s", message.Options["secret"])
			}
			name := message.Options["name"]
			if name == "" {
				return fmt.Errorf("missing name")
			}
			value := message.Options["value"]
			err := p.Store.SetEnv(name, value, secret)
			if err != nil {
				return fmt.Errorf("error setting env - %s", err)
			}
		}
	}

	return nil
}
