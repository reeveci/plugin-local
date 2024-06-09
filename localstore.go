package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/reeveci/plugin-local/encryption"
	"github.com/reeveci/reeve-lib/schema"
)

func NewLocalStore(plugin *LocalPlugin) (*LocalStore, error) {
	s := &LocalStore{
		plugin: plugin,
	}

	err := s.read()
	if err != nil {
		return nil, err
	}

	return s, nil
}

type LocalStore struct {
	plugin *LocalPlugin

	lock sync.Mutex
	data Config
}

func (s *LocalStore) Notify() {
	s.plugin.Lock()
	hasUI := s.plugin.WebUIPresent
	s.plugin.Unlock()
	if hasUI {
		SendEnvBundleMessage(s.plugin, EnvBundle{
			BundleID: "env",
			Env:      s.GetAllEnv(),
		})
	}
}

func (s *LocalStore) SetEnv(name string, value string, secret bool) error {
	if name == "" {
		return fmt.Errorf("missing name")
	}

	s.lock.Lock()
	if value != "" {
		s.data.Env[name] = Env{Value: value, Secret: secret}
	} else {
		delete(s.data.Env, name)
	}
	s.lock.Unlock()

	err := s.write()
	if err != nil {
		return err
	}

	s.Notify()

	return nil
}

func (s *LocalStore) GetEnv(name string) (value string, secret bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	env := s.data.Env[name]
	return env.Value, env.Secret
}

func (s *LocalStore) GetSomeEnv(names []string) map[string]schema.Env {
	result := make(map[string]schema.Env, len(names))

	s.lock.Lock()
	defer s.lock.Unlock()

	for _, name := range names {
		env := s.data.Env[name]
		if env.Value != "" {
			result[name] = schema.Env{Value: env.Value, Priority: s.plugin.Priority, Secret: env.Secret}
		}
	}
	return result
}

func (s *LocalStore) GetAllEnv() map[string]schema.Env {
	s.lock.Lock()
	defer s.lock.Unlock()

	result := make(map[string]schema.Env, len(s.data.Env))
	for name, env := range s.data.Env {
		if env.Value != "" {
			result[name] = schema.Env{Value: env.Value, Priority: s.plugin.Priority, Secret: env.Secret}
		}
	}
	return result
}

func (s *LocalStore) write() error {
	s.lock.Lock()
	encryptedData, err := s.encrypt(s.data)
	s.lock.Unlock()
	if err != nil {
		return err
	}

	configDir := filepath.Join(s.plugin.ConfigPath, PLUGIN_NAME)
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("cannot create config directory %s - %s", configDir, err)
	}

	content, err := json.Marshal(encryptedData)
	if err != nil {
		return fmt.Errorf("cannot stringify config - %s", err)
	}

	configFile := filepath.Join(configDir, "config.json")
	err = os.WriteFile(configFile, content, 0600)
	if err != nil {
		return fmt.Errorf("cannot write config file %s - %s", configFile, err)
	}

	return nil
}

func (s *LocalStore) read() error {
	configDir := filepath.Join(s.plugin.ConfigPath, PLUGIN_NAME)
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("cannot create config directory %s - %s", configDir, err)
	}

	configFile := filepath.Join(configDir, "config.json")
	content, err := os.ReadFile(configFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cannot read config file %s - %s", configFile, err)
	}

	var data Config
	if len(content) > 0 {
		err := json.Unmarshal(content, &data)
		if err != nil {
			return fmt.Errorf("cannot parse config file - %s", err)
		}
	}

	decryptedData, err := s.decrypt(data)
	if err != nil {
		return err
	}

	s.lock.Lock()
	s.data = decryptedData
	s.lock.Unlock()

	return nil
}

func (s *LocalStore) encrypt(source Config) (result Config, err error) {
	result.Env = make(map[string]Env, len(source.Env))
	for key, env := range source.Env {
		if env.Secret {
			encryptedValue, err := encryption.EncryptSecret(s.plugin.SecretKey, env.Value)
			if err != nil {
				return result, fmt.Errorf("error encrypting secret %s - %s", key, err)
			}
			env.Value = encryptedValue
		}
		result.Env[key] = env
	}

	return
}

func (s *LocalStore) decrypt(source Config) (result Config, err error) {
	result.Env = make(map[string]Env, len(source.Env))
	for key, env := range source.Env {
		if env.Secret {
			decryptedValue, err := encryption.DecryptSecret(s.plugin.SecretKey, env.Value)
			if err != nil {
				return result, fmt.Errorf("error decrypting secret %s - %s", key, err)
			}
			env.Value = decryptedValue
		}
		result.Env[key] = env
	}

	return
}
