package main

type Config struct {
	Env map[string]Env `json:"env"`
}

type Env struct {
	Value  string `json:"value"`
	Secret bool   `json:"secret"`
}
