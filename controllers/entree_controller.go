package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ax-vasquez/wedding-site-api/helper"
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
//	@Description  	Gets the selected entree for the given user ID (empty array if no selection has been made), or a list of all available entrees if no user ID is provided
//	@Tags         	entrees
//	@Produce      	json
//	@Success      	200  {object}  V1_API_RESPONSE_ENTREE
//	@Failure      	500  {object}  V1_API_RESPONSE_ENTREE
//	@Param 			user_id  path string true "User ID" Format(uuid)
//	@Router       	/entrees [get]
//	@Router       	/user/{user_id}/entrees [get]
func GetEntrees(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	response := V1_API_RESPONSE_ENTREE{}
	var status int
	var entrees []models.Entree
	// If no error occurred, the parse was successful, meaning a UUID was found and results will be filtered for the given user
	if err == nil {
		// Since we have an ID, we need to check that the user exists before continuing with controller logic.
		if err := helper.MatchUserTypeToUid(c.Request.Context(), id.String()); err != nil {
			status = http.StatusBadRequest
			response.Message = err.Error()
		} else {
			status = http.StatusOK
			entrees, err = models.FindEntreesForUser(ctx, id)
			if err != nil {
				status = http.StatusInternalServerError
				log.Println(err.Error())
				response.Message = "Internal server error"
			}
		}
		// If an error occurred, we ignore it and assume it's because there was no ID in the path - all results will be returned
	} else {
		entrees, err = models.FindEntrees(ctx)
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
//	@Summary      create entree
//	@Description  Create a new entree and return the new record's data to the caller
//	@Tags         entrees
//	@Accept       json
//	@Produce      json
//	@Param		  data body models.Entree true "The input entree data (only `option_name` is required)"
//	@Success      201  {object}  V1_API_RESPONSE_ENTREE
//	@Failure      400  {object}  V1_API_RESPONSE_ENTREE
//	@Failure      500  {object}  V1_API_RESPONSE_ENTREE
//	@Router       /entree [post]
func CreateEntree(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := V1_API_RESPONSE_ENTREE{}
	var status int
	if err := helper.CheckUserType(c.Request.Context(), "ADMIN"); err != nil {
		status = http.StatusUnauthorized
		response.Status = status
		c.JSON(status, response)
		return
	}
	var input models.Entree
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else {
		entrees := []models.Entree{input}
		err := models.CreateEntrees(ctx, &entrees)
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
//	@Summary      deletes an entree
//	@Description  Deletes an entree and returns a response to indicate success or failure
//	@Tags         entrees
//	@Produce      json
//	@Param 		  id  path string true "Entree ID" Format(uuid)
//	@Success      202  {object}  V1_API_RESPONSE_ENTREE
//	@Failure      500  {object}  V1_API_RESPONSE_ENTREE
//	@Router       /entree [delete]
func DeleteEntree(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := V1_API_DELETE_RESPONSE{}
	var status int
	if err := helper.CheckUserType(c.Request.Context(), "ADMIN"); err != nil {
		status = http.StatusUnauthorized
		response.Status = status
		c.JSON(status, response)
		return
	}
	id, _ := uuid.Parse(c.Param("id"))
	result, err := models.DeleteEntree(ctx, id)
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
