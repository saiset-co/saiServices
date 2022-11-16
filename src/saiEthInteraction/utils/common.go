package utils

import "github.com/iamthe1whoknocks/saiEthInteraction/models"

func RemoveContract(slice []models.Contract, s int) []models.Contract {
	return append(slice[:s], slice[s+1:]...)
}
