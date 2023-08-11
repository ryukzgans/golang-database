package repository

import (
	"context"
	"database/sql"
	"errors"
	"golang-database/entity"
	"strconv"
)

type commentRepositoryImpl struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepositoryImpl{DB: db}
}

func (repo *commentRepositoryImpl) Insert(ctx context.Context, comment entity.Comments) (entity.Comments, error) {

	querySql := "INSERT INTO comments(email, comments) VALUES (?, ?)"
	result, err := repo.DB.ExecContext(ctx, querySql, comment.Email, comment.Comment)
	if err != nil {
		return comment, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return comment, err
	}

	comment.Id = int32(id)
	return comment, nil

}
func (repo *commentRepositoryImpl) FindById(ctx context.Context, id int32) (entity.Comments, error) {
	querySql := "SELECT id, email, comments FROM comments WHERE id = ? LIMIT 1"
	rows, err := repo.DB.QueryContext(ctx, querySql, id)
	comment := entity.Comments{}
	if err != nil {
		return comment, err
	}
	defer rows.Close()

	if rows.Next() {
		// ada
		rows.Scan(&comment.Id, &comment.Email, &comment.Comment)
		return comment, nil

	} else {
		// tidak ada
		return comment, errors.New("ID: " + strconv.Itoa(int(id)) + " Not Found")
	}
}

func (repo *commentRepositoryImpl) FindAll(ctx context.Context) ([]entity.Comments, error) {
	querySql := "SELECT id, email, comments FROM comments"
	rows, err := repo.DB.QueryContext(ctx, querySql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []entity.Comments
	for rows.Next() {
		comment := entity.Comments{}
		rows.Scan(&comment.Id, &comment.Email, &comment.Comment)
		comments = append(comments, comment)
	}
	return comments, nil
}
