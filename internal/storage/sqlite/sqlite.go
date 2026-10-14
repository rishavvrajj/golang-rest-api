package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rishavvrajj/golang-rest-api/internal/models"
)

type SQLite struct {
	db *sql.DB
}

// New initializes a new SQLite storage connection and creates the students table if it doesn't exist.
func New(storagePath string) (*SQLite, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		age INTEGER NOT NULL
	);
	`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create students table: %w", err)
	}

	return &SQLite{db: db}, nil
}

func (s *SQLite) CreateStudent(name string, email string, age int) (int64, error) {
	query := "INSERT INTO students(name, email, age) VALUES(?, ?, ?)"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, fmt.Errorf("failed to execute insert statement: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

func (s *SQLite) GetStudentById(id int64) (models.Student, error) {
	query := "SELECT id, name, email, age FROM students WHERE id = ?"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return models.Student{}, fmt.Errorf("failed to prepare select statement: %w", err)
	}
	defer stmt.Close()

	var student models.Student
	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Student{}, fmt.Errorf("student with id %d not found", id)
		}
		return models.Student{}, fmt.Errorf("failed to query student by id: %w", err)
	}

	return student, nil
}

func (s *SQLite) GetStudents() ([]models.Student, error) {
	query := "SELECT id, name, email, age FROM students"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query students list: %w", err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age); err != nil {
			return nil, fmt.Errorf("failed to scan student row: %w", err)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return students, nil
}

func (s *SQLite) UpdateStudent(id int64, name string, email string, age int) error {
	query := "UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, email, age, id)
	if err != nil {
		return fmt.Errorf("failed to execute update statement: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("student with id %d not found", id)
	}

	return nil
}

func (s *SQLite) DeleteStudent(id int64) error {
	query := "DELETE FROM students WHERE id = ?"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to execute delete statement: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("student with id %d not found", id)
	}

	return nil
}
