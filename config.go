package main

import "os"

const (
	_RATE_LIMIT = 100
	_RATE_TPS   = 2000
	// second minute hour day month week
	// this runs once a day
	_DEFAULT_SCHEDULE = "0 0 0 * * *"
)

func getDBConnectionString() string {
	return os.Getenv("DB_CONNECTION_STRING")
}

func getEmbedderUrl() string {
	return os.Getenv("EMBEDDER_BASE_URL")
}

// func getInternalAuthToken() string {
// 	return os.Getenv("INTERNAL_AUTH_TOKEN")
// }

func getLLMServiceAPIKey() string {
	return os.Getenv("LLMSERVICE_API_KEY")
}

func getInstanceMode() string {
	return os.Getenv("INSTANCE_MODE")
}

func getCollectionSchedule() string {
	schedule := os.Getenv("COLLECTION_SCHEDULE")
	if schedule == "" {
		return _DEFAULT_SCHEDULE
	}
	return schedule
}

func getCleanupSchedule() string {
	schedule := os.Getenv("CLEANUP_SCHEDULE")
	if schedule == "" {
		return _DEFAULT_SCHEDULE
	}
	return schedule
}
