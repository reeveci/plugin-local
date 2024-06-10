package main

import (
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"
)

var CLIMethods = map[string]string{
	"set":        "<name> <value> - set environment variable",
	"get":        "<name> - get environment variable",
	"set-secret": "<name> <value> - set environment secret",
	"unset":      "<name> - unset environment variable or secret",
	"list":       "- list environment variables and secrets",
}

func (p *LocalPlugin) CLIMethod(method string, args []string) (string, error) {
	switch method {
	case "set":
		return p.CLISet(args)

	case "get":
		return p.CLIGet(args)

	case "set-secret":
		return p.CLISetSecret(args)

	case "unset":
		return p.CLIUnset(args)

	case "list":
		return p.CLIList(args)

	default:
		return "", fmt.Errorf("unknown method %s", method)
	}
}

func (p *LocalPlugin) CLISet(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("two arguments are expected but there are %v", len(args))
	}

	name := args[0]
	value := args[1]

	if name == "" {
		return "", fmt.Errorf("no name was specified")
	}

	err := p.Store.SetEnv(name, value, false)
	if err != nil {
		return "", err
	}
	return "ok", nil
}

func (p *LocalPlugin) CLIGet(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("one argument is expected but there are %v", len(args))
	}

	name := args[0]

	if name == "" {
		return "", fmt.Errorf("no name was specified")
	}

	value, secret := p.Store.GetEnv(name)
	if value == "" {
		return "", fmt.Errorf("the variable does not exist")
	}
	if secret {
		return "", fmt.Errorf("the value is secret")
	}
	return value, nil
}

func (p *LocalPlugin) CLISetSecret(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("two arguments are expected but there are %v", len(args))
	}

	name := args[0]
	value := args[1]

	if name == "" {
		return "", fmt.Errorf("no name was specified")
	}

	err := p.Store.SetEnv(name, value, true)
	if err != nil {
		return "", err
	}
	return "ok", nil
}

func (p *LocalPlugin) CLIUnset(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("one argument is expected but there are %v", len(args))
	}

	name := args[0]

	if name == "" {
		return "", fmt.Errorf("no name was specified")
	}

	value, _ := p.Store.GetEnv(name)
	if value == "" {
		return "", fmt.Errorf("the variable does not exist")
	}

	err := p.Store.SetEnv(name, "", false)
	if err != nil {
		return "", err
	}
	return "ok", nil
}

func (p *LocalPlugin) CLIList(args []string) (string, error) {
	if len(args) != 0 {
		return "", fmt.Errorf("no arguments are expected but there are %v", len(args))
	}

	env := p.Store.GetAllEnv()
	list := struct {
		Vars    []string `yaml:"vars"`
		Secrets []string `yaml:"secrets"`
	}{
		Vars:    make([]string, 0, len(env)),
		Secrets: make([]string, 0, len(env)),
	}
	for name, env := range env {
		if env.Secret {
			list.Secrets = append(list.Secrets, name)
		} else {
			list.Vars = append(list.Vars, name)
		}
	}
	sort.Strings(list.Vars)
	sort.Strings(list.Secrets)

	result, err := yaml.Marshal(list)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
