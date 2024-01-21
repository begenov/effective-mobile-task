package postgre

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/begenov/effective-mobile-task/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUsers(ctx context.Context, userID, limit, offset *int64, gender *int, nationality *string) ([]model.User, error) {
	var queryBuilder strings.Builder
	var dest []interface{}

	queryBuilder.WriteString(`
		select u.id, u.name, u.surname, u.patronymic, u.age,
		u.gender_id, u.nationality_id, n.name
		from "user" u 
		left join nationality n on n.id = u.nationality_id 
		left join gender g on g.id = u.gender_id 
		where (1=1)
	`)

	if userID != nil && *userID > 0 {
		queryBuilder.WriteString(" and u.id = $" + strconv.Itoa(len(dest)+1))
		dest = append(dest, *userID)
	}

	if gender != nil && (*gender == 1 || *gender == 2) {
		queryBuilder.WriteString(" and u.gender_id = $" + strconv.Itoa(len(dest)+1))
		dest = append(dest, *gender)
	}

	if nationality != nil && len(*nationality) > 0 {
		queryBuilder.WriteString(" and n.name = $" + strconv.Itoa(len(dest)+1))
		dest = append(dest, *nationality)
	}

	if limit != nil && *limit > 0 {
		queryBuilder.WriteString(" limit $" + strconv.Itoa(len(dest)+1))
		dest = append(dest, *limit)
	}

	if offset != nil && *offset > 0 {
		if limit != nil && *limit > 0 {
			calculatedOffset := *offset * *limit
			queryBuilder.WriteString(" offset $" + strconv.Itoa(len(dest)+1))
			dest = append(dest, calculatedOffset)
		}
	}

	rows, err := r.db.QueryContext(ctx, queryBuilder.String(), dest...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var nationality model.Nationality

		if err := rows.Scan(
			&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Age,
			&user.Gender,
			&nationality.ID, &nationality.Name,
		); err != nil {
			return nil, err
		}

		user.Nationality = &nationality
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) Create(ctx context.Context, user model.User) error {
	query := `
		insert into "user" (name, surname, patronymic, age, gender_id, nationality_id)
		values ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query, user.Name, user.Surname, user.Patronymic, user.Age,
		user.Gender, user.Nationality.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user model.User) error {
	query := `
		update "user"
		set name = $2, surname = $3, patronymic = $4, age = $5, 
		gender_id = $6, nationality_id = $7
		where id = $1
	`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Surname, user.Patronymic, user.Age,
		user.Gender, user.Nationality.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, userID int64) error {
	query := `
		delete from "user"
		where id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
