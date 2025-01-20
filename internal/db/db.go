package db

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
    bun.BaseModel `bun:"table:users,alias:u"`
    ID           int64     `bun:"id,pk,autoincrement"`
    Username     string    `bun:"username,unique,notnull"`
    PasswordHash string    `bun:"password,notnull"`
    CreatedAt    time.Time `bun:"created_at,nullzero,default:current_timestamp"`
    UpdatedAt    time.Time `bun:"updated_at,nullzero,default:current_timestamp"`
    DeletedAt    time.Time `bun:"deleted_at,soft_delete"`
}

type Session struct {
    bun.BaseModel `bun:"table:sessions,alias:s"`
    ID            int64     `bun:"id,pk,autoincrement"`
    AccessToken         string    `bun:"access_token,unique,notnull"`
    RefreshToken        string    `bun:"refresh_token,unique,notnull"`
    UserID        int64     `bun:"user_id,notnull"`
    CreatedAt     time.Time `bun:"created_at,nullzero,default:current_timestamp"`
    UpdatedAt     time.Time `bun:"updated_at,nullzero,default:current_timestamp"`
}