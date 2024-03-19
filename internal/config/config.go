package config

import (
	"flag"
	"log"

	"github.com/notnull-co/cfg"
)

var (
	conf Config
)

type Config struct {
	Token struct {
		Key string `cfg:"key"`
	} `cfg:"token"`
	Server struct {
		Port        string `cfg:"port"`
		ProductPort string `cfg:"product_integration"`
	} `cfg:"server"`
	DB struct {
		ConnectionString string `cfg:"connectionString"`
	} `cfg:"db"`
	SQS struct {
		PaymentPendingQueue   string `cfg:"payment_pending_queue"`
		PaymentPayedQueue     string `cfg:"payment_payed_queue"`
		PaymentCancelledQueue string `cfg:"payment_cancelled_queue"`
		OrderQueue            string `cfg:"order_queue"`
		Region                string `cfg:"region"`
		Endpoint              string `cfg:"endpoint"`
	} `cfg:"sqs"`
}

func ParseFromFlags() {
	var configDir string

	flag.StringVar(&configDir, "config-dir", "../../internal/config/", "Configuration file directory")
	flag.Parse()

	parse(configDir)
}

func parse(dirs ...string) {
	if err := cfg.Load(&conf,
		cfg.Dirs(dirs...),
		cfg.UseEnv("app"),
	); err != nil {
		log.Panic(err)
	}
}

func Get() Config {
	return conf
}
