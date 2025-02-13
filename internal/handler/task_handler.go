package handler

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"rashpor.com/todolist/internal/storage"
	"rashpor.com/todolist/internal/storage/models"
)

type TaskHandler struct {
	repo   *storage.TaskRepository
	logger *slog.Logger
}

func NewTaskHandler(repo *storage.TaskRepository, logger *slog.Logger) *TaskHandler {
	return &TaskHandler{repo: repo, logger: logger}
}

// @Summary Создать задачу
// @Description Добавляет новую задачу в список
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Данные новой задачи"
// @Success 201 {object} models.Task
// @Failure 400 {object} models.ErrorResponse "Невалидный JSON"
// @Failure 500 {object} models.ErrorResponse "Ошибка сервера"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	h.logger.Info("Запрос на создание новой заметки")
	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		h.logger.Error("Невалидный JSON", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Невалидный JSON",
		})
	}

	if task.Status == "" {
		task.Status = "new"
	}

	err := h.repo.CreateNewTask(context.Background(), &task)
	if err != nil {
		h.logger.Error("Ошибка создания задачи", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	h.logger.Info("Задача успешно создана")
	return c.Status(fiber.StatusCreated).JSON(task)
}

// @Summary Получить список задач
// @Description Возвращает массив задач, поддерживает пагинацию
// @Tags tasks
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество задач на странице" default(10)
// @Success 200 {array} models.Task
// @Failure 500 {object} models.ErrorResponse "Ошибка сервера"
// @Router /tasks [get]
func (h *TaskHandler) GetAllTasks(c *fiber.Ctx) error {
	h.logger.Info("Запрос на получение всех задач")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	tasks, err := h.repo.GetAllTasks(context.Background(), page, limit)
	if err != nil {
		h.logger.Error("Ошибка получения задач", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при получении задач",
		})
	}
	h.logger.Info("Данные возвращены", slog.Any("page", page), slog.Any("limit", limit))
	return c.Status(fiber.StatusOK).JSON(tasks)
}

// @Summary Обновить задачу
// @Description Обновляет задачу по переданному ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param task body models.Task true "Данные для обновления"
// @Success 200 {object} map[string]string "Сообщение об успешном обновлении"
// @Failure 400 {object} models.ErrorResponse "Некорректный ID или данные"
// @Failure 404 {object} models.ErrorResponse "Задача не найдена"
// @Failure 500 {object} models.ErrorResponse "Ошибка сервера"
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	h.logger.Info("Запрос на обновление заметки " + id)
	taskID, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Error("Некорректный ID", slog.Any("ID", id))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Некорректный ID",
		})
	}

	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		h.logger.Error("Ошибка при парсинге данных", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ошибка при парсинге данных",
		})
	}

	err = h.repo.UpdateTaskOnID(context.Background(), taskID, &task)
	if err != nil {
		if err == models.ErrNoRecord {
			h.logger.Error("Задача с таким ID не найдена", slog.Any("id", id))
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Задача с таким ID не найдена",
			})
		}
		h.logger.Error("Ошибка при обновлении задачи", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при обновлении задачи",
		})
	}

	h.logger.Info("Задача успешно обновлена", slog.Any("id", id))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Задача успешно обновлена",
	})
}

// @Summary Удалить задачу
// @Description Удаляет задачу по переданному ID
// @Tags tasks
// @Param id path int true "ID задачи"
// @Success 200 {object} map[string]string "Сообщение об успешном удалении"
// @Failure 400 {object} models.ErrorResponse "Некорректный ID"
// @Failure 404 {object} models.ErrorResponse "Задача не найдена"
// @Failure 500 {object} models.ErrorResponse "Ошибка сервера"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	h.logger.Info("Запрос на удаление заметки " + id)
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Некорректный ID",
		})
	}

	err = h.repo.DeleteTaskOnID(context.Background(), taskID)
	if err != nil {
		if err == models.ErrNoRecord {
			h.logger.Error("Задача с таким ID не найдена", slog.Any("id", id))
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Задача с таким ID не найдена",
			})
		}
		h.logger.Error("Ошибка при удалении задачи", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при удалении задачи",
		})
	}

	h.logger.Info("Задача успешно удалена", slog.Any("id", id))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Задача успешно удалена",
	})
}
