package db

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Username      string    `bun:"username,unique,notnull"`
	PasswordHash  string    `bun:"password,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,default:current_timestamp"`
	DeletedAt     time.Time `bun:"deleted_at,soft_delete"`
}

type Session struct {
	bun.BaseModel `bun:"table:sessions,alias:s"`
	ID            int64     `bun:"id,pk,autoincrement"`
	AccessToken   string    `bun:"access_token,unique,notnull"`
	RefreshToken  string    `bun:"refresh_token,unique,notnull"`
	UserID        int64     `bun:"user_id,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,default:current_timestamp"`
}

type List struct {
	bun.BaseModel `bun:"table:todos,alias:l"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Name          string    `bun:"name,notnull"`
	Users         []*User   `bun:"rel:has-many:join:id=list_id"`
	CreatedAt     time.Time `bun:"created_at,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,default:current_timestamp"`
}

type Tag struct {
	bun.BaseModel `bun:"table:tags,alias:t"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Name          string    `bun:"name,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,default:current_timestamp"`
}

type ToDoStatus string

const (
	TODO  ToDoStatus = "todo"
	DOING ToDoStatus = "doing"
	DONE  ToDoStatus = "done"
)

type Todo struct {
	bun.BaseModel `bun:"table:todo_items,alias:i"`
	ID            int64      `bun:"id,pk,autoincrement"`
	Title         string     `bun:"title,notnull"`
	Description   string     `bun:"description,notnull"`
	Status        ToDoStatus `bun:"status,notnull"`
	CreatedAt     time.Time  `bun:"created_at,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time  `bun:"updated_at,nullzero,default:current_timestamp"`

	ListID int64 `bun:"list_id,notnull"`
	Tags   []Tag `bun:"m2m:todo_tags,join:Todo=Tag"`
}

type TodoTag struct {
	bun.BaseModel `bun:"table:todo_tags,alias:i"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Todo          *Todo     `bun:"rel:belongs-to,join:id=todo_id"`
	TodoID        int64     `bun:"todo_id,notnull"`
	List          string    `bun:"list_id,notnull"`
	Tag           *Tag      `bun:"rel:belongs-to,join:id=tag_id"`
	TagID         int64     `bun:"tag_id,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,default:current_timestamp"`
}

type DB struct {
	*bun.DB
}
