package storage

import (
	"testing"
	"context"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/nkhamm-spb/red_soft_test/schemas"
)

func TestGetUserById(t *testing.T) {
	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp),
	)
	require.NoError(t, err)
	defer db.Close()
	storage := Storage {db: db}

	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT id, name, surname, age, gender, nationalize FROM users WHERE id = $1;`)).
		WithArgs(11).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "age", "gender", "nationalize"}).
										AddRow(11, "Test", "Testovich", 20, "Male", "Russian"),
		)

	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT email FROM emails WHERE user_id = $1;`)).
		WithArgs(11).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).
						AddRow("test_testovich@test.com"))

	got, err := storage.GetUserById(context.Background(), 11)
	require.NoError(t, err)
	require.Equal(t, schemas.User{ID: 11, Name: "Test", Surname: "Testovich",
								  Age: 20, Gender: "Male", Nationalize: "Russian",
								  Emails: []string{"test_testovich@test.com"},}, *got)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserBySurname(t *testing.T) {
	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp),
	)
	require.NoError(t, err)
	defer db.Close()
	storage := Storage {db: db}

	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT id, name, surname, age, gender, nationalize FROM users WHERE surname = $1;`)).
		WithArgs("Testovich").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "age", "gender", "nationalize"}).
										AddRow(11, "Test", "Testovich", 20, "Male", "Russian"),
		)

	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT email FROM emails WHERE user_id = $1;`)).
		WithArgs(11).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).
						AddRow("test_testovich@test.com"))

	got, err := storage.GetUserBySurname(context.Background(), "Testovich")
	require.NoError(t, err)
	require.Equal(t, schemas.User{ID: 11, Name: "Test", Surname: "Testovich",
								  Age: 20, Gender: "Male", Nationalize: "Russian",
								  Emails: []string{"test_testovich@test.com"},}, *got)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddUser(t *testing.T) {
	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp),
	)
	require.NoError(t, err)
	defer db.Close()
	storage := Storage {db: db}

	mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (name, surname, age, gender, nationalize) VALUES ($1, $2, $3, $4, $5) RETURNING id;`)).
		WithArgs("Test", "Testovich", 20, "Male", "Russian").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
										AddRow(11),
		)

	mock.
		ExpectExec(regexp.QuoteMeta(`INSERT INTO emails (user_id, email) VALUES ($1, $2);`)).
		WithArgs(11, "test_testovich@test.com").
		WillReturnResult(sqlmock.NewResult(0, 1))

	user := schemas.User{Name: "Test", Surname: "Testovich",
								  Age: 20, Gender: "Male", Nationalize: "Russian",
								  Emails: []string{"test_testovich@test.com"},}

	got, err := storage.AddUser(context.Background(), &user)
	require.NoError(t, err)
	require.Equal(t, user, *got)

	require.NoError(t, mock.ExpectationsWereMet())
}
