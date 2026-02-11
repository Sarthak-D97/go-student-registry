package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Sarthak-D97/go_stuAPI/entity"
	"github.com/Sarthak-D97/go_stuAPI/repository"
	"github.com/redis/go-redis/v9"
)

const (
	studentKeyPrefix = "student:"
	studentListKey   = "students_list"
	cacheTTL         = 10 * time.Minute
)

// StudentService interface aligned with Controller calls
type StudentService interface {
	Create(student *entity.Student) error
	FindByID(id uint) (*entity.Student, error)
	FindAll() ([]entity.Student, error)
	Update(student *entity.Student) error
	Delete(id uint) error
}

type studentService struct {
	repo repository.Repository // Ensure your repository interface matches these types
	rdb  *redis.Client
}

// NewStudentService creates a new instance of the service
func NewStudentService(repo repository.Repository, rdb *redis.Client) StudentService {
	return &studentService{
		repo: repo,
		rdb:  rdb,
	}
}

// Create - Aligned to receive pointer
func (s *studentService) Create(student *entity.Student) error {
	// 1. Save to DB
	// We pass the pointer or value depending on your repo implementation.
	// Assuming Repo returns the created struct with ID.
	created, err := s.repo.Create(*student)
	if err != nil {
		return err
	}

	// Update the original pointer with the new ID so the controller can return it
	student.ID = created.ID

	// 2. Cache (Async)
	go func(st entity.Student) {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("%s%d", studentKeyPrefix, st.ID)

		pipe := s.rdb.Pipeline()
		// Use JSON for HSet if fields aren't flat, or use standard Set for simplicity
		// Here we stick to your HSet logic, ensure struct has redis tags
		pipe.HSet(ctx, cacheKey, st)
		pipe.Expire(ctx, cacheKey, cacheTTL)

		// Invalidate list cache
		pipe.Del(ctx, studentListKey)

		if _, err := pipe.Exec(ctx); err != nil {
			slog.Error("failed to cache student", "error", err)
		}
	}(created)

	slog.Info("student created successfully", slog.Uint64("student_id", uint64(created.ID)))
	return nil
}

// FindByID - Aligned to accept uint
func (s *studentService) FindByID(id uint) (*entity.Student, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("%s%d", studentKeyPrefix, id)

	// 1. Check Redis
	var cachedStudent entity.Student
	// Using HGetAll assuming the data was stored as a Hash
	if err := s.rdb.HGetAll(ctx, cacheKey).Scan(&cachedStudent); err == nil && cachedStudent.ID != 0 {
		slog.Info("serving student from cache (hash)", slog.Uint64("id", uint64(id)))
		return &cachedStudent, nil
	}

	// 2. Check DB
	// Cast uint to int64 if your repo expects int64
	student, err := s.repo.GetByID(int64(id))
	if err != nil {
		return nil, err
	}

	// 3. Cache (Async)
	go func(st *entity.Student, key string) {
		ctx := context.Background()
		pipe := s.rdb.Pipeline()
		pipe.HSet(ctx, key, st)
		pipe.Expire(ctx, key, cacheTTL)
		if _, err := pipe.Exec(ctx); err != nil {
			slog.Error("failed to cache student", "error", err)
		}
	}(student, cacheKey)

	slog.Info("student fetched successfully", slog.Uint64("student_id", uint64(id)))
	return student, nil
}

// FindAll - Renamed from GetAllStudents
func (s *studentService) FindAll() ([]entity.Student, error) {
	ctx := context.Background()

	// 1. Check Redis
	val, err := s.rdb.Get(ctx, studentListKey).Result()
	if err == nil {
		var cachedStudents []entity.Student
		if jsonErr := json.Unmarshal([]byte(val), &cachedStudents); jsonErr == nil {
			slog.Info("serving student list from cache")
			return cachedStudents, nil
		}
	}

	// 2. Check DB
	students, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	// 3. Cache (Async)
	go func(st []entity.Student) {
		data, _ := json.Marshal(st)
		s.rdb.Set(context.Background(), studentListKey, data, cacheTTL)
	}(students)

	slog.Info("students fetched successfully", slog.Int("count", len(students)))
	return students, nil
}

// Update - Aligned to accept pointer
func (s *studentService) Update(student *entity.Student) error {
	// Cast ID to int64 for repo
	if err := s.repo.Update(int64(student.ID), *student); err != nil {
		return err
	}

	// Invalidate/Update Cache (Async)
	go func(st entity.Student) {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("%s%d", studentKeyPrefix, st.ID)

		pipe := s.rdb.Pipeline()
		pipe.HSet(ctx, cacheKey, st) // Update individual cache
		pipe.Expire(ctx, cacheKey, cacheTTL)
		pipe.Del(ctx, studentListKey) // Invalidate list cache

		if _, err := pipe.Exec(ctx); err != nil {
			slog.Error("failed to update student cache", "error", err)
		}
	}(*student)

	slog.Info("student updated successfully", slog.Uint64("student_id", uint64(student.ID)))
	return nil
}

// Delete - Aligned to accept uint
func (s *studentService) Delete(id uint) error {
	if err := s.repo.Delete(int64(id)); err != nil {
		return err
	}

	// Clear Cache (Async)
	go func(uid uint) {
		ctx := context.Background()
		pipe := s.rdb.Pipeline()
		pipe.Del(ctx, fmt.Sprintf("%s%d", studentKeyPrefix, uid))
		pipe.Del(ctx, studentListKey)

		if _, err := pipe.Exec(ctx); err != nil {
			slog.Error("failed to update student cache after delete", "error", err)
		}
	}(id)

	slog.Info("student deleted successfully", slog.Uint64("student_id", uint64(id)))
	return nil
}
