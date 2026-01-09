package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Lang string // "en", es, tr
	API_SECRET_KEY string
	Port string
}

func LoadConfig() Config {

	if err := godotenv.Load(); err != nil {
		log.Println("Uyarı: .env dosyası bulunamadı, ortam değişkenleri okunacak")
	}

	lang := strings.TrimSpace(strings.ToLower(os.Getenv("APP_LANG"))) // En en

	if lang == "" {
		lang = "en"
	}

	switch lang {
	case "tr", "en", "es":
		// bunlardan birisi geliyorsa lang olduğu gibi kalsın
	default:
		lang = "en"
	}

	apiSecretKey := os.Getenv("API_SECRET_KEY")
	if apiSecretKey == "" {
		// Güvenlik önlemi: Eğer key yoksa uygulama çalışmamalı veya uyarılmalı
		log.Fatal("HATA: API_SECRET_KEY ortam değişkeni ayarlanmamış!")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Lang: lang,
		API_SECRET_KEY: apiSecretKey,
		Port: port,
	}
}