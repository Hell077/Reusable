package migrate

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunMigrate() error {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	cmd := exec.Command("migrate", "-database", dsn, "-path", "./migrations", "up")

	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "no change") {
			fmt.Println("Миграция уже была применена, изменений нет.")
			return nil
		}
		return fmt.Errorf("ошибка при запуске миграции: %s", err)
	}

	fmt.Println(string(out))
	return nil
}
