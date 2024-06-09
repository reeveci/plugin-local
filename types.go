package main

type Config struct {
	Env map[string]Env
}

type Env struct {
	Value, Encrypted string
	Secret           bool
}

type ExternalConfig struct {
	Env map[string]ExternalEnv `json:"env"`
}

type ExternalEnv struct {
	Value  string `json:"value"`
	Secret bool   `json:"secret"`
}
