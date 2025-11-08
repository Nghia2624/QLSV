package postgres

import (
	"database/sql"
	"qlsvgo/internal/domain/model"
)

type UserPostgresRepository struct {
	DB *sql.DB
}

func (r *UserPostgresRepository) Create(user *model.User) error {
	var refID interface{}
	if user.RefID == "" {
		refID = nil
	} else {
		refID = user.RefID
	}

	_, err := r.DB.Exec(`INSERT INTO users (id, username, email, password_hash, role, ref_id) VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, user.Username, user.Email, user.Password, user.Role, refID)
	return err
}

func (r *UserPostgresRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	var refID sql.NullString
	err := r.DB.QueryRow(`SELECT id, username, email, password_hash, role, ref_id FROM users WHERE username = $1`, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &refID)
	if err != nil {
		return nil, err
	}
	if refID.Valid {
		user.RefID = refID.String
	} else {
		user.RefID = ""
	}
	return &user, nil
}

func (r *UserPostgresRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	var refID sql.NullString
	err := r.DB.QueryRow(`SELECT id, username, email, password_hash, role, ref_id FROM users WHERE id = $1`, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &refID)
	if err != nil {
		return nil, err
	}
	if refID.Valid {
		user.RefID = refID.String
	} else {
		user.RefID = ""
	}
	return &user, nil
}

func (r *UserPostgresRepository) GetAll() ([]*model.User, error) {
	rows, err := r.DB.Query(`SELECT id, username, email, password_hash, role, ref_id FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		var refID sql.NullString
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &refID)
		if err != nil {
			return nil, err
		}
		if refID.Valid {
			user.RefID = refID.String
		} else {
			user.RefID = ""
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *UserPostgresRepository) Update(user *model.User) error {
	var refID interface{}
	if user.RefID == "" {
		refID = nil
	} else {
		refID = user.RefID
	}

	_, err := r.DB.Exec(`UPDATE users SET username = $2, email = $3, password_hash = $4, role = $5, ref_id = $6 WHERE id = $1`,
		user.ID, user.Username, user.Email, user.Password, user.Role, refID)
	return err
}

func (r *UserPostgresRepository) Delete(id string) error {
	_, err := r.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}
