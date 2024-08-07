package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"tasks_microservice/internal/database"
	"tasks_microservice/internal/models"

	"github.com/streadway/amqp"
)

type Publisher interface {
	Publish(queueName string, body []byte, replyTo string, correlationId string) error
}

type Handlers struct {
	DB        *database.Database
	Publisher Publisher
}

type TasksValueRequest struct {
	Data models.TasksValue `json:"data"`
}

type TasksValueResponse struct {
    ID          int    `json:"valueId"`
    Type        string `json:"type"`
    StringValue *string `json:"stringValue,omitempty"` // Use pointer to differentiate between null and empty string
    NumberValue *int    `json:"numberValue,omitempty"` // Use pointer to differentiate between null and zero
    ColumnID    int    `json:"columnId"`
    TaskID      int    `json:"taskId"`
    OptionID    *int    `json:"optionId,omitempty"`
}


func NewHandlers(db *database.Database, publisher Publisher) *Handlers {
	return &Handlers{DB: db, Publisher: publisher}
}

func (h *Handlers) HandleTaskValues(d amqp.Delivery) {
	var requestData TasksValueRequest
	if err := json.Unmarshal(d.Body, &requestData); err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		d.Nack(false, false)
		return
	}

	switch requestData.Data.Action {
	case "create":
		log.Println("Handling create action")
		h.CreateTasksValue(d)
	case "update":
		log.Println("Handling update action")
		h.UpdateTasksValue(d)
	case "delete":
		log.Println("Handling delete action")
		h.DeleteTasksValue(d)
	case "get":
		log.Println("Handling get action")
		h.GetTasksValue(d)
	case "deleteAllForTask":
        log.Println("Handling delete all for task action")
        h.DeleteAllValuesForTask(d)  
	default:
		log.Printf("Unknown action: %s", requestData.Data.Action)
		d.Nack(false, true)
	}
}

func (h *Handlers) CreateTasksValue(d amqp.Delivery) {
	var requestData TasksValueRequest
	if err := json.Unmarshal(d.Body, &requestData); err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	value := requestData.Data

	if err := value.Validate(); err != nil {
		log.Printf("Validation error: %v", err)
		return
	}

	query := `INSERT INTO tasks_values (type, string_value, number_value, column_id, task_id, option_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := h.DB.Conn.QueryRow(query, value.Type, value.StringValue, value.NumberValue, value.ColumnID, value.TaskID, value.OptionID).Scan(&value.ID)
	if err != nil {
		log.Printf("Error inserting tasks value: %v", err)
		return
	}

	log.Printf("Successfully inserted tasks value: %+v", value)
}

func (h *Handlers) UpdateTasksValue(d amqp.Delivery) {
	var requestData TasksValueRequest
	if err := json.Unmarshal(d.Body, &requestData); err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	value := requestData.Data

	if err := value.ValidateUpdate(); err != nil {
		log.Printf("Validation error: %v", err)
		return
	}

	query := `SELECT id FROM tasks_values WHERE id = $1`
	var existingID int
	err := h.DB.Conn.QueryRow(query, value.ID).Scan(&existingID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Tasks value with ID %d does not exist", value.ID)
		} else {
			log.Printf("Error checking existence of tasks value: %v", err)
		}
		return
	}

	updateQuery := `UPDATE tasks_values SET type = $1, string_value = $2, number_value = $3, option_id = $4 WHERE id = $5`
	_, err = h.DB.Conn.Exec(updateQuery, value.Type, value.StringValue, value.NumberValue, value.OptionID, value.ID)
	if err != nil {
		log.Printf("Error updating tasks value: %v", err)
		return
	}

	log.Printf("Successfully updated tasks value: %+v", value)
}

func (h *Handlers) DeleteTasksValue(d amqp.Delivery) {
	var requestData TasksValueRequest
	if err := json.Unmarshal(d.Body, &requestData); err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	value := requestData.Data

	query := `SELECT id FROM tasks_values WHERE id = $1`
	var existingID int
	err := h.DB.Conn.QueryRow(query, value.ID).Scan(&existingID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Tasks value with ID %d does not exist", value.ID)
		} else {
			log.Printf("Error checking existence of tasks value: %v", err)
		}
		return
	}

	deleteQuery := `DELETE FROM tasks_values WHERE id = $1`
	_, err = h.DB.Conn.Exec(deleteQuery, value.ID)
	if err != nil {
		log.Printf("Error deleting tasks value: %v", err)
		return
	}
	log.Printf("Successfully deleted tasks value with ID %d", value.ID)
}

func (h *Handlers) GetTasksValue(d amqp.Delivery) {
    var requestData TasksValueRequest
    if err := json.Unmarshal(d.Body, &requestData); err != nil {
        log.Printf("Error unmarshalling message: %v", err)
        return
    }

    value := requestData.Data
    if h.DB == nil || h.DB.Conn == nil {
        log.Fatal("DB or DB.Conn is nil")
        return
    }

    if err := value.ValidateGet(); err != nil {
        log.Printf("Validation error: %v", err)
        return
    }

    var retrievedValue models.TasksValue
    query := `SELECT id, type, string_value, number_value, column_id, task_id, option_id FROM tasks_values WHERE id = $1`

    err := h.DB.Conn.QueryRow(query, value.ID).Scan(
        &retrievedValue.ID,
        &retrievedValue.Type,
        &retrievedValue.StringValue,
        &retrievedValue.NumberValue,
        &retrievedValue.ColumnID,
        &retrievedValue.TaskID,
        &retrievedValue.OptionID,
    )

    if err == sql.ErrNoRows {
        log.Printf("No record found with ID: %d", value.ID)
        err = h.Publisher.Publish(d.ReplyTo, []byte(`{"error": "Record not found"}`), "", d.CorrelationId)
        if err != nil {
            log.Printf("Failed to send not found response: %v", err)
        }
        return
    } else if err != nil {
        log.Printf("Error retrieving tasks value: %v", err)
        return
    }

    response := TasksValueResponse{
        ID:          retrievedValue.ID,
        Type:        retrievedValue.Type,
        StringValue: retrievedValue.StringValue,
        NumberValue: retrievedValue.NumberValue,
        ColumnID:    retrievedValue.ColumnID,
        TaskID:      retrievedValue.TaskID,
        OptionID:    retrievedValue.OptionID,
    }

    responseBytes, err := json.Marshal(response)
    if err != nil {
        log.Printf("Error marshalling response: %v", err)
        return
    }

    err = h.Publisher.Publish(d.ReplyTo, responseBytes, "", d.CorrelationId)
    if err != nil {
        log.Printf("Error publishing response: %v", err)
    }

    log.Printf("Successfully retrieved and sent tasks value: %+v", response)
}

func (h *Handlers) DeleteAllValuesForTask(d amqp.Delivery) {
    var requestData TasksValueRequest
    if err := json.Unmarshal(d.Body, &requestData); err != nil {
        log.Printf("Error unmarshalling message: %v", err)
        d.Nack(false, false)
        return
    }

    value := requestData.Data
    if h.DB == nil || h.DB.Conn == nil {
        log.Fatal("DB or DB.Conn is nil")
        return
    }

    query := `DELETE FROM tasks_values WHERE task_id = $1`
    _, err := h.DB.Conn.Exec(query, value.TaskID)
    if err != nil {
        log.Printf("Error deleting all values for task ID %d: %v", value.TaskID, err)
        d.Nack(false, false)
        return
    }

    log.Printf("Successfully deleted all values for task ID %d", value.TaskID)
}