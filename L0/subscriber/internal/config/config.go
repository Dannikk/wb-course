package config

import (
	"os"

	"github.com/spf13/viper"
)

type PostgresCfg struct {
	User     string
	Password string
	DB       string
	Hostname string
}

type StanCfg struct {
	ClusterID   string
	ClientID    string
	Hostname    string
	Port        string
	Subject     string
	Qgroup      string
	Durable     string
	StartSeq    uint64
	DeliverLast bool
	DeliverAll  bool
	NewOnly     bool
	StartDelta  string
}

type Config struct {
	Pgsql PostgresCfg
	Stan  StanCfg
}

func NewConfig(path string) (Config, error) {
	c := Config{}

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return c, err
	}

	pghost, ok := os.LookupEnv("PGHOST")
	if !ok {
		pghost = "localhost"
	}
	c.Pgsql = PostgresCfg{
		User:     viper.GetString("POSTGRES_USER"),
		Password: viper.GetString("POSTGRES_PW"),
		DB:       viper.GetString("POSTGRES_DB"),
		Hostname: pghost,
	}

	viper.SetDefault("NATS_DELIVER_ALL", true)
	viper.SetDefault("NATS_PORT", "4222")
	c.Stan = StanCfg{
		ClusterID:   viper.GetString("NATS_CLUSTER_ID"),
		ClientID:    viper.GetString("NATS_CLIENT_ID"),
		Hostname:    viper.GetString("NATS_HOSTNAME"),
		Port:        viper.GetString("NATS_PORT"),
		Subject:     viper.GetString("NATS_SUBJECT"),
		Qgroup:      viper.GetString("NATS_QGROUP"),
		Durable:     viper.GetString("NATS_DURABLE"),
		StartSeq:    viper.GetUint64("NATS_START_SEQ"),
		DeliverLast: viper.GetBool("NATS_DELIVER_LAST"),
		DeliverAll:  viper.GetBool("NATS_DELIVER_ALL"),
		NewOnly:     viper.GetBool("NATS_NEW_ONLY"),
		StartDelta:  viper.GetString("NATS_START_DELTA"),
	}

	return c, nil
}
