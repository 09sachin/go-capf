package config

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)



var (
	host     = getEnv("DB_HOST")
    port     = getEnv("DB_PORT")
    user     = getEnv("DB_USER")
    password = getEnv("DB_PASS")
    dbname   = getEnv("DB_NAME")
)

func getEnv(key string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        fmt.Printf("Warning: Environment variable %s is not set.\n", key)
        return ""
    }
    return value
}



func connectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database")
	return db, nil
}

// ExecuteQuery executes a SQL query and returns the result
func ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}


// InsertData inserts data into the database using a custom query
func InsertData(query string, args ...interface{}) error {
    db, err := connectDB()
    if err != nil {
        return err
    }
    defer db.Close()

    _, err = db.Exec(query, args...)
    if err != nil {
        return err
    }

    fmt.Println("Data inserted successfully")
    return nil
}

// UpdateData updates data in the database using a custom query
func UpdateData(query string, args ...interface{}) error {
    db, err := connectDB()
    if err != nil {
        return err
    }
    defer db.Close()

    _, err = db.Exec(query, args...)
    if err != nil {
        return err
    }

    fmt.Println("Data updated successfully")
    return nil
}
