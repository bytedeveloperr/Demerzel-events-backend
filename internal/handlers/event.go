package handlers

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

func GetGroupEventsHandler(c *gin.Context) {

	id := c.Param("id")

	group := models.Group{
		ID: id,
	}

	events, err := group.GetGroupEvents(db.DB)

	if err != nil {
		response.Error(c, 500, "Can't process your request")
		return
	}

	response.Success(c, 200, "List of events", map[string]interface{}{"events": events})
}

func CreateEventHandler(c *gin.Context) {
	var input models.NewEvent

	// Error if JSON request is invalid
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Unable to parse payload: %s", err.Error()))
		return
	}

	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "Unable to read user from context")
		return
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid context user type")
		return
	}

	input.CreatorId = user.Id

	// Check if description field is empty or is a string
	if input.Description == "" {
		response.Error(c, http.StatusBadRequest, "Description field is empty")
		return
	}

	if reflect.ValueOf(input.Description).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "Description is not a string")
		return
	}

	// Check if thumbnail field is empty or is a string
	if input.Thumbnail == "" {
		response.Error(c, http.StatusBadRequest, "Thumbnail field is empty")
		return
	}

	if reflect.ValueOf(input.Thumbnail).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "Thumbnail is not a string")
		return
	}

	// Check if location field is empty or is a string
	if input.Location == "" {
		response.Error(c, http.StatusBadRequest, "Location field is empty")
		return
	}

	if reflect.ValueOf(input.Location).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "Location is not a string")
		return
	}

	// Check if title field is empty or is a string
	if input.Title == "" {
		response.Error(c, http.StatusBadRequest, "Title field is empty")
		return
	}

	if reflect.ValueOf(input.Title).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "Title is not a string")
		return
	}

	// Check if start_time field is empty or is a string
	if input.StartTime == "" {
		response.Error(c, http.StatusBadRequest, "StartTime field is empty")
		return
	}

	if reflect.ValueOf(input.StartTime).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "StartTime is not a string")
		return
	}

	// Check if end_time field is empty or is a string
	if input.EndTime == "" {
		response.Error(c, http.StatusBadRequest, "EndTime field is empty")
		return
	}

	if reflect.ValueOf(input.EndTime).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "EndTime is not a string")
		return
	}

	// Check if start_date field is empty or is a string
	if input.StartDate == "" {
		response.Error(c, http.StatusBadRequest, "StartDate field is empty")
		return
	}

	if !models.IsValidDate(input.StartDate) {
		response.Error(c, http.StatusBadRequest, "Invalid StartDate. Should follow format 2023-09-21")
		return
	}

	if reflect.ValueOf(input.StartDate).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "StartDate is not a string")
		return
	}
	// Check if end_date field is empty or is a string
	if input.EndDate == "" {
		response.Error(c, http.StatusBadRequest, "EndDate field is empty")
		return
	}

	if !models.IsValidDate(input.EndDate) {
		response.Error(c, http.StatusBadRequest, "Invalid EndDate. Should follow format 2023-09-21")
		return
	}

	if reflect.ValueOf(input.EndDate).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "EndDate is not a string")
		return
	}

	createdEvent, err := models.CreateEvent(db.DB, &input)

	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Event Created", map[string]interface{}{"event": createdEvent})

}

func GetEventHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	if eventID == "" {
		response.Error(c, http.StatusBadRequest, "Event ID is required")
		return
	}

	event, err := models.GetEventByID(db.DB, eventID)

	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Event details fetched", map[string]interface{}{"event": event})
}

func JoinEventHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	if eventID == "" {
		response.Error(c, http.StatusBadRequest, "Event ID is required")
		return
	}

	_, err := models.GetEventByID(db.DB, eventID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "Unable to read user from context")
		return
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid context user type")
		return
	}

	event, err := models.AttachUserToEvent(db.DB, user.Id, eventID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Unable to join event:"+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Joined Event", map[string]*models.Event{"event": event})
}

func LeaveEventHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	if eventID == "" {
		response.Error(c, http.StatusBadRequest, "Event ID is required")
		return
	}

	_, err := models.GetEventByID(db.DB, eventID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "Unable to read user from context")
		return
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid context user type")
		return
	}

	event, err := models.DetachUserFromEvent(db.DB, user.Id, eventID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Unable to leave event:"+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Removed from Event", map[string]*models.Event{"event": event})

}

// ListEventsHandler lists all events
func ListEventsHandler(c *gin.Context) {
	startDate := c.Query("start_date")
	events, err := models.ListEvents(db.DB, startDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.Success(c, http.StatusOK, "Events retrieved successfully", map[string]interface{}{"events": events})
}

func ListFriendsEventsHandler(c *gin.Context) {
	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "Unable to read user from context")
		return
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid context user type")
		return
	}

	userGroups, _, err := services.GetGroupsByUserId(user.Id)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Unable to get groups which user belongs to:"+err.Error())
		return
	}

	if len(userGroups) == 0 {
		events := make([]models.Event, 0)
		response.Success(c, http.StatusOK, "Friend Events", map[string]interface{}{"events": events})

		return
	}

	var userGroupIds []string
	for _, group := range userGroups {
		userGroupIds = append(userGroupIds, group.ID)
	}

	events, err := models.ListEventsInGroups(db.DB, userGroupIds)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Unable to get events: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Events", map[string]interface{}{"events": events})

	return
}

func SubscribeUserToEvent(c *gin.Context) {
	eventID := c.Param("id")
	rawUser, exists := c.Get("user")

	if !exists {
		response.Error(c, http.StatusConflict, "error: unable to retrieve user from context")
		return
	}

	user, ok := rawUser.(*models.User)

	if !ok {
		response.Error(c, http.StatusConflict, "error: invalid user type in context")
		return
	}

	event, err := models.GetEventByID(db.DB, eventID)
	if event == nil {
		response.Error(c, http.StatusNotFound, "Event does not exist")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = models.SubscribeUserToEvent(db.DB, user.Id, eventID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User successfully subscribed to event", nil)
}

func UnsubscribeFromEvent(c *gin.Context) {
	eventID := c.Param("id")
	rawUser, exists := c.Get("user")

	if !exists {
		response.Error(c, http.StatusConflict, "error: unable to retrieve user from context")
		return
	}
	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusConflict, "error: invalid user type in context")
		return
	}

	event, err := models.GetEventByID(db.DB, eventID)
	if event == nil {
		response.Error(c, http.StatusNotFound, "Event does not exist")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = models.UnsubscribeUserFromEvent(db.DB, user.Id, eventID)
	if err != nil {
		response.Error(c, http.StatusConflict, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User successfully unsubscribed to event", nil)
}

func GetUserEventSubscriptions(c *gin.Context) {
	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusConflict, "error: unable to retrieve user from context")
		return
	}
	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusConflict, "error: invalid user type in context")
		return
	}

	events, err := models.GetUserEventSubscriptions(db.DB, user.Id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User event subscriptions retrieved", map[string]interface{}{"events": events})
}
