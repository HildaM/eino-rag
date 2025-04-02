package conf

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 配置
type Config struct {
	// rag 配置
	Rag struct {
		Dimension int64 `mapstructure:"dimension"` // 服务端嵌入维度
	} `mapstructure:"rag"`

	// 模型配置
	DeekSeek struct {
		APIKey  string `mapstructure:"api_key"`  // 模型 API Key
		ModelID string `mapstructure:"model_id"` // 模型 ID
		BaseURL string `mapstructure:"base_url"` // 模型 API Base URL
	} `mapstructure:"DeekSeek"`

	// Redis 配置
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
	} `mapstructure:"Redis"`
}

var (
	// 全局配置实例
	AppConfig Config
)

// Init 初始化配置系统
func Init() error {
	v := viper.New()
	v.SetConfigType("yaml")   // 配置文件类型
	v.SetConfigName("config") // 配置文件名
	v.AddConfigPath(".")      // 配置文件路径

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		log.Printf("Init failed, read cfg err: %v\n", err)
		return err
	}
	log.Printf("Init success, cfg: %v\n", v.ConfigFileUsed())

	// 监听配置文件变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s\n", e.Name)
		if err := v.Unmarshal(&AppConfig); err != nil {
			log.Printf("Config file changed, but unmarshal to cfg failed, err: %v\n", err)
		}
	})

	// 解析配置文件
	if err := v.Unmarshal(&AppConfig); err != nil {
		log.Printf("Init failed, unmarshal to cfg failed, err: %v\n", err)
		return err
	}
	return nil
}

// GetCfg 获取全局配置实例
func GetCfg() *Config {
	return &AppConfig
}
