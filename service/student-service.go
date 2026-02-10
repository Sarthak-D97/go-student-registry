package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Sarthak-D97/go_stuAPI/internal/storage"
	"github.com/Sarthak-D97/go_stuAPI/internal/types"
	"github.com/redis/go-redis/v9"
)

const (
	studentKeyPrefix = "student:"
	studentListKey   = "students_list"
	cacheTTL         = 10 * time.Minute
)

type StudentService interface {
	CreateStudent(student types.Student) (types.Student, error)
	GetStudentByID(id int64) (*types.Student, error)
	GetAllStudents() ([]types.Student, error)
	UpdateStudent(id int64, student types.Student) error
	DeleteStudent(id int64) error
}

type studentService struct {
	storage storage.Storage
	rdb     *redis.Client
}

func NewStudentService(storage storage.Storage, rdb *redis.Client) StudentService {
	return &studentService{
		storage: storage,
		rdb:     rdb,
	}
}

func (s *studentService) CreateStudent(student types.Student) (types.Student, error) {
	lastID, err := s.storage.CreateStudent(student.Name, student.Email, student.Age)
	if err != nil {
		return types.Student{}, err
	}
	student.ID = int(lastID)

	go func(st types.Student, id int64) {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("%s%d", studentKeyPrefix, id)
		pipe := s.rdb.Pipeline()
		pipe.HSet(ctx, cacheKey, st)
		pipe.Expire(ctx, cacheKey, cacheTTL)
		pipe.Del(ctx, studentListKey)
		if _, err := pipe.Exec(ctx); err != nil {
			slog.Error("failed to cache student", "error", err)
		}
	}(student, lastID)

	slog.Info("student created successfully", slog.Int64("student_id", lastID))
	return student, nil
}

func (s *studentService) GetStudentByID(id int64) (*types.Student, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("%s%d", studentKeyPrefix, id)

	var cachedStudent types.Student
	err := s.rdb.HGetAll(ctx, cacheKey).Scan(&cachedStudent)
	if err == nil && cachedStudent.ID != 0 {
		slog.Info("serving student from cache (hash)", slog.Int64("id", id))
		return &cachedStudent, nil
	}

	student, err := s.storage.GetStudentById(id)
	if err != nil {
		return nil, err
	}

	go func(st *types.Student, cacheKey string) {
		ctx := context.Background()
		pipe := s.rdb.Pipeline()
		pipe.HSet(ctx, cacheKey, st)
		pipe.Expire(ctx, cacheKey, cacheTTL)
		if _, err := pipe.Exec(ctx); err != nil {
			slog.Error("failed to cache student", "error", err)
		}
	}(student, cacheKey)

	slog.Info("student fetched successfully", slog.Int64("student_id", id))
	return student, nil
}

func (s *studentService) GetAllStudents() ([]types.Student, error) {
	ctx := context.Background()

	val, err := s.rdb.Get(ctx, studentListKey).Result()
	if err == nil {
		var cachedStudents []types.Student
		if jsonErr := json.Unmarshal([]byte(val), &cachedStudents); jsonErr == nil {
			slog.Info("serving student list from cache")
			return cachedStudents, nil
		}
	}

	students, err := s.storage.GetAllStudents()
	if err != nil {
		return nil, err
	}

	go func(students []types.Student) {
		data, _ := json.Marshal(students)
		s.rdb.Set(context.Background(), studentListKey, data, cacheTTL)
	}(students)

	slog.Info("students fetched successfully", slog.Int("count", len(students)))
	return students, nil
}

func (s *studentService) UpdateStudent(id int64, student types.Student) error {
	if err := s.storage.UpdateStudent(id, student.Name, student.Email, student.Age); err != nil {
		return err
	}

	go func(id int64, st types.Student) {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("%s%d", studentKeyPrefix, id)

		pipe := s.rdb.Pipeline()
		pipe.HSet(ctx, cacheKey, st)
		pipe.Expire(ctx, cacheKey, cacheTTL)
		pipe.Del(ctx, studentListKey)
		if _, err := pipe.Exec(ctx); err != nil {
			slog.Error("failed to update student cache", "error", err)
		}
	}(id, student)

	slog.Info("student updated successfully", slog.Int64("student_id", id))
	return nil
}

func (s *studentService) DeleteStudent(id int64) error {
	if err := s.storage.DeleteStudent(id); err != nil {
		return err
	}

	go func(id int64) {
		ctx := context.Background()
		pipe := s.rdb.Pipeline()
		pipe.Del(ctx, fmt.Sprintf("%s%d", studentKeyPrefix, id))
		pipe.Del(ctx, studentListKey)
		if _, err := pipe.Exec(ctx); err != nil {
			slog.Error("failed to update student cache after delete", "error", err)
		}
	}(id)

	slog.Info("student deleted successfully", slog.Int64("student_id", id))
	return nil
}

