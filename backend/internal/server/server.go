package server

import (
	logger2 "University-Selection-Service/pkg/logger"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type ProfileData struct {
	EgeScores          string            `json:"egeScores"`
	Gpa                string            `json:"gpa"`
	Olympiads          string            `json:"olympiads"`
	DesiredSpecialty   string            `json:"desiredSpecialty"`
	EducationType      string            `json:"educationType"`
	Country            string            `json:"country"`
	UniversityLocation string            `json:"universityLocation"`
	Financing          string            `json:"financing"`
	ImportanceFactors  ImportanceFactors `json:"importanceFactors"`
}

type ImportanceFactors struct {
	LocalUniversityRating int  `json:"localUniversityRating"`
	Prestige              int  `json:"prestige"`
	ScholarshipPrograms   int  `json:"scholarshipPrograms"`
	EducationQuality      int  `json:"educationQuality"`
	Dormitory             bool `json:"dormitory"`
	ScientificLabs        bool `json:"scientificLabs"`
	SportsInfrastructure  bool `json:"sportsInfrastructure"`
}

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

		var profile ProfileData
		err := json.NewDecoder(r.Body).Decode(&profile)
		if err != nil {
			log.Error(ctx, "Ошибка при декодировании JSON", zap.Error(err))
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		log.Info(ctx, "Получены данные профиля",
			zap.Any("profile", profile))

		// Здесь можно добавить сохранение в базу данных или другую обработку

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
