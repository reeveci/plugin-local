package main

import (
	"fmt"
)

var CLIMethods = map[string]string{
	"var-set":    "<name> <value> - set environment variable",
	"var-get":    "<name> - get environment variable",
	"secret-set": "<name> <value> - set environment secret",
	"unset":      "<name> - unset environment variable or secret",
}

func (p *LocalPlugin) CLIMethod(method string, args []string) (string, error) {
	switch method {
	case "var-set":
		return p.CLISetVariable(args)

	case "var-get":
		return p.CLIGetVariable(args)

	case "secret-set":
		return p.CLISetSecret(args)

	case "unset":
		return p.CLIUnset(args)

	default:
		return "", fmt.Errorf("unknown method %s", method)
	}
}

func (p *LocalPlugin) CLISetVariable(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("var-set expects two arguments but got %v", len(args))
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
		return "", fmt.Errorf("secret-set expects two arguments but got %v", len(args))
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
