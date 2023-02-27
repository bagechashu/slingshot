package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type cfg struct {
	Server
	DB
}

var Cfg = new(cfg)

func Show(cmd *cobra.Command, args []string) {
	fmt.Printf("Config: \n%+v\n", Cfg)
}

func InitConfig(cfgfile string) {
	viper.SetConfigFile(cfgfile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
	if err := viper.Unmarshal(Cfg); err != nil {
		panic(fmt.Errorf("unmarshal conf server failed, err:%s ", err))
	}
}
