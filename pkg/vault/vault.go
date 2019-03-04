package vault

import (
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/token"
	// "github.com/hashicorp/vault/helper/parseutil"
	"github.com/readytalk/stim/pkg/log"

	"errors"
	"time"
)

type Vault struct {
	client      *api.Client
	config      *Config
	tokenHelper token.InternalTokenHelper
	newLogin    bool
}

type Config struct {
	Noprompt bool
	Address  string
	Username string
	Timeout  time.Duration
	Log      log.Logger
}

func New(config *Config) (*Vault, error) {
	// Ensure that the Vault address is set
	if config.Address == "" {
		return nil, errors.New("Vault address not set")
	}

	v := &Vault{config: config}
	log.SetLogger(config.Log)

	// if v.config.Timeout == 0 {
	// 	v.config.Timeout = time.Second * 10 // No need to wait over a minite from default
	// }
	var err error
	// var clientTimeout time.Duration
	// clientTimeout, err = parseutil.ParseDurationSecond(10)
	// if err != nil {
	// 	return nil, err
	// }
	// v.config.Timeout = clientTimeout
	log.Debug("Vault Timeout: ", v.config.Timeout)

	// Configure new Vault Client
	apiConfig := api.DefaultConfig()
	apiConfig.Address = v.config.Address // Since we read the env we can override
	// apiConfig.HttpClient.Timeout = v.config.Timeout

	// Create our new API client
	v.client, err = api.NewClient(apiConfig)
	if err != nil {
		return nil, err
	}

	// Ensure Vault is up and Healthy
	_, err = v.isVaultHealthy()
	if err != nil {
		return nil, err
	}

	// Run Login logic
	err = v.Login()
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (v *Vault) GetUser() string {
	return v.config.Username
}
