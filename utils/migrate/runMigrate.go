package migrate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func findMigrationsPath(name string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("не удалось получить текущий каталог: %w", err)
	}

	for {
		migrationsPath := filepath.Join(dir, name)
		if _, err := os.Stat(migrationsPath); err == nil {
			return migrationsPath, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	return "", fmt.Errorf("папка migrations не найдена")
}

func RunMigrate(name string) error {
	migrationsPath, err := findMigrationsPath(name)
	if err != nil {
		return err
	}
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	cmd := exec.Command("migrate", "-database", dsn, "-path", migrationsPath, "up")
	fmt.Println("Выполнение миграции:", cmd.String())

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
