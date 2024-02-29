package domain

import "movie-management-service/utils"

// DBStatus returns DB connection status.
func (ms *MovieService) DBStatus() (bool, error) {
	if err := ms.db.Ping(); err != nil {
		utils.ErrorLogger.Println(err)
		return false, err
	}

	return true, nil
}
