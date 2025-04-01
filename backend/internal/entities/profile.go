package entities

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
