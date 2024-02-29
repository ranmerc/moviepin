package domain

import (
	"database/sql"
	"errors"
	"sync"

	"movie-management-service/model"
	"movie-management-service/utils"
)

var (
	// Error failed to get movies.
	ErrFailedGetMovies = errors.New("failed to get movies")

	// Error failed to get movie.
	ErrFailedGetMovie = errors.New("failed to get movie")

	// Error failed to add movie.
	ErrFailedAdd = errors.New("failed to add movie")

	// Error failed to delete movie.
	ErrFailedDelete = errors.New("failed to delete movie")

	// Error failed to update movie.
	ErrFailedUpdate = errors.New("failed to update movie")

	// Error failed to get movie rating.
	ErrFailedGetRating = errors.New("failed to get movie rating")

	// Error failed to replace movies.
	ErrFailedReplace = errors.New("failed to replace movies")

	// Error returned when movie does not exist.
	ErrNotExists = errors.New("movie does not exist")
)

// Returns slice of all movies present.
func (ms MovieService) GetMovies() ([]*model.Movie, error) {
	sqlStatement := `SELECT 
										"movieID", 
										"title", 
										"releaseDate", 
										"genre", 
										"director", 
										"description" 
									FROM 
										"movies";`

	rows, err := ms.db.Query(sqlStatement)

	if err != nil {
		utils.ErrorLogger.Println(err)
		return nil, ErrFailedGetMovies
	}

	defer rows.Close()

	movies := make([]*model.Movie, 0)

	for rows.Next() {
		movie := &model.Movie{}

		if err := rows.Scan(&movie.ID, &movie.Title, &movie.ReleaseDate, &movie.Genre, &movie.Director, &movie.Description); err != nil {
			return nil, ErrFailedGetMovies
		}

		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrFailedGetMovies
	}

	return movies, nil
}

// Replaces movies collection with passed in collection
func (ms MovieService) ReplaceMovies(movies []model.Movie) error {
	tx, err := ms.db.Begin()

	if err != nil {
		return ErrFailedReplace
	}

	defer tx.Rollback()

	// Reviews are deleted first because of foreign key constraint.
	sqlStatement := `DELETE FROM "reviews";`

	_, err = tx.Exec(sqlStatement)

	if err != nil {
		return ErrFailedReplace
	}

	sqlStatement = `DELETE FROM "movies";`

	_, err = tx.Exec(sqlStatement)

	if err != nil {
		return ErrFailedReplace
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, movie := range movies {
		wg.Add(1)

		go func(movie model.Movie) {
			defer wg.Done()

			sqlStatement = `INSERT INTO "movies" (
												"movieID", "title", "releaseDate", 
												"genre", "director", "description"
											) 
											VALUES 
												($1, $2, $3, $4, $5, $6);`

			_, err := tx.Exec(sqlStatement, movie.ID, movie.Title, movie.ReleaseDate, movie.Genre, movie.Director, movie.Description)

			if err != nil {
				mu.Lock()
				tx.Rollback()
				mu.Unlock()

				return
			}
		}(movie)
	}

	wg.Wait()

	if err = tx.Commit(); err != nil {
		return ErrFailedReplace
	}

	return nil
}

// Returns particular movie.
func (ms MovieService) GetMovie(id string) (*model.Movie, error) {
	sqlStatement := `SELECT 
										"movieID", 
										"title", 
										"releaseDate", 
										"genre", 
										"director", 
										"description" 
									FROM 
										"movies" 
									WHERE 
										"movieID" = $1;`

	row := ms.db.QueryRow(sqlStatement, id)

	movie := &model.Movie{}

	if err := row.Scan(&movie.ID, &movie.Title, &movie.ReleaseDate, &movie.Genre, &movie.Director, &movie.Description); err != nil {

		if err == sql.ErrNoRows {
			return nil, ErrNotExists
		}

		utils.ErrorLogger.Println(err)
		return nil, ErrFailedGetMovie
	}

	return movie, nil
}

// Adds movie to the database.
func (ms MovieService) AddMovie(newMovie model.Movie) error {
	sqlStatement := `INSERT INTO "movies" (
											"movieID", "title", "releaseDate", 
											"genre", "director", "description"
										) 
										VALUES 
											($1, $2, $3, $4, $5, $6);`

	if _, err := ms.db.Exec(sqlStatement, newMovie.ID, newMovie.Title, newMovie.ReleaseDate, newMovie.Genre, newMovie.Director, newMovie.Description); err != nil {
		utils.ErrorLogger.Println(err)
		return ErrFailedAdd
	}

	return nil
}

// Deletes a movie from the database.
func (ms MovieService) DeleteMovie(id string) error {
	sqlStatement := `DELETE FROM 
											"movies" 
										WHERE 
											"movieID" = $1;`

	result, err := ms.db.Exec(sqlStatement, id)

	if err != nil {
		utils.ErrorLogger.Println(err)
		return ErrFailedDelete
	}

	num, err := result.RowsAffected()

	if err != nil {
		utils.ErrorLogger.Println(err)
		return ErrFailedDelete
	}

	if num == 0 {
		return ErrNotExists
	}

	return nil
}

// Updates a movie in the database.
func (ms MovieService) UpdateMovie(id string, movie model.Movie) error {
	sqlStatement := `UPDATE 
											"movies" 
										SET 
											"movieID" = $1, 
											"title" = $2, 
											"releaseDate" = $3, 
											"genre" = $4, 
											"director" = $5, 
											"description" = $6 
										WHERE 
											"movieID" = $7;`

	result, err := ms.db.Exec(sqlStatement, movie.ID, movie.Title, movie.ReleaseDate, movie.Genre, movie.Director, movie.Description, id)

	if err != nil {
		utils.ErrorLogger.Println(err)
		return ErrFailedUpdate
	}

	num, err := result.RowsAffected()

	if err != nil {
		utils.ErrorLogger.Println(err)
		return ErrFailedUpdate
	}

	if num == 0 {
		return ErrNotExists
	}

	return nil
}

// Returns movie details along with its rating.
func (ms MovieService) GetMovieRating(id string) (*model.MovieReview, error) {
	sqlStatement := `SELECT 
										m."movieID", 
										m."title", 
										m."releaseDate", 
										m."genre", 
										m."director", 
										m."description", 
										COALESCE(
											TRUNC(
												ROUND(
													AVG(r.rating)
												), 
												1
											),
											0
										) AS "rating"
									FROM 
										"movies" m 
										LEFT JOIN "reviews" r ON m."movieID" = r."movieID" 
									WHERE 
										m."movieID" = $1 
									GROUP BY 
										m."movieID";`

	row := ms.db.QueryRow(sqlStatement, id)

	mr := &model.MovieReview{}

	if err := row.Scan(&mr.ID, &mr.Title, &mr.ReleaseDate, &mr.Genre, &mr.Director, &mr.Description, &mr.Rating); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotExists
		}

		utils.ErrorLogger.Println(err)
		return nil, ErrFailedGetRating
	}

	return mr, nil
}
