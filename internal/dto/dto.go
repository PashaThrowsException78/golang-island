package dto

type CalculateIslandsRequest struct {
	Matrix [][]bool `json:"matrix"`

	IslandId int `json:"islandId"`
}

type GetIslandsResponse struct {
	IslandId int `json:"islandId"`
}
