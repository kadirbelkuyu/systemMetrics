package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"systemMetric/config"
)

func ConnectPostgres(c *config.Config) (*sql.DB, error) {
	DSN := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Dbname,
		c.Database.Password,
	)

	db, err := sql.Open("postgres", DSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
