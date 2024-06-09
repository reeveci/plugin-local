package main

import (
	"fmt"
)

var CLIMethods = map[string]string{
	"set-variable": "<name> <value> - set environment variable",
	"get-variable": "<name> - get environment variable",
	"set-secret":   "<name> <value> - set secret environment variable",
	"unset":        "<name> - unset environment variable",
}

func (p *LocalPlugin) CLIMethod(method string, args []string) (string, error) {
	switch method {
	case "set-variable":
		return p.CLISetVariable(args)

	case "get-variable":
		return p.CLIGetVariable(args)

	case "set-secret":
		return p.CLISetSecret(args)

	case "unset":
		return p.CLIUnset(args)

	default:
		return "", fmt.Errorf("unknown method %s", method)
	}
}

func (p *LocalPlugin) CLISetVariable(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("set-variable expects two arguments but got %v", len(args))
	}

	name := args[0]
	value := args[1]

	if name == "" {
		return "", fmt.Errorf("no name was specified")
	}

	err := p.Store.SetEnv(name, value, false)
	if err != nil {
		return "", fmt.Errorf("setting variable failed - %s", err)
	}
	return "ok", nil
}

func (p *LocalPlugin) CLIGetVariable(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("encrypt expects one argument but got %v", len(args))
	}

	name := args[0]

	if name == "" {
		return "", fmt.Errorf("no name was specified")
	}

	value, secret := p.Store.GetEnv(name)
	if value == "" {
		return "", fmt.Errorf("the variable is not set")
	}
	if secret {
		return "", fmt.Errorf("the variable is a secret")
	}
	return value, nil
}

func (p *LocalPlugin) CLISetSecret(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("set-secret expects two arguments but got %v", len(args))
	}

	name := args[0]
	value := args[1]

	if name == "" {
		return "", fmt.Errorf("no name was specified")
	}

	err := p.Store.SetEnv(name, value, true)
	if err != nil {
		return "", fmt.Errorf("setting secret failed - %s", err)
	}
	return "ok", nil
}

func (p *LocalPlugin) CLIUnset(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("encrypt expects one argument but got %v", len(args))
	}

	name := args[0]

	if name == "" {
		return "", fmt.Errorf("no name was specified")
	}

	err := p.Store.SetEnv(name, "", false)
	if err != nil {
		return "", fmt.Errorf("unsetting key failed - %s", err)
	}
	return "ok", nil
}
