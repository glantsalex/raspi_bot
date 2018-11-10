package config

import (
	"github.com/cihub/seelog"
	"os"
)

const (

	CFG_WS_API_ROOT	 	= "WS_API_ROOT"
	CFG_HOST	 		= "HOST"
	CFG_PORT			= "PORT"
)

type Config map[string]string

func init(){
	for key, _ := range cfg{
		if tmp := os.Getenv( key ); tmp != ""{
			cfg[key] = tmp
			seelog.Debugf("key from environment: %s=%s", key,tmp )
		}
	}
}

var cfg  = Config {

	CFG_WS_API_ROOT  	: "/v1/ws/",
	CFG_HOST	:"",
	CFG_PORT	:"8090",
}

func GetConfig() *Config {
	return &cfg
}

func (c *Config)WsApiRoot()   			string {return cfg[CFG_WS_API_ROOT]}

func (c *Config)Host()		string	{return cfg[CFG_HOST]}
func (c *Config)Port()		string	{return cfg[CFG_PORT]}

