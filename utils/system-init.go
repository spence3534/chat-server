package utils

import (
	"context"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

// 设置yaml配置

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app", viper.Get("app"))
	fmt.Println("config mysql", viper.Get("mysql"))
}

func InitMySQL() {
	// sql打印
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		})

	// 连接mysql数据基本配置
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN: viper.GetString("mysql.dsn"),
	}), &gorm.Config{
		SkipDefaultTransaction:                   false,
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   newLogger,
	})
	if err != nil {
		fmt.Println("连接数据库失败")
	}
}

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	ctx := context.Background()

	pong, err := Redis.Ping(ctx).Result()

	if err != nil {
		fmt.Println("init redis...", err)
	} else {
		fmt.Println("redis init success...", pong)
	}
}

func InitOss() (*oss.Client, error) {
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		fmt.Println("get oss key Error", err)
		return nil, err
	}

	id := provider.GetCredentials().GetAccessKeyID()
	secret := provider.GetCredentials().GetAccessKeySecret()

	client, err := oss.New("https://oss-cn-hangzhou.aliyuncs.com", id, secret, oss.SetCredentialsProvider(&provider))
	if err != nil {
		fmt.Println("oss new error", err)
	}

	return client, nil
}

var (
	PublishKey = "websocket"
)

// Publish 发布消息到redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("Publish...", err)
	err = Redis.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 订阅redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Redis.Subscribe(ctx, channel)
	fmt.Println("Subscribe1111...", ctx)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Subscribe...", msg.Payload)
	return msg.Payload, err
}
