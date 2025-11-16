package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
	DBMaxConns int
	DBTimeout  time.Duration
}

func Load() *Config {
	cfg := &Config{}

	defaultPort := getEnv("HTTP_PORT", "8080")
	defaultReadTimeout := getEnv("READ_TIMEOUT", "10s")
	defaultWriteTimeout := getEnv("WRITE_TIMEOUT", "10s")

	defaultDBHost := getEnv("DB_HOST", "localhost")
	defaultDBPort := getEnv("DB_PORT", "5432")
	defaultDBName := getEnv("DB_NAME", "postgres")
	defaultDBUser := getEnv("DB_USER", "postgres")
	defaultDBPassword := getEnv("DB_PASSWORD", "postgres")
	defaultDBSSLMode := getEnv("DB_SSLMODE", "disable")
	defaultDBMaxConns := getEnvInt("DB_MAX_CONNS", 50)
	defaultDBTimeout := getEnv("DB_TIMEOUT", "30s")

	flag.StringVar(&cfg.HTTPPort, "http-port", defaultPort, "HTTP server port")
	flag.DurationVar(&cfg.ReadTimeout, "read-timeout", parseDuration(defaultReadTimeout), "Read timeout")
	flag.DurationVar(&cfg.WriteTimeout, "write-timeout", parseDuration(defaultWriteTimeout), "Write timeout")

	flag.StringVar(&cfg.DBHost, "db-host", defaultDBHost, "Database host")
	flag.StringVar(&cfg.DBPort, "db-port", defaultDBPort, "Database port")
	flag.StringVar(&cfg.DBName, "db-name", defaultDBName, "Database name")
	flag.StringVar(&cfg.DBUser, "db-user", defaultDBUser, "Database user")
	flag.StringVar(&cfg.DBPassword, "db-password", defaultDBPassword, "Database password")
	flag.StringVar(&cfg.DBSSLMode, "db-ssl-mode", defaultDBSSLMode, "Database SSL mode")
	flag.IntVar(&cfg.DBMaxConns, "db-max-conns", defaultDBMaxConns, "Database max connections")
	flag.DurationVar(&cfg.DBTimeout, "db-timeout", parseDuration(defaultDBTimeout), "Database connection timeout")

	flag.Parse()

	return cfg
}

func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBUser,
		c.DBPassword,
		c.DBSSLMode,
	)
}

func (c *Config) Validate() error {
	if dbPort, err := strconv.Atoi(c.DBPort); err != nil {
		return fmt.Errorf("invalid DB port: %w", err)
	} else if dbPort < 1 || dbPort > 65535 {
		return fmt.Errorf("DB port %d is out of range 1-65535", dbPort)
	}

	if c.DBHost == "" {
		return fmt.Errorf("DB host is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB name is required")
	}
	if c.DBUser == "" {
		return fmt.Errorf("DB user is required")
	}

	if c.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be positive")
	}
	if c.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be positive")
	}
	if c.DBTimeout <= 0 {
		return fmt.Errorf("DB timeout must be positive")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func parseDuration(value string) time.Duration {
	dur, err := time.ParseDuration(value)
	if err != nil {
		if strings.Contains(strings.ToLower(value), "min") {
			return 10 * time.Minute
		}
		return 10 * time.Second
	}
	return dur
}
