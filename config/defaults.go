package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx"
)

// Service Ports
const (
	HTTPAddr    = ":8000"
	HTTPSAddr   = ":8001"
	GRPCAddr    = ":8002"
	MetricsAddr = ":8003"
	HealthAddr  = ":8004"
	AliveAddr   = ":8005"
)

func EnvPostgresConfig() pgx.ConnPoolConfig {
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	return pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     uint16(port),
			Database: os.Getenv("DB_DATABASE"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		// TLSConfig: os.Getenv("DB_DATABASE"),
	}
}

func GRPCServiceAddr(hostname string) string {
	return fmt.Sprint(hostname, GRPCAddr)
}
