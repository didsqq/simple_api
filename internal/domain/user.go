package domain

type User struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

type Conditions struct {
	Limit  int // определенное количество строк
	Offset int // с какой строки начинать выборку
}

type UpdateUserInput struct {
	ID    int64   `db:"id"`
	Name  *string `db:"name"`
	Email *string `db:"email"`
}
