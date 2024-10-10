package dto

type CalculateIslandsRequest struct {
	Matrix [][]bool `json:"matrix"`

	IslandId int `json:"islandId,string"`
}

type GetIslandsResponse struct {
	IslandId int `json:"islandId,string"`
}
