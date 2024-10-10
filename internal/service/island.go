package service

import (
	"errors"
	"golang-island/internal/data"
	"golang-island/internal/dto"
	"golang-island/internal/storage"
	"golang.org/x/exp/slog"
	"strconv"
	"time"
)

type IslandService interface {
	PutTask(request dto.CalculateIslandsRequest) error

	GetResult(islandId int) (int, error)

	IsReady(islandId int) (bool, error)
}

type IslandServiceImpl struct {
	repo *storage.MockRepository

	log *slog.Logger
}

func NewService(log *slog.Logger) IslandService {
	return &IslandServiceImpl{repo: storage.NewRepo(log), log: log}
}

func (service *IslandServiceImpl) PutTask(request dto.CalculateIslandsRequest) error {
	matrix := request.Matrix
	id := request.IslandId

	if matrix == nil || len(matrix) == 0 {
		return errors.New("matrix is empty")
	}

	if id <= 0 {
		return errors.New("id must be positive")
	}

	if service.repo.ExistsById(id) {
		return errors.New("already exists by id " + strconv.Itoa(id))
	}

	go service.calculate(matrix, id)

	return nil
}

func (service *IslandServiceImpl) calculate(matrix [][]bool, id int) {
	service.log.Info("start calculation for id=" + strconv.Itoa(id))

	service.repo.PutIfEmpty(
		id,
		data.Data{IslandCount: -1},
	)

	rows := len(matrix)
	cols := len(matrix[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	count := 0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if matrix[i][j] && !visited[i][j] {
				dfs(matrix, visited, i, j, rows, cols)
				count++
			}
		}
	}

	service.repo.Put(
		id,
		data.Data{IslandCount: count, CalculationDate: time.Now()},
	)

	service.log.Info("end calculation for id=" + strconv.Itoa(id) + ", result=" + strconv.Itoa(count))
}

func (service *IslandServiceImpl) GetResult(id int) (int, error) {

	result, err := service.repo.GetById(id)

	if err != nil {
		return 0, err
	}

	return result.IslandCount, err
}

func (service *IslandServiceImpl) IsReady(id int) (bool, error) {

	result, err := service.repo.GetById(id)

	if err != nil {
		return false, err
	} else if result.IslandCount == -1 {
		return false, nil
	}

	return true, nil
}

func dfs(matrix, visited [][]bool, i, j, rows, cols int) {
	if i < 0 || i >= rows || j < 0 || j >= cols || visited[i][j] || !matrix[i][j] {
		return
	}

	visited[i][j] = true

	dfs(matrix, visited, i+1, j, rows, cols)
	dfs(matrix, visited, i-1, j, rows, cols)
	dfs(matrix, visited, i, j+1, rows, cols)
	dfs(matrix, visited, i, j-1, rows, cols)
}
