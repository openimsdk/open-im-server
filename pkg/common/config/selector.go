package config

import "encoding/json"

type EngineSelector struct {
	Engine string `yaml:"engine" mapstructure:"engine" json:"engine"`
}

func (e EngineSelector) String() string {
	return e.Engine
}

func (e *EngineSelector) UnmarshalYAML(unmarshal func(any) error) error {
	var engine string
	if err := unmarshal(&engine); err == nil {
		e.Engine = engine
		return nil
	}
	var cfg struct {
		Engine string `yaml:"engine"`
	}
	if err := unmarshal(&cfg); err != nil {
		return err
	}
	e.Engine = cfg.Engine
	return nil
}

func (e *EngineSelector) UnmarshalJSON(data []byte) error {
	var engine string
	if err := json.Unmarshal(data, &engine); err == nil {
		e.Engine = engine
		return nil
	}
	var cfg struct {
		Engine string `json:"engine"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	e.Engine = cfg.Engine
	return nil
}
