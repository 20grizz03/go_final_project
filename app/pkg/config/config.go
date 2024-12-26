package config

import (
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

var DB *sql.DB

// загружаем переменные окружения
func LoadEnviroment() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}
}

// создаем БД
func MakeDB() {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = filepath.Join(filepath.Dir(appPath), "database", "scheduler.db")
	}

	log.Printf("Путь к базе данных: %s", dbFile)

	install := false
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Если база данных новая, создаём таблицы и индексы
	if install {
		err = createSchema(DB)
		if err != nil {
			log.Fatalf("Ошибка создания схемы базы данных: %v", err)
		}
		log.Println("База данных успешно создана и проинициализирована.")
	}

}

// создаем таблицы и индексы в базе данных
func createSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT CHECK (LENGTH(repeat) <= 128)
	);

	CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
	`
	_, err := db.Exec(query)
	return err
}

// закрываем соединение с базой данных
func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Fatalf("Ошибка закрытия базы данных: %v", err)
		}
		log.Println("Соединение с базой данных закрыто.")
	}
}
