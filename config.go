package main

import (
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-github-server/mq"
)

type Config struct {
	MQ     mq.Config `json:"mq"`
	Org    string    `json:"org"`
	Topics Topics    `json:"topics"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := new(Config)
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return nil, err
	}

	cfg.SetDefault()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) SetDefault() {
	if c.Org == "" {
		c.Org = "src-openeuler"
	}
}

func (c *Config) Validate() error {
	if _, err := utils.BuildRequestBody(c, ""); err != nil {
		return err
	}

	return nil
}

type Topics struct {
	ApprovedPkg string `json:"approved_pkg" required:"true"`
	MergedPR    string `json:"merged_pr"    required:"true"`
	CreatedRepo string `json:"created_repo" required:"true"`
}
