package main

import (
	"github.com/reeveci/reeve-lib/schema"
)

func (p *LocalPlugin) Resolve(env []string) (map[string]schema.Env, error) {
	return p.Store.GetSomeEnv(env), nil
}
