package controllers

import (
	"log"
	"net/http"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EntreeData struct {
	Entrees []models.Entree `json:"entrees"`
}

type V1_API_RESPONSE_ENTREE struct {
	V1_API_RESPONSE
	Data EntreeData `json:"data"`
}

// GetEntrees gets one or all entrees
//
//	@Summary      	gets one or all entrees
//	@Description  	gets 1 entree if an ID is found in the route, otherwise returns all entrees
//	@Tags         	entrees
//	@Accept       	json
//	@Produce      	json
//	@Success      	200  {object}  V1_API_RESPONSE_ENTREE
//	@Failure      	500  {object}  V1_API_RESPONSE_ENTREE
//	@Param 			entree_id  path uuid.UUID true "Entree ID"
//	@Router       	/entree [get]
//	@Router       	/user/{entree_id}/entrees [get]
func GetEntrees(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	response := V1_API_RESPONSE_ENTREE{}
	var status int
	var entrees []models.Entree
	// If no error occurred, the parse was successful, meaning a UUID was found and results will be filtered for the given user
	if err == nil {
		status = http.StatusOK
		entrees, err = models.FindEntreesForUser(id)
		if err != nil {
			status = http.StatusInternalServerError
			log.Println(err.Error())
			response.Message = "Internal server error"
		}
		// If an error occurred, we ignore it and assume it's because there was no ID in the path - all results will be returned
	} else {
		entrees, err = models.FindEntrees()
		if err != nil {
			status = http.StatusInternalServerError
			log.Println(err.Error())
			response.Message = "Internal server error"
		} else {
			status = http.StatusOK
		}
	}
	response.Status = status
	response.Data.Entrees = entrees
	c.JSON(status, response)
}

// CreateEntree creates an entree
//
//	@Summary      Create entree
//	@Description  create a new entree
//	@Tags         entrees
//	@Accept       json
//	@Produce      json
//	@Success      202  {object}  V1_API_RESPONSE_ENTREE
//	@Failure      400  {object}  V1_API_RESPONSE_ENTREE
//	@Failure      500  {object}  V1_API_RESPONSE_ENTREE
//	@Router       /entree [post]
func CreateEntree(c *gin.Context) {
	// TODO: Add logic to reject unauthorized requests (and certainly do not deploy until all auth logic is wired up)
	response := V1_API_RESPONSE_ENTREE{}
	var status int
	var input models.Entree
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = "\"option_name\" is required"
	} else {
		entrees := []models.Entree{input}
		err := models.CreateEntrees(&entrees)
		if err != nil {
			status = http.StatusInternalServerError
			log.Println(err.Error())
			response.Message = "Internal server error"
		} else {
			status = http.StatusCreated
			response.Message = "Created entree"
			response.Data.Entrees = entrees
		}
	}
	response.Status = status
	c.JSON(status, response)
}

// DeleteEntree deletes an entree
//
//	@Summary      gets one or all entrees
//	@Description  gets 1 entree if an ID is found in the route, otherwise returns all entrees
//	@Tags         entrees
//	@Accept       json
//	@Produce      json
//	@Success      200  {object}  V1_API_RESPONSE_ENTREE
//	@Failure      500  {object}  V1_API_RESPONSE_ENTREE
//	@Router       /entree [delete]
func DeleteEntree(c *gin.Context) {
	// TODO: Add logic to reject unauthorized requests (and certainly do not deploy until all auth logic is wired up)
	response := V1_API_DELETE_RESPONSE{}
	var status int
	id, _ := uuid.Parse(c.Param("id"))
	result, err := models.DeleteEntree(id)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Internal server error"
	} else {
		status = http.StatusAccepted
		response.Message = "Deleted entree"
		response.Data = DeleteRecordResponse{
			DeletedRecords: int(*result),
		}
	}
	response.Status = status
	c.JSON(status, response)
}
