package logs

import "gAPIManagement/api/logs/models"

type LogWorkRequest struct {
	Name      string
	LogToSave models.RequestLogging
}
