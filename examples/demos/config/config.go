// ══════════════════════════════════════════════════════════════════════════════
// FILE: config.go - MT5 CONNECTION CONFIGURATION LOADER
// ══════════════════════════════════════════════════════════════════════════════
//
// 🎯 WHAT IS THIS?
//   Universal configuration loader for MT5 connection settings.
//   Used by ALL demo examples to connect to MT5 servers.
//
//  LOADING PRIORITY:
//   1️. config.json file (if exists)
//   2️. Environment variables (fallback)
//   3️. Error if neither found
//
// 📁 METHOD 1: Using config.json (Recommended)
// ─────────────────────────────────────────────────────────────────────────────
//   Create file: examples/demos/config/config.json
//
//   {
//     "user": 591129415,
//     "password": "YourPassword",
//     "host": "mt5.mrpc.pro",
//     "port": 443,
//     "grpc_server": "mt5.mrpc.pro:443",
//     "mt_cluster": "FxPro-MT5 Demo",
//     "test_symbol": "EURUSD",
//     "test_volume": 0.01
//   }
//
// 🌍 METHOD 2: Using Environment Variables
// ─────────────────────────────────────────────────────────────────────────────
//   REQUIRED:
//     MT5_USER        - MT5 account number (uint64)
//     MT5_PASSWORD    - MT5 account password
//     MT5_HOST        - MT5 server host (e.g., "mt5.mrpc.pro")
//
//   OPTIONAL:
//     MT5_PORT        - Server port (default: 443)
//     MT5_GRPC_SERVER - Full gRPC server address (default: HOST:PORT)
//     MT5_CLUSTER     - MT5 cluster name (e.g., "FxPro-MT5 Demo")
//     MT5_TEST_SYMBOL - Symbol for testing (default: "EURUSD")
//     MT5_TEST_VOLUME - Volume for testing (default: 0.01)
//
//   Example (Linux/Mac):
//     export MT5_USER=591129415
//     export MT5_PASSWORD="YourPassword"
//     export MT5_HOST="mt5.mrpc.pro"
//     export MT5_PORT=443
//     export MT5_GRPC_SERVER="mt5.mrpc.pro:443"
//     export MT5_CLUSTER="FxPro-MT5 Demo"
//
//   Example (Windows PowerShell):
//     $env:MT5_USER="591129415"
//     $env:MT5_PASSWORD="YourPassword"
//     $env:MT5_HOST="mt5.mrpc.pro"
//
// 💡 SMART DEFAULTS:
//   • If MT5_GRPC_SERVER not set → auto-constructs from MT5_HOST:MT5_PORT
//   • If MT5_PORT not set → defaults to 443
//   • If MT5_TEST_SYMBOL not set → defaults to "EURUSD"
//   • If MT5_TEST_VOLUME not set → defaults to 0.01
//
// 📖 USAGE IN CODE:
//   cfg, err := config.LoadConfig()
//   if err != nil {
//       log.Fatal(err)
//   }
//   // Now use cfg.User, cfg.Password, cfg.GrpcServer, etc.
//
// ══════════════════════════════════════════════════════════════════════════════

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// MT5Config contains connection settings for MT5
type MT5Config struct {
	User       uint64  `json:"user"`
	Password   string  `json:"password"`
	Host       string  `json:"host"`
	Port       int32   `json:"port"`
	GrpcServer string  `json:"grpc_server"`
	MtCluster  string  `json:"mt_cluster"`
	TestSymbol string  `json:"test_symbol"`
	TestVolume float64 `json:"test_volume"`
}

// LoadConfig loads configuration from file or environment variables
// Priority: 1. config.json file, 2. environment variables
func LoadConfig() (*MT5Config, error) {
	// Try to load from config.json first
	config, err := loadFromFile("config/config.json")
	if err == nil {
		fmt.Println("✓ Loaded configuration from config.json")
		return config, nil
	}

	// If config file not found, try environment variables
	config, err = loadFromEnv()
	if err == nil {
		fmt.Println("✓ Loaded configuration from environment variables")
		return config, nil
	}

	return nil, fmt.Errorf("no configuration found: %w", err)
}

// loadFromFile loads configuration from JSON file
func loadFromFile(filename string) (*MT5Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config MT5Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv() (*MT5Config, error) {
	user := os.Getenv("MT5_USER")
	password := os.Getenv("MT5_PASSWORD")
	host := os.Getenv("MT5_HOST")
	port := os.Getenv("MT5_PORT")
	grpcServer := os.Getenv("MT5_GRPC_SERVER")

	if user == "" || password == "" || host == "" {
		return nil, fmt.Errorf("required environment variables not set")
	}

	userInt, err := strconv.ParseUint(user, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid MT5_USER: %w", err)
	}

	portInt := int32(443) // default
	if port != "" {
		p, err := strconv.ParseInt(port, 10, 32)
		if err == nil {
			portInt = int32(p)
		}
	}

	// If GRPC_SERVER not set, construct from host
	if grpcServer == "" {
		grpcServer = fmt.Sprintf("%s:%d", host, portInt)
	}

	return &MT5Config{
		User:       userInt,
		Password:   password,
		Host:       host,
		Port:       portInt,
		GrpcServer: grpcServer,
		MtCluster:  os.Getenv("MT5_CLUSTER"),
		TestSymbol: getEnvOrDefault("MT5_TEST_SYMBOL", "EURUSD"),
		TestVolume: getEnvFloatOrDefault("MT5_TEST_VOLUME", 0.01),
	}, nil
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvFloatOrDefault(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}
	return defaultValue
}
