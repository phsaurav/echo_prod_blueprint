package repository

import (
	"database/sql"

	"github.com/phsaurav/go_echo_base/internal/post"
	"github.com/phsaurav/go_echo_base/internal/user"
)

type Storage struct {
	Posts post.PostRepo
	Users user.UserRepo
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: post.NewRepository(db),
		Users: user.NewRepository(db),
		//Comments:  comment.NewRepository(db),
		//Followers: follower.NewRepository(db),
		//Roles:     role.NewRepository(db),
	}
}
