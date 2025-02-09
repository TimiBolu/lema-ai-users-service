package config

import (
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
)

type ServerEnvironment string

const (
	ServerEnvironmentDevelopment ServerEnvironment = "development"
	ServerEnvironmentProduction  ServerEnvironment = "production"
)

type Config struct {
	PORT            string
	DB_NAME         string
	FRONTEND_APPS   string
	SERVER_BASE_URL string
	APP_ENV         ServerEnvironment
}

var EnvConfig = Config{}

func InitEnvSchema() {
	loadENV()

	envConfigReflection := reflect.ValueOf(&EnvConfig).Elem()
	envConfigType := envConfigReflection.Type()

	for i := 0; i < envConfigReflection.NumField(); i++ {
		field := envConfigType.Field(i)
		fieldName := field.Name
		envVariableValue := os.Getenv(fieldName)

		if envVariableValue == "" {
			log.Fatalf("You must set your %s environment variable.", fieldName)
		}

		switch field.Type.Kind() {
		case reflect.String:
			envConfigReflection.FieldByName(fieldName).SetString(envVariableValue)
		case reflect.Int:
			val, err := strconv.Atoi(envVariableValue)
			if err != nil {
				log.Fatalf("Invalid value for %s: %v", fieldName, err)
			}
			envConfigReflection.FieldByName(fieldName).SetInt(int64(val))
		default:
			log.Fatalf("Unsupported field type %s for field %s", field.Type.Kind(), fieldName)
		}
	}
}

func loadENV() {
	if goEnv := os.Getenv("GO_ENV"); goEnv == "" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	}
}
