package database

import (
	"database/sql"
	"day06/ex01/internal/credentials"
	"fmt"
	_ "github.com/lib/pq"
)

func runDB() (*sql.DB, error) {
	db, err := sql.Open("postgres",
    	fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", credentials.AC.DBUser, credentials.AC.DBPassword, credentials.AC.DBName))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS articles (
            id SERIAL PRIMARY KEY,
            articlename VARCHAR(255) NOT NULL,
            linktoarticle VARCHAR(255) NOT NULL
        )
    `
    _, err = db.Exec(createTableQuery)
    if err != nil {
        return nil, err
    }

	return db, nil
}

func GetArticles() (map[string]string, error) {
	db, err := runDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make(map[string]string)
	for rows.Next() {
		var id int
		var name, link string
		err := rows.Scan(&id, &name, &link)
		if err != nil {
			return nil, err
		}
		articles[link] = name
	}
	return articles, nil
}