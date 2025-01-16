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

type List struct {
	bun.BaseModel `bun:"table:lists,alias:l"`

	ID           int64     `bun:"id,pk,autoincrement"`
	
}
