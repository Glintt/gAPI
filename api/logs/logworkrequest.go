package logs

import "github.com/Glintt/gAPI/api/logs/models"

type LogWorkRequest struct {
	Name      string
	LogToSave models.RequestLogging
}
