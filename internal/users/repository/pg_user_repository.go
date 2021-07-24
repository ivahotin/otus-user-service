package repository

import (
	"context"

	"example.com/arch/user-service/internal/users"
	"example.com/arch/user-service/internal/users/models"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	getByIdQuery    = "select id, username, firstname, lastname, email, phone from users where id = $1"
	insertUserQuery = "insert into users (username, firstname, lastname, email, phone) values ($1, $2, $3, $4, $5) returning id"
	deleteUserQuery = "delete from users where id = $1"
	updateUserQuery = "update users set username = $1, firstname = $2, lastname = $3, email = $4, phone = $5 where id = $6"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (repo *PostgresUserRepository) CreateUser(user *models.User) (models.UserId, error) {
	row := repo.db.QueryRow(
		context.Background(),
		insertUserQuery,
		user.Username,
		user.Firstname,
		user.Lastname,
		user.Email,
		user.Phone,
	)
	var userId models.UserId
	err := row.Scan(&userId)
	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (repo *PostgresUserRepository) GetUser(userId models.UserId) (models.User, error) {
	var user models.User
	err := pgxscan.Get(context.Background(), repo.db, &user, getByIdQuery, userId)
	if err != nil {
		if pgxscan.NotFound(err) {
			return models.User{}, users.ErrNotFound
		}

		return models.User{}, err
	}
	return user, nil
}

func (repo *PostgresUserRepository) UpdateUser(user *models.User) error {
	tag, err := repo.db.Exec(
		context.Background(),
		updateUserQuery,
		user.Username,
		user.Firstname,
		user.Lastname,
		user.Email,
		user.Phone,
		user.Id,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return users.ErrNotFound
	}
	return nil
}

func (repo *PostgresUserRepository) DeleteUser(userId models.UserId) error {
	tag, err := repo.db.Exec(context.Background(), deleteUserQuery, userId)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return users.ErrNotFound
	}
	return nil
}
