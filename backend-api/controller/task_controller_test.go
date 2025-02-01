package controller

import (
	"bytes"
	"encoding/json"
	"go-rest-api/model"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// mockTaskUsecase はタスクユースケースのモック
type mockTaskUsecase struct {
	mockGetAllTasks func(userId uint) ([]model.TaskResponse, error)
	mockGetTaskById func(userId, taskId uint) (model.TaskResponse, error)
	mockCreateTask  func(task model.Task) (model.TaskResponse, error)
	mockUpdateTask  func(task model.Task, userId, taskId uint) (model.TaskResponse, error)
	mockDeleteTask  func(userId, taskId uint) error
}

func (m *mockTaskUsecase) GetAllTasks(userId uint) ([]model.TaskResponse, error) {
	return m.mockGetAllTasks(userId)
}

func (m *mockTaskUsecase) GetTaskById(userId, taskId uint) (model.TaskResponse, error) {
	return m.mockGetTaskById(userId, taskId)
}

func (m *mockTaskUsecase) CreateTask(task model.Task) (model.TaskResponse, error) {
	return m.mockCreateTask(task)
}

func (m *mockTaskUsecase) UpdateTask(task model.Task, userId, taskId uint) (model.TaskResponse, error) {
	return m.mockUpdateTask(task, userId, taskId)
}

func (m *mockTaskUsecase) DeleteTask(userId, taskId uint) error {
	return m.mockDeleteTask(userId, taskId)
}

// テストヘルパー関数
func setupTaskTest(t *testing.T) (*echo.Echo, *httptest.ResponseRecorder) {
	e := echo.New()
	rec := httptest.NewRecorder()
	return e, rec
}

func createJWTToken(userId uint) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userId),
	})
	return token
}

func TestTaskController_GetAllTasks(t *testing.T) {
	e, rec := setupTaskTest(t)

	tests := []struct {
		name            string
		userId          uint
		mockBehavior    func(*mockTaskUsecase)
		expectedCode    int
		expectedTasks   []model.TaskResponse
		expectError     bool
		expectedMessage string
	}{
		{
			name:   "正常系：タスク一覧取得成功",
			userId: 1,
			mockBehavior: func(m *mockTaskUsecase) {
				m.mockGetAllTasks = func(userId uint) ([]model.TaskResponse, error) {
					return []model.TaskResponse{
						{
							ID:        1,
							Title:     "Test Task 1",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
					}, nil
				}
			},
			expectedCode:  http.StatusOK,
			expectedTasks: []model.TaskResponse{{ID: 1, Title: "Test Task 1"}},
			expectError:   false,
		},
		{
			name:   "異常系：タスク取得エラー",
			userId: 1,
			mockBehavior: func(m *mockTaskUsecase) {
				m.mockGetAllTasks = func(userId uint) ([]model.TaskResponse, error) {
					return nil, echo.NewHTTPError(http.StatusInternalServerError, "タスクの取得に失敗しました")
				}
			},
			expectedCode:    http.StatusInternalServerError,
			expectError:     true,
			expectedMessage: "タスクの取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mock := &mockTaskUsecase{}
			tt.mockBehavior(mock)

			// コントローラーの作成
			tc := NewTaskController(mock)

			// リクエストの準備
			req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			c := e.NewContext(req, rec)
			c.Set("user", createJWTToken(tt.userId))

			// テスト実行
			err := tc.GetAllTasks(c)

			// レスポンスの検証
			if tt.expectError {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, he.Code)
				assert.Equal(t, tt.expectedMessage, he.Message)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, rec.Code)

				var response []model.TaskResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTasks[0].ID, response[0].ID)
				assert.Equal(t, tt.expectedTasks[0].Title, response[0].Title)
			}

			// レコーダーをリセット
			rec.Body.Reset()
		})
	}
}

func TestTaskController_GetTaskById(t *testing.T) {
	e, rec := setupTaskTest(t)

	tests := []struct {
		name            string
		userId          uint
		taskId          uint
		mockBehavior    func(*mockTaskUsecase)
		expectedCode    int
		expectedTask    model.TaskResponse
		expectError     bool
		expectedMessage string
	}{
		{
			name:   "正常系：タスク取得成功",
			userId: 1,
			taskId: 1,
			mockBehavior: func(m *mockTaskUsecase) {
				m.mockGetTaskById = func(userId, taskId uint) (model.TaskResponse, error) {
					return model.TaskResponse{
						ID:        1,
						Title:     "Test Task",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil
				}
			},
			expectedCode: http.StatusOK,
			expectedTask: model.TaskResponse{ID: 1, Title: "Test Task"},
			expectError:  false,
		},
		{
			name:   "異常系：タスク取得エラー",
			userId: 1,
			taskId: 999,
			mockBehavior: func(m *mockTaskUsecase) {
				m.mockGetTaskById = func(userId, taskId uint) (model.TaskResponse, error) {
					return model.TaskResponse{}, echo.NewHTTPError(http.StatusInternalServerError, "タスクの取得に失敗しました")
				}
			},
			expectedCode:    http.StatusInternalServerError,
			expectError:     true,
			expectedMessage: "タスクの取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mock := &mockTaskUsecase{}
			tt.mockBehavior(mock)

			// コントローラーの作成
			tc := NewTaskController(mock)

			// リクエストの準備
			req := httptest.NewRequest(http.MethodGet, "/tasks/:taskId", nil)
			c := e.NewContext(req, rec)
			c.SetParamNames("taskId")
			c.SetParamValues(strconv.FormatUint(uint64(tt.taskId), 10))
			c.Set("user", createJWTToken(tt.userId))

			// テスト実行
			err := tc.GetTaskById(c)

			// レスポンスの検証
			if tt.expectError {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, he.Code)
				assert.Equal(t, tt.expectedMessage, he.Message)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, rec.Code)

				var response model.TaskResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTask.ID, response.ID)
				assert.Equal(t, tt.expectedTask.Title, response.Title)
			}

			// レコーダーをリセット
			rec.Body.Reset()
		})
	}
}

func TestTaskController_CreateTask(t *testing.T) {
	e, rec := setupTaskTest(t)

	tests := []struct {
		name            string
		userId          uint
		inputTask       model.Task
		mockBehavior    func(*mockTaskUsecase)
		expectedCode    int
		expectedTask    model.TaskResponse
		expectError     bool
		expectedMessage string
	}{
		{
			name:   "正常系：タスク作成成功",
			userId: 1,
			inputTask: model.Task{
				Title: "New Task",
			},
			mockBehavior: func(m *mockTaskUsecase) {
				m.mockCreateTask = func(task model.Task) (model.TaskResponse, error) {
					return model.TaskResponse{
						ID:        1,
						Title:     task.Title,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil
				}
			},
			expectedCode: http.StatusCreated,
			expectedTask: model.TaskResponse{ID: 1, Title: "New Task"},
			expectError:  false,
		},
		{
			name:   "異常系：タスク作成エラー",
			userId: 1,
			inputTask: model.Task{
				Title: "",
			},
			mockBehavior: func(m *mockTaskUsecase) {
				m.mockCreateTask = func(task model.Task) (model.TaskResponse, error) {
					return model.TaskResponse{}, echo.NewHTTPError(http.StatusBadRequest, "タスクの作成に失敗しました")
				}
			},
			expectedCode:    http.StatusBadRequest,
			expectError:     true,
			expectedMessage: "タスクの作成に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mock := &mockTaskUsecase{}
			tt.mockBehavior(mock)

			// コントローラーの作成
			tc := NewTaskController(mock)

			// リクエストの準備
			jsonBody, _ := json.Marshal(tt.inputTask)
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.Set("user", createJWTToken(tt.userId))

			// テスト実行
			err := tc.CreateTask(c)

			// レスポンスの検証
			if tt.expectError {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, he.Code)
				assert.Equal(t, tt.expectedMessage, he.Message)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, rec.Code)

				var response model.TaskResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTask.ID, response.ID)
				assert.Equal(t, tt.expectedTask.Title, response.Title)
			}

			// レコーダーをリセット
			rec.Body.Reset()
		})
	}
}

func TestTaskController_UpdateTask(t *testing.T) {
	e, rec := setupTaskTest(t)

	tests := []struct {
		name            string
		userId          uint
		taskId          uint
		inputTask       model.Task
		mockBehavior    func(*mockTaskUsecase)
		expectedCode    int
		expectedTask    model.TaskResponse
		expectError     bool
		expectedMessage string
	}{
		{
			name:   "正常系：タスク更新成功",
			userId: 1,
			taskId: 1,
			inputTask: model.Task{
				Title: "Updated Task",
			},
			mockBehavior: func(m *mockTaskUsecase) {
				m.mockUpdateTask = func(task model.Task, userId, taskId uint) (model.TaskResponse, error) {
					return model.TaskResponse{
						ID:        taskId,
						Title:     task.Title,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil
				}
			},
			expectedCode: http.StatusOK,
			expectedTask: model.TaskResponse{ID: 1, Title: "Updated Task"},
			expectError:  false,
		},
		{
			name:   "異常系：タスク更新エラー",
			userId: 1,
			taskId: 999,
			inputTask: model.Task{
				Title: "Updated Task",
			},
			mockBehavior: func(m *mockTaskUsecase) {
				m.mockUpdateTask = func(task model.Task, userId, taskId uint) (model.TaskResponse, error) {
					return model.TaskResponse{}, echo.NewHTTPError(http.StatusInternalServerError, "タスクの更新に失敗しました")
				}
			},
			expectedCode:    http.StatusInternalServerError,
			expectError:     true,
			expectedMessage: "タスクの更新に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
