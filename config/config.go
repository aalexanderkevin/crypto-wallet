package config

import (
	"sync"

	"github.com/jinzhu/configor"
)

type Config struct {
	Service   Service
	LogLevel  string `default:"INFO" env:"LOG_LEVEL"`
	LogFormat string `default:"json" env:"LOG_FORMAT"`
	Version   string
	Redis     Redis
	Ethereum  Ethereum
	Tron      Tron
	Bitcoin   Bitcoin
	Postgres  Postgres
	JwtSecret string `required:"true" env:"JWT_SECRET"`
}

type Postgres struct {
	Client    string `default:"postgres" env:"POSTGRES_CLIENT"`
	Host      string `default:"127.0.0.1" env:"POSTGRES_HOST"`
	Username  string `default:"root" env:"POSTGRES_USER"`
	Password  string `required:"true" env:"POSTGRES_PASSWORD"`
	Port      uint   `default:"5432" env:"POSTGRES_PORT"`
	Database  string `default:"gits" env:"POSTGRES_DATABASE"`
	Migration struct {
		Path string `default:"database/migration" env:"POSTGRES_MIGRATION_PATH"`
	}
	MaxIdleConnections int  `default:"25" env:"POSTGRES_MAX_IDLE_CONN"`
	MaxOpenConnections int  `default:"0" env:"POSTGRES_MAX_OPEN_CONN"`
	MaxConnLifeTime    int  `default:"90" env:"POSTGRES_MAX_CONN_LIFETIME"`
	Debug              bool `default:"false" env:"POSTGRES_DEBUG"`
}

type Service struct {
	Name   string `default:"crypto-wallet" env:"SERVICE_NAME"`
	Scheme string `default:"http" env:"SERVICE_SCHEME"`
	Host   string `default:"0.0.0.0" env:"SERVICE_HOST"`
	Port   string `default:"9004" env:"SERVICE_PORT"`
	Secret string `default:"secret" env:"SERVICE_SECRET"`
	Path   struct {
		V1  string `default:"/v1" env:"SERVICE_PATH_API"`
		Btc string `default:"/btc" env:"SERVICE_PATH_BTC"`
	}

	SeedPhraseEncryptionKey string `env:"SEED_PHRASE_ENCRYPTION_KEY"`
}

type Ethereum struct {
	NetUrl      string `default:"https://cloudflare-eth.com" env:"ETH_NET_URL"`
	Passphrase  string `default:"passphrase" env:"ETH_PASSPHRASE"`
	KeyStoreDir string `default:"./keystore" env:"ETH_KEY_STORE_DIR"`
}

type Tron struct {
	NetgRPCUrl  string `default:"grpc.shasta.trongrid.io:50051" env:"TRON_NET_GRPC_URL"`
	NetUrl      string `default:"https://api.shasta.trongrid.io" env:"TRON_NET_URL"`
	ApiKey      string `default:"09c32c49-d972-494c-96fd-eb5b1b0a4414" env:"TRON_API_KEY"`
	Passphrase  string `default:"passphrase" env:"TRON_PASSPHRASE"`
	KeyStoreDir string `default:"./keystore" env:"TRON_KEY_STORE_DIR"`
}

type Bitcoin struct {
	Chain               string `default:"test3" env:"BTC_CHAIN"`
	Token               string `default:"a843ce1e9a1c48ac9c621e12b9e8762a" env:"BTC_TOKEN"`
	WebhookURL          string `env:"BTC_WEBHOOK_URL"`
	MinimalConfirmation int    `default:"6" env:"BTC_MINIMAL_CONFIRMATION"`
}

type Redis struct {
	Host string `default:"localhost" env:"REDIS_HOST" json:"-"`
	Port uint   `default:"6379" env:"REDIS_PORT"`
}

var config *Config
var configLock = &sync.Mutex{}

// Instance
func Instance() Config {
	if config == nil {
		err := Load()
		if err != nil {
			panic(err)
		}
	}
	return *config
}

func Load() error {
	tmpConfig := Config{}
	err := configor.Load(&tmpConfig)
	if err != nil {
		return err
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &tmpConfig

	return nil
}
