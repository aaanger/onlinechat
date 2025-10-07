package users

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewUserRepository(db *sql.DB, logger *logrus.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) CreateUser(user User) (*User, error) {
	query := `
		INSERT INTO users (email, username, password_hash, created_at, updated_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at, last_seen, is_active
	`

	now := time.Now()
	row := r.db.QueryRow(query, user.Email, user.Username, user.Password, now, now, true)

	err := row.Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
		&user.LastSeen, &user.IsActive,
	)
	if err != nil {
		r.logger.WithError(err).Error("Failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, username, password_hash, created_at, updated_at, 
		       last_seen, is_active
		FROM users 
		WHERE email = $1 AND is_active = true
	`

	user := &User{}
	row := r.db.QueryRow(query, email)

	err := row.Scan(
		&user.ID, &user.Email, &user.Username, &user.Password,
		&user.CreatedAt, &user.UpdatedAt, &user.LastSeen,
		&user.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		r.logger.WithError(err).Error("Failed to get user by email")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByID(id int) (*User, error) {
	query := `
		SELECT id, email, username, password_hash, created_at, updated_at, 
		       last_seen, is_active
		FROM users 
		WHERE id = $1 AND is_active = true
	`

	user := &User{}
	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&user.ID, &user.Email, &user.Username, &user.Password,
		&user.CreatedAt, &user.UpdatedAt, &user.LastSeen,
		&user.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		r.logger.WithError(err).Error("Failed to get user by ID")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(username string) (*User, error) {
	query := `
		SELECT id, email, username, password_hash, created_at, updated_at, 
		       last_seen, is_active
		FROM users 
		WHERE username = $1 AND is_active = true
	`

	user := &User{}
	row := r.db.QueryRow(query, username)

	err := row.Scan(
		&user.ID, &user.Email, &user.Username, &user.Password,
		&user.CreatedAt, &user.UpdatedAt, &user.LastSeen,
		&user.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		r.logger.WithError(err).Error("Failed to get user by username")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(id int, update UserUpdate) (*User, error) {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if update.Username != nil {
		setParts = append(setParts, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, *update.Username)
		argIndex++
	}

	if len(setParts) == 0 {
		return r.GetByID(id)
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)

	setClause := ""
	for i, part := range setParts {
		if i > 0 {
			setClause += ", "
		}
		setClause += part
	}

	query := fmt.Sprintf(`
		UPDATE users 
		SET %s 
		WHERE id = $%d AND is_active = true
		RETURNING id, email, username, password_hash, created_at, updated_at, 
		          last_seen, is_active
	`, setClause, argIndex)

	user := &User{}
	row := r.db.QueryRow(query, args...)

	err := row.Scan(
		&user.ID, &user.Email, &user.Username, &user.Password,
		&user.CreatedAt, &user.UpdatedAt, &user.LastSeen,
		&user.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		r.logger.WithError(err).Error("Failed to update user")
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) UpdateLastSeen(id int) error {
	query := `UPDATE users SET last_seen = $1 WHERE id = $2 AND is_active = true`

	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		r.logger.WithError(err).Error("Failed to update last seen")
		return fmt.Errorf("failed to update last seen: %w", err)
	}

	return nil
}

func (r *UserRepository) DeleteUser(id int) error {
	query := `UPDATE users SET is_active = false, updated_at = $1 WHERE id = $2`

	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		r.logger.WithError(err).Error("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = $1 AND is_active = true`

	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		r.logger.WithError(err).Error("Failed to check email existence")
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

func (r *UserRepository) UsernameExists(username string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = $1 AND is_active = true`

	var count int
	err := r.db.QueryRow(query, username).Scan(&count)
	if err != nil {
		r.logger.WithError(err).Error("Failed to check username existence")
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return count > 0, nil
}
