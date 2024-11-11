package services

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/justjcurtis/flxvwr/models"
	"github.com/justjcurtis/flxvwr/utils"
	"github.com/spf13/viper"
)

type ConfigService struct {
	config      models.Config
	subscribers []func(models.Config)
}

func NewConfigService() *ConfigService {
	configPath, err := utils.GetConfigPath()
	if err != nil {
		log.Fatal(err)
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	viper.SetDefault("app.name", "flxvwr")
	viper.SetDefault("app.version", "0.0.1")
	viper.SetDefault("shuffle", true)
	viper.SetDefault("delay", 10.0)

	viper.SafeWriteConfig()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No config file found, using defaults")
	}
	cs := &ConfigService{
		config: models.Config{
			Delay:   viper.GetDuration("delay") * time.Second,
			Shuffle: viper.GetBool("shuffle"),
		},
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		cs.Update()
	})
	viper.WatchConfig()
	return cs
}

func (cs *ConfigService) OnChange() {
	for _, fn := range cs.subscribers {
		fn(cs.config)
	}
}

func (cs *ConfigService) Update() {
	cs.config.Delay = viper.GetDuration("delay") * time.Second
	cs.config.Shuffle = viper.GetBool("shuffle")
	cs.OnChange()
}

func (cs *ConfigService) Subscribe(handler func(models.Config)) {
	cs.subscribers = append(cs.subscribers, handler)
}

func (cs *ConfigService) Save() {
	viper.Set("delay", cs.config.Delay.Seconds())
	viper.Set("shuffle", cs.config.Shuffle)
	viper.WriteConfig()
	cs.OnChange()
}

func (cs *ConfigService) SetDelay(delay time.Duration) {
	cs.config.Delay = delay
	cs.Save()
}

func (cs *ConfigService) GetDelay() time.Duration {
	return cs.config.Delay
}

func (cs *ConfigService) SetShuffle(shuffle bool) {
	cs.config.Shuffle = shuffle
	cs.Save()
}

func (cs *ConfigService) GetShuffle() bool {
	return cs.config.Shuffle
}
