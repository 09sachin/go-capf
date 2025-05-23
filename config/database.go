package config

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)



var (
	// host      = getEnv("DB_HOST")
    localhost = getEnv("DB_HOST_LOCAL")
    // port      = getEnv("DB_PORT")
    localport = getEnv("DB_PORT_LOCAL")
    // user      = getEnv("DB_USER")
    localuser = getEnv("DB_USER_LOCAL")
    // password  = getEnv("DB_PASS")
    localpass = getEnv("DB_PASS_LOCAL")
    // dbname    = getEnv("DB_NAME")
    localname = getEnv("DB_NAME_LOCAL")
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
		localhost, localport, localuser, localpass, localname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectDBLocal() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		localhost, localport, localuser, localpass, localname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ExecuteQueryLocal(query string, args ...interface{}) (*sql.Rows, error) {
	db, err := connectDBLocal()
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
    db, err := connectDBLocal()
    if err != nil {
        return err
    }
    defer db.Close()

    _, err = db.Exec(query, args...)
    if err != nil {
        return err
    }

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

    return nil
}
