package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"study-tracker-api/models"
	"sync"
)

var (
	lessons     []models.Lesson
	assignments []models.Assignment
	mutex       sync.RWMutex
)

const (
	SCHEDULE_FILE    = "data/schedule.json"
	ASSIGNMENTS_FILE = "data/assignments.json"
)

func InitStorage() {
	os.MkdirAll("data", 0755)
	loadLessons()
	loadAssignments()
}

func loadLessons() {
	data, err := os.ReadFile(SCHEDULE_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			lessons = []models.Lesson{}
			return
		}
		panic(fmt.Sprintf("Ошибка загрузки занятий: %v", err))
	}
	err = json.Unmarshal(data, &lessons)
	if err != nil {
		panic(fmt.Sprintf("Ошибка парсинга занятий: %v", err))
	}
}

func saveLessons() error {
	data, err := json.MarshalIndent(lessons, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %w", err)
	}
	err = os.WriteFile(SCHEDULE_FILE, data, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}
	fmt.Println("✅ Файл schedule.json успешно сохранён")
	return nil
}

func loadAssignments() {
	data, err := os.ReadFile(ASSIGNMENTS_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			assignments = []models.Assignment{}
			return
		}
		panic(fmt.Sprintf("Ошибка загрузки заданий: %v", err))
	}
	err = json.Unmarshal(data, &assignments)
	if err != nil {
		panic(fmt.Sprintf("Ошибка парсинга заданий: %v", err))
	}
}

func saveAssignments() error {
	data, err := json.MarshalIndent(assignments, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ASSIGNMENTS_FILE, data, 0644)
}

// --- Утилиты для работы с ID ---
func generateID(sliceLen int) string {
	return fmt.Sprintf("%d", sliceLen+1)
}

// --- Работа с занятиями ---
func GetLessons() []models.Lesson {
	mutex.RLock()
	defer mutex.RUnlock()
	return lessons
}

func GetLessonByID(id string) (*models.Lesson, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	for _, l := range lessons {
		if l.ID == id {
			return &l, true
		}
	}
	return nil, false
}

func CreateLesson(lesson models.Lesson) (models.Lesson, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// Проверка на пересечение
	for _, existing := range lessons {
		if existing.Date == lesson.Date && existing.OverlapsWith(&lesson) {
			return models.Lesson{}, fmt.Errorf("время занято")
		}
	}

	lesson.ID = generateID(len(lessons))
	lessons = append(lessons, lesson)
	err := saveLessons()
	return lesson, err
}

func UpdateLesson(id string, updated models.Lesson) (*models.Lesson, error) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, l := range lessons {
		if l.ID == id {
			// Проверка пересечения (кроме самого себя)
			for j, other := range lessons {
				if j != i && other.Date == updated.Date && other.OverlapsWith(&updated) {
					return nil, fmt.Errorf("время занято")
				}
			}
			lessons[i] = updated
			lessons[i].ID = id // сохраняем ID
			err := saveLessons()
			return &lessons[i], err
		}
	}
	return nil, fmt.Errorf("занятие не найдено")
}

func DeleteLesson(id string) error {
	mutex.Lock()
	defer mutex.Unlock()
	for i, l := range lessons {
		if l.ID == id {
			lessons = append(lessons[:i], lessons[i+1:]...)
			return saveLessons()
		}
	}
	return fmt.Errorf("занятие не найдено")
}

// --- Работа с заданиями ---
func GetAssignments() []models.Assignment {
	mutex.RLock()
	defer mutex.RUnlock()
	return assignments
}

func GetAssignmentByID(id string) (*models.Assignment, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	for _, a := range assignments {
		if a.ID == id {
			return &a, true
		}
	}
	return nil, false
}

func CreateAssignment(assignment models.Assignment) (models.Assignment, error) {
	mutex.Lock()
	defer mutex.Unlock()
	assignment.ID = generateID(len(assignments))
	assignments = append(assignments, assignment)
	err := saveAssignments()
	return assignment, err
}

func UpdateAssignment(id string, updated models.Assignment) (*models.Assignment, error) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, a := range assignments {
		if a.ID == id {
			assignments[i] = updated
			assignments[i].ID = id
			err := saveAssignments()
			return &assignments[i], err
		}
	}
	return nil, fmt.Errorf("задание не найдено")
}

func DeleteAssignment(id string) error {
	mutex.Lock()
	defer mutex.Unlock()
	for i, a := range assignments {
		if a.ID == id {
			assignments = append(assignments[:i], assignments[i+1:]...)
			return saveAssignments()
		}
	}
	return fmt.Errorf("задание не найдено")
}

// --- Получение нагрузки ---
func GetWorkload() map[string][]string {
	mutex.RLock()
	defer mutex.RUnlock()

	workload := make(map[string][]string)

	for _, l := range lessons {
		event := fmt.Sprintf("📚 %s (%s–%s)", l.Subject, l.StartTime, l.EndTime)
		workload[l.Date] = append(workload[l.Date], event)
	}

	for _, a := range assignments {
		status := "⏳"
		if a.Completed {
			status = "✅"
		}
		event := fmt.Sprintf("%s %s (%s)", status, a.Title, a.Subject)
		workload[a.DueDate] = append(workload[a.DueDate], event)
	}

	return workload
}
