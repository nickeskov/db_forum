package repository

import (
	"github.com/jackc/pgx"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/utils/database/driver/pgx/codes"
	"github.com/pkg/errors"
)

type Repository struct {
	db *pgx.ConnPool
}

func NewRepository(db *pgx.ConnPool) Repository {
	return Repository{
		db: db,
	}
}

func (repo Repository) Create(user models.User) error {
	_, err := repo.db.Exec(
		`	INSERT INTO users (nickname, email, fullname, about) 
				VALUES ($1, $2, $3, $4)`,
		user.Nickname,
		user.Email,
		user.Fullname,
		user.About,
	)

	pgxErr := codes.ExtractErrorCode(err)
	if pgxErr != nil && pgxErr.Error() == codes.ErrCodeUnique {
		err = models.ErrAlreadyExist
	}

	return err
}

func (repo Repository) UpdateByNickname(user models.User) (models.User, error) {
	err := repo.db.QueryRow(
		`	UPDATE users
				SET email=COALESCE(NULLIF($2, ''), email),
					fullname=COALESCE(NULLIF($3, ''), fullname),
					about=COALESCE(NULLIF($4, ''), about)
				WHERE nickname = $1
				RETURNING nickname, email, fullname, about`,
		user.Nickname,
		user.Email,
		user.Fullname,
		user.About,
	).Scan(
		&user.Nickname,
		&user.Email,
		&user.Fullname,
		&user.About,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, models.ErrDoesNotExist
		}

		if pgxErr := codes.ExtractErrorCode(err); pgxErr != nil {
			err = models.ErrConflict
		}
	}

	return user, err
}

func (repo Repository) GetByNickname(nickname string) (user models.User, err error) {
	err = repo.db.QueryRow(
		`	SELECT 	nickname,
						email,
						fullname,
						about
				FROM users
				WHERE nickname = $1`,
		nickname,
	).Scan(
		&user.Nickname,
		&user.Email,
		&user.Fullname,
		&user.About,
	)

	switch {
	case err == pgx.ErrNoRows:
		return user, models.ErrDoesNotExist
	case err != nil:
		return user, errors.Wrapf(err,
			"error in user repository GetByName, nickname=%s: %v", nickname, err)
	}

	return user, nil
}

func (repo Repository) GetWithSameNicknameAndEmail(nickname, email string) (users models.Users, err error) {
	rows, err := repo.db.Query(
		`	SELECT 	nickname,
						email,
						fullname,
						about
				FROM users
				WHERE nickname = $1 OR email = $2`,
		nickname, email,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error in user repository "+
			"GetWithSameNicknameAndEmail, nickname=%s, email=%s: %v", nickname, email, err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.Nickname,
			&user.Email,
			&user.Fullname,
			&user.About,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "error in user repository "+
				"GetWithSameNicknameAndEmail while scanning, nickname=%s, email=%s: %v",
				nickname, email, err)
		}

		users = append(users, user)
	}

	return users, nil
}
