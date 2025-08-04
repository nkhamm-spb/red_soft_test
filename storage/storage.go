package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"

	"github.com/nkhamm-spb/red_soft_test/config"
	"github.com/nkhamm-spb/red_soft_test/schemas"
)

type StorageInterface interface {
	GetUserById(ctx context.Context, id int) (*schemas.User, error)
	GetUserBySurname(ctx context.Context, surname string) (*schemas.User, error)
	AddUser(ctx context.Context, user *schemas.User) (*schemas.User, error)
	GetAll(ctx context.Context) ([]schemas.User, error)
	EditUser(ctx context.Context, id int, editData map[string]interface{}) (*schemas.User, error)
}

type Storage struct {
	db *sql.DB
}

func New(ctx context.Context, config *config.Storage) (*Storage, error) {
	db, err := sql.Open("postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			config.User, config.Password, config.Name))
	if err != nil {
		return nil, fmt.Errorf("Error open db: %v", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("Error ping database: %v", err)
	}

	createTables := `
		CREATE TABLE IF NOT EXISTS users (
			id           SERIAL PRIMARY KEY,
			name         TEXT NOT NULL,
			surname      TEXT NOT NULL,
			age          INT NOT NULL,
			gender       TEXT NOT NULL,
			nationalize  TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS emails (
			user_id    INT NOT NULL,
			email      TEXT NOT NULL
		);`

	if _, err := db.ExecContext(ctx, createTables); err != nil {
		return nil, fmt.Errorf("Error create table: %v", err)
	}

	log.Printf("Connected to db!")

	return &Storage{db}, nil
}

func (storage *Storage) GetUserById(ctx context.Context, id int) (*schemas.User, error) {
	user := schemas.User{}

	err := storage.db.QueryRowContext(ctx,
		`SELECT id, name, surname, age, gender, nationalize FROM users WHERE id = $1;`,
		id).Scan(&user.ID, &user.Name, &user.Surname, &user.Age, &user.Gender, &user.Nationalize)
	if err != nil {
		return nil, fmt.Errorf("Error query: %v", err)
	}

	rows, err := storage.db.QueryContext(ctx,
		`SELECT email FROM emails WHERE user_id = $1;`,
		id)
	if err != nil {
		return nil, fmt.Errorf("Error query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		user.Emails = append(user.Emails, email)
	}

	return &user, nil
}

func (storage *Storage) GetUserBySurname(ctx context.Context, surname string) (*schemas.User, error) {
	user := schemas.User{}

	err := storage.db.QueryRowContext(ctx,
		`SELECT id, name, surname, age, gender, nationalize FROM users WHERE surname = $1;`,
		surname).Scan(&user.ID, &user.Name, &user.Surname, &user.Age, &user.Gender, &user.Nationalize)
	if err != nil {
		return nil, fmt.Errorf("Error query: %v", err)
	}

	rows, err := storage.db.QueryContext(ctx,
		`SELECT email FROM emails WHERE user_id = $1;`,
		user.ID)
	if err != nil {
		return nil, fmt.Errorf("Error query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		user.Emails = append(user.Emails, email)
	}

	return &user, nil
}

func (storage *Storage) AddUser(ctx context.Context, user *schemas.User) (*schemas.User, error) {
	err := storage.db.QueryRowContext(ctx,
		`INSERT INTO users (name, surname, age, gender, nationalize) VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		user.Name, user.Surname, user.Age, user.Gender, user.Nationalize).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("Error query: %v", err)
	}

	for _, email := range user.Emails {
		_, err := storage.db.ExecContext(ctx,
			`INSERT INTO emails (user_id, email) VALUES ($1, $2);`,
			user.ID, email)

		if err != nil {
			return nil, fmt.Errorf("Error query: %v", err)
		}
	}

	return user, nil
}

func (storage *Storage) GetAll(ctx context.Context) ([]schemas.User, error) {
	users := make([]schemas.User, 0)

	rows, err := storage.db.QueryContext(ctx,
		`SELECT id, name, surname, age, gender, nationalize FROM users`)
	if err != nil {
		return nil, fmt.Errorf("Error query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := schemas.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Age, &user.Gender, &user.Nationalize); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	for i := 0; i < len(users); i++ {
		emailRows, err := storage.db.QueryContext(ctx,
			`SELECT email FROM emails WHERE user_id = $1;`,
			users[i].ID)
		if err != nil {
			return nil, fmt.Errorf("Error query: %v", err)
		}

		for emailRows.Next() {
			var email string
			if err := emailRows.Scan(&email); err != nil {
				return nil, err
			}
			users[i].Emails = append(users[i].Emails, email)
		}
		emailRows.Close()
	}

	return users, nil
}

func (storage *Storage) EditUser(ctx context.Context, id int, editData map[string]interface{}) (*schemas.User, error) {
	tx, err := storage.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("Error begin transaction: %v", err)
	}
	defer tx.Rollback()

	if emails, ok := editData["Emails"]; ok {
		_, err = tx.ExecContext(ctx,
			`DELETE FROM emails WHERE user_id = $1;`,
			id)
		if err != nil {
			return nil, fmt.Errorf("Error exec: %v", err)
		}

		listEmails, ok := emails.([]interface{})

		if !ok {
			return nil, fmt.Errorf("Wrong format for Emails")
		}

		for _, email := range listEmails {
			stringEmail, ok := email.(string)
			if !ok {
				return nil, fmt.Errorf("Wrong format for Emails")
			}

			_, err := tx.ExecContext(ctx,
				`INSERT INTO emails (user_id, email) VALUES ($1, $2);`,
				id, stringEmail)

			if err != nil {
				return nil, fmt.Errorf("Error exec: %v", err)
			}
		}
	}

	var updates []string
	var args []interface{}

	queryCounter := 1

	if name, ok := editData["name"]; ok {
		stringName, ok := name.(string)

		if !ok {
			return nil, fmt.Errorf("Wrong format for name")
		}

		updates = append(updates, fmt.Sprintf("name = $%d", queryCounter))
		args = append(args, stringName)
		queryCounter++
	}

	if surname, ok := editData["surname"]; ok {
		stringSurname, ok := surname.(string)

		if !ok {
			return nil, fmt.Errorf("Wrong format for surname")
		}

		updates = append(updates, fmt.Sprintf("surname = $%d", queryCounter))
		args = append(args, stringSurname)
		queryCounter++
	}

	if gender, ok := editData["gender"]; ok {
		stringGender, ok := gender.(string)

		if !ok {
			return nil, fmt.Errorf("Wrong format for gender")
		}

		updates = append(updates, fmt.Sprintf("gender = $%d", queryCounter))
		args = append(args, stringGender)
		queryCounter++
	}

	if age, ok := editData["age"]; ok {
		floatAge, ok := age.(float64)

		if !ok {
			return nil, fmt.Errorf("Wrong format for age")
		}

		updates = append(updates, fmt.Sprintf("age = $%d", queryCounter))
		args = append(args, int(floatAge))
		queryCounter++
	}

	if nationalize, ok := editData["nationalize"]; ok {
		stringNationalize, ok := nationalize.(string)

		if !ok {
			return nil, fmt.Errorf("Wrong format for nationalize")
		}

		updates = append(updates, fmt.Sprintf("nationalize = $%d", queryCounter))
		args = append(args, stringNationalize)
		queryCounter++
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d;", strings.Join(updates, ", "), queryCounter)
	args = append(args, id)
	log.Println(query)

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Error exec: %v", err)
	}

	user := schemas.User{}

	// Дублирование так как нужно получать пользователя через транзакцию.
	// TODO: Убрать дублирование с помощью generic
	err = tx.QueryRowContext(ctx,
		`SELECT id, name, surname, age, gender, nationalize FROM users WHERE id = $1;`,
		id).Scan(&user.ID, &user.Name, &user.Surname, &user.Age, &user.Gender, &user.Nationalize)
	if err != nil {
		return nil, fmt.Errorf("Error query: %v", err)
	}

	rows, err := tx.QueryContext(ctx,
		`SELECT email FROM emails WHERE user_id = $1;`,
		id)
	if err != nil {
		return nil, fmt.Errorf("Error query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		user.Emails = append(user.Emails, email)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("Error commit: %v", err)
	}
	log.Printf("transaction committed")

	return &user, nil
}
