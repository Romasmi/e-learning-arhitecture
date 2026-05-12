package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Db             Database
	Server         Server
	Kafka          Kafka
	AuthServiceURL string `mapstructure:"auth_service_url"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Server struct {
	GRPCPort uint `mapstructure:"grpc_port"`
	HTTPPort uint `mapstructure:"http_port"`
}

type Kafka struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
}

func bindEnvRecursive(viperInstance *viper.Viper, prefix string, val reflect.Value) error {
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			tag = FirstCharToLowerCase(field.Name)
		}

		var fieldPath string
		if prefix != "" {
			fieldPath = prefix + "." + tag
		} else {
			fieldPath = tag
		}

		if field.Type.Kind() == reflect.Struct {
			if err := bindEnvRecursive(viperInstance, fieldPath, val.Field(i)); err != nil {
				return err
			}
		} else {
			envVarName := strings.ToUpper(strings.ReplaceAll(fieldPath, ".", "_"))
			if err := viperInstance.BindEnv(fieldPath, envVarName); err != nil {
				return err
			}
		}
	}

	return nil
}

func bindAllEnvVars(viperInstance *viper.Viper) error {
	return bindEnvRecursive(viperInstance, "", reflect.ValueOf(&Config{}).Elem())
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v2 := viper.New()
	v2.SetConfigName("override")
	v2.SetConfigType("yaml")
	v2.AddConfigPath(configPath)
	if err := v2.ReadInConfig(); err == nil {
		err := v.MergeConfigMap(v2.AllSettings())
		if err != nil {
			return nil, err
		}
	}

	if err := bindAllEnvVars(v); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	fmt.Println("cfg", cfg)
	return &cfg, nil
}

func FirstCharToLowerCase(str string) string {
	firstChar := str[:1]
	return strings.ToLower(firstChar) + str[1:]
}
