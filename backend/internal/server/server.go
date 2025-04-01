package server

import (
	"University-Selection-Service/internal/entities"
	logger2 "University-Selection-Service/pkg/logger"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func createProfileHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger2.GetLoggerFromCtx(ctx)

		// Настройка CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Обработка preflight-запросов
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var profile entities.ProfileData
		err := json.NewDecoder(r.Body).Decode(&profile)
		if err != nil {
			log.Error(ctx, "Ошибка при декодировании JSON", zap.Error(err))
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		log.Info(ctx, "Получены данные профиля",
			zap.Any("profile", profile))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Профиль успешно создан",
		})
	}
}

func StartServer(ctx context.Context) error {
	logger := logger2.GetLoggerFromCtx(ctx)
	mux := http.NewServeMux()
	mux.HandleFunc("/profile", createProfileHandler(ctx))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		logger.Error(ctx, "error starting server", zap.Error(err))
		return err
	}
	return nil
}
