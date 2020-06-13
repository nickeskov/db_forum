package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/utils/database/driver/pgx/codes"
	sqlHelpers "github.com/nickeskov/db_forum/pkg/sql"
	"github.com/pkg/errors"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return Repository{
		db: db,
	}
}

func (repo Repository) Create(user models.User) error {
	ctx := context.Background()

	_, err := repo.db.Exec(ctx,
		`	INSERT INTO users (nickname, email, fullname, about) 
				VALUES ($1, $2, $3, $4)`,
		user.Nickname,
		user.Email,
		user.Fullname,
		user.About,
	)

	pgxErr := codes.ExtractPgx4ErrorCode(err)
	if pgxErr != nil && pgxErr.Error() == codes.ErrCodeUnique {
		return models.ErrAlreadyExist
	}

	return errors.WithStack(err)
}

func (repo Repository) UpdateByNickname(user models.User) (models.User, error) {
	ctx := context.Background()

	row := repo.db.QueryRow(ctx,
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
	)

	if err := scanUser(row, &user); err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, models.ErrDoesNotExist
		}

		if pgxErr := codes.ExtractPgx4ErrorCode(err); pgxErr != nil {
			switch pgxErr.Error() {
			case codes.ErrCodeUnique:
				return models.User{}, models.ErrConflict
			default:
				return models.User{}, errors.Wrapf(err,
					"some error while updating user by nickname, user=%+v", user)
			}
		}
	}

	return user, nil
}

func (repo Repository) GetByNickname(nickname string) (user models.User, err error) {
	ctx := context.Background()

	row := repo.db.QueryRow(ctx,
		`	SELECT 	nickname,
						email,
						fullname,
						about
				FROM users
				WHERE nickname = $1`,
		nickname,
	)

	err = scanUser(row, &user)

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
	ctx := context.Background()

	rows, err := repo.db.Query(ctx,
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

		if err := scanUser(rows, &user); err != nil {
			return nil, errors.Wrapf(err, "error in user repository "+
				"GetWithSameNicknameAndEmail while scanning, nickname=%s, email=%s: %v",
				nickname, email, err)
		}

		users = append(users, user)
	}

	return users, nil
}

func scanUser(scanner sqlHelpers.Scanner, userDst *models.User) error {
	err := scanner.Scan(
		&userDst.Nickname,
		&userDst.Email,
		&userDst.Fullname,
		&userDst.About,
	)
	return err
}
