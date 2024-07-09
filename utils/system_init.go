package utils

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

/*
viper第三方库用于读取不同数据源和格式的配置文件
*/
func InitConfig() {
	viper.SetConfigName("app")      //配置文件名
	viper.AddConfigPath("./config") //查找config文件
	err := viper.ReadInConfig()     //读取配置文件
	if err != nil {
		fmt.Println("--读取失败----", err)
	}
	fmt.Println("config mysql:", viper.Get("mysql"))
	fmt.Println("redis init:", viper.Get("redis"))
}
func InitMySQL() {
	//自定义日志模板，打印SQL语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")),
		&gorm.Config{Logger: newLogger})
	//DB.AutoMigrate(&models.UserBasic{}) //创建表?判断是否表结构存在
}

func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis."),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	//pong, err := Red.Ping().Result()
	//if err != nil {
	//	fmt.Println("redis init------", err)
	//} else {
	//	fmt.Println("redis inited-----", pong)
	//}
}

const (
	PublishKey = "websocket"
)

// 发布消息到Redis
func Publish(c context.Context, channel string, msg string) error {
	var err error
	err = Red.Publish(c, channel, msg).Err()
	fmt.Println("Publish------", msg)
	if err != nil {
		fmt.Println("发送失败：", err)
	}
	return err
}

// 订阅Redis消息
func Subscribe(c context.Context, channel string) (string, error) {
	sub := Red.Subscribe(c, channel)
	fmt.Println("Subscribe1------", c)
	msg, err := sub.ReceiveMessage(c)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Subscribe2------", msg.Payload)
	return msg.Payload, err
}
