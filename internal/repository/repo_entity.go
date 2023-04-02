package repository

type repoTask struct {
	Id          string `db:"taskId"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Status      int16  `db:"status"`
	CreatedAt   int64  `db:"createdAt"`
	UpdatedAt   *int64 `db:"updatedAt"`
	IsDeleted   int    `db:"isDeleted"`
}
