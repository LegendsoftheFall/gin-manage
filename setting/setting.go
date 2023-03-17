package setting

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

// Conf 全局变量 保存程序的所有配置选项
var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	Port         int    `mapstructure:"port"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int64  `mapstructure:"machine_id"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
	*QiniuConfig `mapstructure:"qiniu"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host        string `mapstructure:"host"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	DbName      string `mapstructure:"dbname"`
	Port        int    `mapstructure:"port"`
	MaxConn     int    `mapstructure:"max_conn"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Port     int    `mapstructure:"port"`
	PoolSize int    `mapstructure:"pool_size"`
}

type QiniuConfig struct {
	Path          string `mapstructure:"path"`
	Prefix        string `mapstructure:"prefix"`
	Bucket        string `mapstructure:"bucket"`
	Domain        string `mapstructure:"domain"`
	AccessKey     string `mapstructure:"accessKey"`
	SecretKey     string `mapstructure:"secretKey"`
	UseHttps      bool   `mapstructure:"useHttps"`
	UseCdnDomains bool   `mapstructure:"useCdnDomains"`
}

func Init() (err error) {
	viper.SetConfigFile("./conf/config.yaml")
	//viper.SetConfigName("config") // 配置文件名称(无扩展名)
	//viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	//viper.AddConfigPath(".")      // 查找配置文件所在的路径
	err = viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		return // 处理读取配置文件的错误
	}
	// 将读取到的数据反序列化到变量中
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}
	viper.WatchConfig() //热加载
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置文件已修改")
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})
	return
}
