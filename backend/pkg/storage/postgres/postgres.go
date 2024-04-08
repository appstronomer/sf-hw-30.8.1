package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Служит ООП-абстракцией для удобной работы с задачами и метками
type Storage struct {
	db *pgxpool.Pool
}

// Создание нового хранилища Storage
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// ООП-абстракция над конкретной меткой
type Label struct {
	ID   int
	Name string
}

// Создание новой метки
func (s *Storage) NewLabel(label Label) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO labels (name)
		VALUES ($1) RETURNING id;
		`,
		label.Name,
	).Scan(&id)
	return id, err
}

// ООП-абстракция над конкретной задачей
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Создание новой задачи
func (s *Storage) NewTask(task Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (
			opened, 
			closed, 
			author_id, 
			assigned_id, 
			title, 
			content
		)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
		`,
		task.Opened,
		task.Closed,
		task.AuthorID,
		task.AssignedID,
		task.Title,
		task.Content,
	).Scan(&id)
	return id, err
}

// Получение задачи по ID
func (s *Storage) Task(taskID int) (Task, error) {
	row := s.db.QueryRow(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE id = $1;
	`,
		taskID,
	)
	var task Task
	err := row.Scan(
		&task.ID,
		&task.Opened,
		&task.Closed,
		&task.AuthorID,
		&task.AssignedID,
		&task.Title,
		&task.Content,
	)
	return task, err
}

// Фильтрация задач по автору и/или метке. Если в качестве authorID и/или label будет
// передано значение по умолчанию (0 и/или "" соответственно), то этот параметр(ы) не
// будет(-ут) участвовать в фильтрации.
func (s *Storage) Tasks(authorID int, label string) ([]Task, error) {
	var rows pgx.Rows
	var err error
	if label != "" {
		rows, err = s.queryTasksLabel(authorID, label)
	} else {
		rows, err = s.queryTasks(authorID)
	}
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)

	}
	return tasks, rows.Err()
}

// Удуление задачи по ID.
func (s *Storage) DeleteTask(taskID int) error {
	tag, err := s.db.Exec(context.Background(), `
		DELETE 
		FROM tasks
		WHERE id = $1;
	`,
		taskID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delete task id=%v: no rows affected", taskID)
	}
	return nil
}

// Обновление задачи по свойству ID структуры Task.
func (s *Storage) UpdateTask(task Task) error {
	tag, err := s.db.Exec(context.Background(), `
		UPDATE tasks
		SET
			opened = $1,
			closed = $2,
			author_id = $3,
			assigned_id = $4,
			title = $5,
			content = $6
		WHERE id = $7;
	`,
		task.Opened,
		task.Closed,
		task.AuthorID,
		task.AssignedID,
		task.Title,
		task.Content,
		task.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("update task id=%v : no rows affected", task.ID)
	}
	return nil
}

// Добавление связи задача-метка
func (s *Storage) TaskAddLabel(taskID int, labelID int) error {
	tag, err := s.db.Exec(context.Background(), `
		INSERT INTO tasks_labels (task_id, label_id) 
		VALUES ($1, $2);
	`,
		taskID,
		labelID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("task id=%v add label id=%v: no rows affected", taskID, labelID)
	}
	return nil
}

// Вспомогательный метод, формирующий выборку с фильтрацией только по автору
func (s *Storage) queryTasks(authorID int) (pgx.Rows, error) {
	return s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR author_id = $1)
		ORDER BY id;
	`,
		authorID,
	)
}

// Вспомогательный метод, формирующий выборку с фильтрацией по автору и/или метке
func (s *Storage) queryTasksLabel(authorID int, label string) (pgx.Rows, error) {
	return s.db.Query(context.Background(), `
		SELECT 
			t.id,
			t.opened,
			t.closed,
			t.author_id,
			t.assigned_id,
			t.title,
			t.content
		FROM tasks t
		LEFT JOIN tasks_labels tl ON t.id = tl.task_id
		LEFT JOIN labels l ON tl.label_id = l.id
		WHERE
			($1 = 0 OR t.author_id = $1) AND
			($2 = '' OR l.name = $2)
		ORDER BY id;
	`,
		authorID,
		label,
	)
}
