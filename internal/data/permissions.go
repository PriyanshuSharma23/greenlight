package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	for i := range p {
		if p[i] == code {
			return true
		}
	}

	return false
}

type PermissionModel struct {
	DB *sql.DB
}

func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	stmt := `
          SELECT permissions.code FROM user_permissions
          INNER JOIN permissions ON permissions.id = user_permissions.permission_id
          INNER JOIN users ON users.id = user_permissions.user_id
          WHERE user_permissions.user_id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r, err := m.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}

	var permissions Permissions

	for r.Next() {
		var permission string
		err := r.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = r.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (m PermissionModel) AddForUser(userID int64, codes ...string) error {
	stmt := `
          INSERT INTO user_permissions (user_id, permission_id)
          SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)
          `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, stmt, userID, pq.Array(codes))
	return err
}
