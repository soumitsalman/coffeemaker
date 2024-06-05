package main

import "os"

// defaults
const (
	_RATE_LIMIT = 100
	_RATE_TPS   = 2000
	// second minute hour day month week
	// this runs once a day
	_ONCE_A_DAY = "0 0 0 * * *"
	_PORT       = ":8080"
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

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = _PORT
	} else {
		port = ":" + port
	}
	return port
}

func getCollectionSchedule() string {
	schedule := os.Getenv("COLLECTION_SCHEDULE")
	if schedule == "" {
		return _ONCE_A_DAY
	}
	return schedule
}

func getRectifySchedule() string {
	schedule := os.Getenv("RECTIFY_SCHEDULE")
	if schedule == "" {
		return _ONCE_A_DAY
	}
	return schedule
}

func getCleanupSchedule() string {
	schedule := os.Getenv("CLEANUP_SCHEDULE")
	if schedule == "" {
		return _ONCE_A_DAY
	}
	return schedule
}
