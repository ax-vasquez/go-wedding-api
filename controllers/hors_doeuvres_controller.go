package controllers

import (
	"log"
	"net/http"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HorsDoeuvresData struct {
	HorsDoeuvres []models.HorsDoeuvres `json:"hors_doeuvres"`
}

type V1_API_RESPONSE_HORS_DOEVRES struct {
	V1_API_RESPONSE
	Data HorsDoeuvresData `json:"data"`
}

// Get a list of hors doeuvres
//
// If a user ID is specified, then this will return the list of HorsDoeuvres containing
// one item (data for the hors doeuvres they have selected), or zero items if no selection
// has been made, yet. If no user ID is specified, then data for all possible HorsDoeuvres
// is returned.
func GetHorsDoeuvres(c *gin.Context) {
	idStr := c.Param("id")
	var response V1_API_RESPONSE_HORS_DOEVRES
	var status int
	var horsDoeuvres []models.HorsDoeuvres
	id, err := uuid.Parse(idStr)
	if err == nil {
		status = http.StatusOK
		horsDoeuvres, err = models.FindHorsDoeuvresForUser(id)
		if err != nil {
			status = http.StatusInternalServerError
			log.Println(err.Error())
			response.Message = "Internal server error"
		} else {
			status = http.StatusOK
		}
	} else {
		var err error
		horsDoeuvres, err = models.FindHorsDoeuvres()
		if err != nil {
			status = http.StatusInternalServerError
			log.Println(err.Error())
			response.Message = "Internal server error"
		} else {
			status = http.StatusOK
		}
	}
	response.Status = status
	response.Data = HorsDoeuvresData{
		HorsDoeuvres: horsDoeuvres,
	}
	c.JSON(status, response)
}

// Controller to handle creating a new hors doeuvres
//
// The route that uses this controller must be protected so that
// only site admins can use this endpoint. All other requests
// should be rejected.
func CreateHorsDoeuvres(c *gin.Context) {
	// TODO: Add logic to reject unauthorized requests (and certainly do not deploy until all auth logic is wired up)
	response := V1_API_RESPONSE_HORS_DOEVRES{}
	var status int
	var input models.HorsDoeuvres
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = "\"option_name\" is required"
	} else {
		horsDoeuvres := []models.HorsDoeuvres{input}
		err := models.CreateHorsDoeuvres(&horsDoeuvres)
		if err != nil {
			status = http.StatusInternalServerError
			response.Message = "Internal server error"
			log.Println("Error inserting hors_doeuvres record: ", err.Error())
		} else {
			status = http.StatusCreated
			response.Message = "Created new hors doeuvres"
			response.Data.HorsDoeuvres = horsDoeuvres
		}
	}
	response.Status = status
	c.JSON(status, response)
}

// Delete an hors doeuvres
func DeleteHorsDoeuvres(c *gin.Context) {
	response := V1_API_DELETE_RESPONSE{}
	var status int
	id, _ := uuid.Parse(c.Param("id"))
	result, err := models.DeleteHorsDoeuvres(id)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Internal server error"
		log.Println("Error deleting hors doeuvres: ", err.Error())
	} else {
		status = http.StatusAccepted
		response.Message = "Deleted hors doeuvres"
		response.Data.DeletedRecords = int(*result)
	}
	response.Status = status
	c.JSON(status, response)
}
