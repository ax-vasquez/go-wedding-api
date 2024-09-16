package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/ax-vasquez/wedding-site-api/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetEntrees gets one or all entrees
//
//	@Summary      	gets one or all entrees
//	@Description  	Gets the selected entree for the given user ID (empty array if no selection has been made), or a list of all available entrees if no user ID is provided
//	@Tags         	entrees
//	@Produce      	json
//	@Success      	200  {object}  types.V1_API_RESPONSE_ENTREE
//	@Failure      	500  {object}  types.V1_API_RESPONSE_ENTREE
//	@Param 			user_id  path string true "User ID" Format(uuid)
//	@Router       	/entrees [get]
//	@Router       	/user/{user_id}/entrees [get]
func GetEntrees(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	idStr := c.Param("id")
	response := types.V1_API_RESPONSE_ENTREE{}
	var status int
	var entrees []models.Entree
	// If an ID param was given, attempt lookup by the ID
	if len(idStr) > 0 {
		id, err := uuid.Parse(idStr)
		// If the given ID is invalid, return error response
		if err != nil {
			status = http.StatusBadRequest
			response.Status = status
			response.Message = err.Error()
			c.JSON(status, response)
			return
		}
		entree, err := models.FindEntreeById(ctx, id)
		entrees = append(entrees, *entree)
		// If an error occurs in the DB during lookup, return error response
		if err != nil {
			status = http.StatusInternalServerError
			log.Println(err.Error())
			response.Status = status
			response.Message = "Internal server error"
			c.JSON(status, response)
			return
		}
		status = http.StatusOK
		response.Status = status
		response.Data.Entrees = entrees
		c.JSON(status, response)
		return
	}
	// If no ID param was given, return all entrees (which will be empty should an error occur)
	entrees, err := models.FindEntrees(ctx)
	if err != nil {
		status = http.StatusInternalServerError
		log.Println(err.Error())
		response.Message = "Internal server error"
	} else {
		status = http.StatusOK
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
//	@Success      201  {object}  types.V1_API_RESPONSE_ENTREE
//	@Failure      400  {object}  types.V1_API_RESPONSE_ENTREE
//	@Failure      500  {object}  types.V1_API_RESPONSE_ENTREE
//	@Router       /entree [post]
func CreateEntree(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := types.V1_API_RESPONSE_ENTREE{}
	var status int
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
//	@Success      202  {object}  types.V1_API_RESPONSE_ENTREE
//	@Failure      500  {object}  types.V1_API_RESPONSE_ENTREE
//	@Router       /entree [delete]
func DeleteEntree(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := types.V1_API_DELETE_RESPONSE{}
	var status int
	id, _ := uuid.Parse(c.Param("id"))
	result, err := models.DeleteEntree(ctx, id)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Internal server error"
	} else {
		status = http.StatusAccepted
		response.Message = "Deleted entree"
		response.Data = types.DeleteRecordResponse{
			DeletedRecords: int(*result),
		}
	}
	response.Status = status
	c.JSON(status, response)
}
