package users

import (
	"context"
	"database/sql"
	"errors"
	"social-media-app/internal/shared/consts"
	usersDTO "social-media-app/internal/shared/dto/users"
	mainErrors "social-media-app/internal/shared/errors/main"
	usersErrors "social-media-app/internal/shared/errors/users"
	dbUtils "social-media-app/internal/shared/utils"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) repository {
	return repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (first_name, last_name, password, email)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, consts.QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Firstname,
		user.Lastname,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return usersErrors.ErrDuplicateEmail
		}

		if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
			return usersErrors.ErrDuplicateUsername
		}

		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, userID int64) (User, error) {
	query := `
		SELECT id, COALESCE(first_name, ''), COALESCE(last_name, ''), email, username, created_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, consts.QueryTimeoutDuration)
	defer cancel()

	var user User

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Firstname,
		&user.Lastname,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, mainErrors.ErrNotFound
		}

		return User{}, err
	}

	return user, nil
}

func (r *repository) Follow(ctx context.Context, payload usersDTO.FollowPayload) error {
	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, follower_id) DO NOTHING
	`

	ctx, cancel := context.WithTimeout(ctx, consts.QueryTimeoutDuration)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, payload.UserID, payload.FollowerID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UnFollow(ctx context.Context, payload usersDTO.FollowPayload) error {
	query := `
		DELETE FROM followers
		WHERE user_id = $1 
		AND follower_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, consts.QueryTimeoutDuration)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, payload.UserID, payload.FollowerID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CreateAndInvite(ctx context.Context, payload *User, token string, invitationExp time.Duration) error {
	err := dbUtils.WithTransaction(r.db, ctx, func(tx *sql.Tx) error {
		if err := r.Create(ctx, tx, payload); err != nil {
			return err
		}

		err := r.createUserInvitation(ctx, tx, token, invitationExp, payload.ID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) createUserInvitation(
	ctx context.Context,
	tx *sql.Tx,
	token string,
	invitationExp time.Duration,
	userID int64,
) error {
	query := `
		INSERT INTO user_invitations (token, user_id, expiry)
		VALUES ($1, $2, $3)	
	`

	ctx, cancel := context.WithTimeout(ctx, consts.QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}

	return nil
}
