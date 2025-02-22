package auth

import (
	"ToDo/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ToDo/configs"
	"ToDo/internal/user"
	"ToDo/pkg/res"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService — структура для мокинга IAuthService
type MockAuthService struct {
	mock.Mock
}

// Register — реализация метода Register для мока, соответствует IAuthService
func (m *MockAuthService) Register(ctx context.Context, email, password, name string) (string, error) {
	// Записываем вызов метода и возвращаем заранее заданные значения
	args := m.Called(ctx, email, password, name)
	return args.String(0), args.Error(1) // Возвращаем userID (string) и ошибку
}

// Login — реализация метода Login для мока, соответствует IAuthService
func (m *MockAuthService) Login(ctx context.Context, email, password string) (*models.User, error) {
	// Записываем вызов метода и возвращаем заранее заданные значения
	args := m.Called(ctx, email, password)
	return args.Get(0).(*models.User), args.Error(1) // Возвращаем *user.User и ошибку
}

// TestAuthHandler_Register — тесты для хендлера Register
func TestAuthHandler_Register(t *testing.T) {
	// Таблица тестов с различными сценариями для проверки поведения Register
	tests := []struct {
		name           string                                            // Название теста для вывода в логах
		body           RegisterRequest                                   // Тело запроса, которое будет отправлено в хендлер
		mockRegister   func(m *MockAuthService)                          // Функция для настройки поведения мока
		expectedStatus int                                               // Ожидаемый HTTP-статус ответа
		checkResponse  func(t *testing.T, rr *httptest.ResponseRecorder) // Функция проверки ответа
	}{
		{
			name: "Successful registration", // Сценарий успешной регистрации
			body: RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			mockRegister: func(m *MockAuthService) {
				// Настраиваем мок: при вызове Register с указанными параметрами возвращаем "user123" и nil
				m.On("Register", mock.Anything, "john@example.com", "password123", "John Doe").
					Return("user123", nil)
			},
			expectedStatus: http.StatusOK, // Ожидаем успешный статус 200
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// Проверяем тело ответа: десериализуем в RegisterResponse
				var resp RegisterResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				assert.NoError(t, err, "failed to unmarshal response")
				assert.NotEmpty(t, resp.Token, "token should not be empty")
			},
		},
		{
			name: "User already exists", // Сценарий, когда пользователь уже существует
			body: RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			mockRegister: func(m *MockAuthService) {
				// Настраиваем мок: возвращаем пустую строку и ошибку ErrUserAlreadyExists
				m.On("Register", mock.Anything, "john@example.com", "password123", "John Doe").
					Return("", user.ErrUserAlreadyExists)
			},
			expectedStatus: http.StatusConflict, // Ожидаем 409 Conflict
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// Проверяем тело ответа: десериализуем в ErrorResponse
				var resp res.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				assert.NoError(t, err, "failed to unmarshal response")
				assert.Equal(t, user.ErrUserAlreadyExists.Error(), resp.Error, "unexpected error message")
			},
		},
		{
			name: "Invalid request body", // Сценарий с неверным телом запроса
			body: RegisterRequest{
				Name:     "John Doe",
				Email:    "invalid", // Неверный email, не пройдет валидацию "email"
				Password: "password123",
			},
			mockRegister:   func(m *MockAuthService) {},    // Мок не вызывается, так как валидация провалится
			expectedStatus: http.StatusUnprocessableEntity, // Ожидаем 422 из-за ошибки валидации
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// Здесь тело ответа не проверяем, так как валидация возвращает ошибку в req.HandleBody
			},
		},
	}

	// Проходим по каждому тесту в таблице
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новый мок для каждого теста
			mockService := new(MockAuthService)
			tt.mockRegister(mockService) // Настраиваем поведение мока для Register

			// Создаем тестовую конфигурацию с секретом для JWT
			cfg := &configs.Config{
				Auth: struct {
					Secret        string        `mapstructure:"SECRET"`
					TokenLifetime time.Duration `mapstructure:"TOKEN_LIFETIME"`
				}{
					Secret: "test-secret",
				},
			}

			// Инициализируем хендлер с конфигурацией и мок-сервисом
			handler := &AuthHandler{
				Config:      cfg,
				AuthService: mockService,
			}

			// Сериализуем тело запроса в JSON
			bodyBytes, _ := json.Marshal(tt.body)
			// Создаем фиктивный HTTP-запрос
			req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			// Создаем рекордер для записи ответа от хендлера
			rr := httptest.NewRecorder()

			// Выполняем хендлер Register с нашим запросом
			handler.Register()(rr, req)

			// Проверяем, что статус ответа соответствует ожидаемому
			assert.Equal(t, tt.expectedStatus, rr.Code, "unexpected status code")

			// Выполняем пользовательскую проверку ответа
			tt.checkResponse(t, rr)

			// Проверяем, что все ожидаемые вызовы мока были выполнены
			mockService.AssertExpectations(t)
		})
	}
}

// TestAuthHandler_Login — тесты для хендлера Login
func TestAuthHandler_Login(t *testing.T) {
	// Таблица тестов с различными сценариями для проверки поведения Login
	tests := []struct {
		name           string                                            // Название теста для вывода в логах
		body           LoginRequest                                      // Тело запроса, которое будет отправлено в хендлер
		mockLogin      func(m *MockAuthService)                          // Функция для настройки поведения мока
		expectedStatus int                                               // Ожидаемый HTTP-статус ответа
		checkResponse  func(t *testing.T, rr *httptest.ResponseRecorder) // Функция проверки ответа
	}{
		{
			name: "Successful login", // Сценарий успешного входа
			body: LoginRequest{
				Email:    "john@example.com",
				Password: "password123",
			},
			mockLogin: func(m *MockAuthService) {
				// Настраиваем мок: при вызове Login возвращаем объект User и nil
				m.On("Login", mock.Anything, "john@example.com", "password123").
					Return(&models.User{
						ID:       "user123",
						Email:    "john@example.com",
						Password: "hashed_password",
						Name:     "John Doe",
					}, nil)
			},
			expectedStatus: http.StatusOK, // Ожидаем успешный статус 200
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// Проверяем тело ответа: десериализуем в LoginResponse
				var resp LoginResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				assert.NoError(t, err, "failed to unmarshal response")
				assert.NotEmpty(t, resp.Token, "token should not be empty")
			},
		},
		{
			name: "Invalid credentials", // Сценарий с неверными учетными данными
			body: LoginRequest{
				Email:    "john@example.com",
				Password: "wrongpassword",
			},
			mockLogin: func(m *MockAuthService) {
				// Настраиваем мок: возвращаем nil и ErrUserNotFound
				m.On("Login", mock.Anything, "john@example.com", "wrongpassword").
					Return((*models.User)(nil), user.ErrUserNotFound)
			},
			expectedStatus: http.StatusUnauthorized, // Ожидаем 401 Unauthorized
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// Проверяем тело ответа: десериализуем в ErrorResponse
				var resp res.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				assert.NoError(t, err, "failed to unmarshal response")
				assert.Equal(t, "invalid credentials", resp.Error, "unexpected error message")
			},
		},
		{
			name: "Invalid request body", // Сценарий с неверным телом запроса
			body: LoginRequest{
				Email: "john@example.com", // Пустой пароль не пройдет валидацию
			},
			mockLogin:      func(m *MockAuthService) {},    // Мок не вызывается из-за ошибки валидации
			expectedStatus: http.StatusUnprocessableEntity, // Ожидаем 422 из-за ошибки валидации
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// Здесь тело ответа не проверяем, так как валидация возвращает ошибку в req.HandleBody
			},
		},
	}

	// Проходим по каждому тесту в таблице
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новый мок для каждого теста
			mockService := new(MockAuthService)
			tt.mockLogin(mockService) // Настраиваем поведение мока для Login

			// Создаем тестовую конфигурацию с секретом для JWT
			cfg := &configs.Config{
				Auth: struct {
					Secret        string        `mapstructure:"SECRET"`
					TokenLifetime time.Duration `mapstructure:"TOKEN_LIFETIME"`
				}{
					Secret: "test-secret",
				},
			}

			// Инициализируем хендлер с конфигурацией и мок-сервисом
			handler := &AuthHandler{
				Config:      cfg,
				AuthService: mockService,
			}

			// Сериализуем тело запроса в JSON
			bodyBytes, _ := json.Marshal(tt.body)
			// Создаем фиктивный HTTP-запрос
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			// Создаем рекордер для записи ответа от хендлера
			rr := httptest.NewRecorder()

			// Выполняем хендлер Login с нашим запросом
			handler.Login()(rr, req)

			// Проверяем, что статус ответа соответствует ожидаемому
			assert.Equal(t, tt.expectedStatus, rr.Code, "unexpected status code")

			// Выполняем пользовательскую проверку ответа
			tt.checkResponse(t, rr)

			// Проверяем, что все ожидаемые вызовы мока были выполнены
			mockService.AssertExpectations(t)
		})
	}
}
