package banco

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	DBName     string `json:"db_name"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
}

type APIConfig struct {
	APIPort   int    `json:"api_port"`
	APISecret string `json:"api_secret"`
}

var (
	Driver        = "postgres"
	APIConfigInfo APIConfig
)

func LoadAPIConfig(filePath string) (APIConfig, error) {
	var config APIConfig

	// Lê o conteúdo do arquivo JSON
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	// Decodifica o conteúdo do arquivo JSON para a estrutura APIConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}

func LoadDatabaseConfig(filePath string) (DatabaseConfig, error) {
	var dbConfig DatabaseConfig
	file, err := os.Open(filePath)
	if err != nil {
		return dbConfig, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&dbConfig)
	if err != nil {
		return dbConfig, fmt.Errorf("failed to decode config file: %w", err)
	}
	return dbConfig, nil
}

func Connection() (*sql.DB, error) {
	// Carrega as configurações do db do arquivo json
	dbConfig, err := LoadDatabaseConfig("config/config.api.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	// Carrega as configurações da API do arquivo JSON
	APIConfigInfo, err = LoadAPIConfig("config/config.api.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load API config: %w", err)
	}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBUser, dbConfig.DBName, dbConfig.DBPassword,
	)

	db, err := sql.Open(Driver, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	//log.Printf("Successfully connected to database at %s:%s", dbConfig.DBHost, dbConfig.DBPort)

	return db, nil

}
