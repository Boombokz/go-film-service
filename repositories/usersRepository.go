package repositories

import (
	"context"
	"filmservice/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository struct {
	db *pgxpool.Pool
}

func NewUsersRepository(conn *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{db: conn}
}

func (r *UsersRepository) FindByEmail(c *gin.Context, email string) (models.User, error) {
	var user models.User
	row := r.db.QueryRow(c, "select id, name, email, password_hash from users where email = $1", email)

	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UsersRepository) FindAll(c *gin.Context) ([]models.User, error) {
	rows, err := r.db.Query(c, "select id, name, email, password_hash from users order by id")
	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersRepository) FindById(c context.Context, id int) (models.User, error) {
	var user models.User
	row := r.db.QueryRow(c, "select id, name, email from users where id = $1", id)

	err := row.Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UsersRepository) Create(c *gin.Context, user models.User) (int, error) {
	var id int

	err := r.db.QueryRow(c, "insert into users (name, email, password_hash) values ($1, $2, $3) returning id",
		user.Name, user.Email, user.PasswordHash).Scan(&id)

	return id, err
}

func (r *UsersRepository) Update(c *gin.Context, updatedUser models.User) error {
	_, err := r.db.Exec(c, "update users set name = $1, email = $2 where id = $3", updatedUser.Name,
		updatedUser.Email, updatedUser.Id)

	return err
}

func (r *UsersRepository) ChangePassword(c *gin.Context, user models.User) error {
	_, err := r.db.Exec(c, "update users set password_hash = $1 where id = $2", user.PasswordHash, user.Id)

	return err
}

func (r *UsersRepository) Delete(c context.Context, id int) error {
	_, err := r.db.Exec(c, "delete from users where id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
