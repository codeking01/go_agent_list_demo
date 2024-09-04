package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Postgres driver
	"log"
	"math/rand"
	"time"
)

/*
// pg中建表

CREATE TABLE terminals (
                           id SERIAL PRIMARY KEY,
                           mac VARCHAR(17),
                           cpu_usage FLOAT,
                           mem_usage FLOAT,
                           last_modify_time TIMESTAMP WITHOUT TIME ZONE
);
*/

var db *sql.DB

func InitPostGreSQLDB(cfg *Config) {
	var err error
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	// Open the database connection
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Ping the database to ensure the connection is established
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Optionally, configure the connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * 60)

	fmt.Println("Successfully connected to PostgreSQL")

	// Insert initial data into the database
	InsertInitialData(db)
}

func InsertInitialData(db *sql.DB) {
	fmt.Println("Inserting initial data...")

	for i := 1; i <= 50000; i++ {
		terminal := Terminal{
			ID:             i,
			MAC:            fmt.Sprintf("00:1A:C2:%02X:%02X:%02X", rand.Intn(256), rand.Intn(256), rand.Intn(256)),
			CPUUsage:       rand.Float64() * 100,
			MemUsage:       rand.Float64() * 100,
			LastModifyTime: time.Now().UTC().Truncate(time.Minute), // Truncate to the nearest minute
		}
		// Format time as string without timezone information
		lastModifyTimeFormatted := time.Now().UTC().Truncate(time.Minute).Format("2006-01-02 15:04:05")
		query := `INSERT INTO terminals (id, mac, cpu_usage, mem_usage, last_modify_time) VALUES ($1, $2, $3, $4, $5)`
		_, err := db.Exec(query, terminal.ID, terminal.MAC, terminal.CPUUsage, terminal.MemUsage, lastModifyTimeFormatted)
		if err != nil {
			log.Fatalf("Failed to insert data: %v", err)
		}
	}

	fmt.Println("Initial data insertion completed.")
}

func AgentSaveTerminal(terminal *Terminal) error {
	_, err := db.Exec("INSERT INTO terminals (id, mac, cpu_usage, mem_usage, last_modify_time) VALUES ($1, $2, $3, $4, $5)",
		terminal.ID, terminal.MAC, terminal.CPUUsage, terminal.MemUsage, terminal.LastModifyTime)
	return err
}

func AgentLogChange(changeLog *ChangeLog) error {
	_, err := db.Exec("INSERT INTO change_logs (terminal_id, old_mac, new_mac, timestamp, last_modify_time) VALUES ($1, $2, $3, $4, $5)",
		changeLog.TerminalID, changeLog.OldMAC, changeLog.NewMAC, changeLog.Timestamp, changeLog.LastModifyTime)
	return err
}
