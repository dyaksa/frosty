package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kadzany/frosty/workflow"
	"github.com/stretchr/testify/assert"
)

func TestWorkflowHandler_CreateNode(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	handler := WorkflowHandler{DB: db}
	node := workflow.Node{
		Title:       "Test Node",
		Type:        "Task",
		Description: "Test Description",
	}
	nodeJSON, _ := json.Marshal(node)

	mock.ExpectQuery("INSERT INTO nodes").
		WithArgs("Test Node", "Task", "Test Description").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New().String()))

	req, _ := http.NewRequest("POST", "/workflow/node", bytes.NewBuffer(nodeJSON))
	resw := httptest.NewRecorder()

	handler.CreateNode(resw, req)

	assert.Equal(t, http.StatusCreated, resw.Code)
}

func TestWorkflowHandler_GetNode(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	handler := WorkflowHandler{DB: db}
	nodeID := uuid.New()

	mock.ExpectQuery("SELECT id::uuid, title, type, description, created_at, updated_at, deleted_at FROM nodes WHERE id = ?").
		WithArgs(nodeID.String()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "type", "description", "created_at", "updated_at", "deleted_at"}).
			AddRow(nodeID.String(), "Test Node", "Task", "Test Description", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), nil))

	req, _ := http.NewRequest("GET", "/workflow/node/"+nodeID.String(), nil)
	resw := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"id": nodeID.String()})

	handler.GetNode(resw, req)

	assert.Equal(t, http.StatusOK, resw.Code)
}

func TestWorkflowHandler_AddRelationship(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	handler := WorkflowHandler{DB: db}
	relationship := workflow.NodeClosure{
		Ancestor:   uuid.New(),
		Descendant: uuid.New(),
	}
	relationshipJSON, _ := json.Marshal(relationship)

	mock.ExpectExec("INSERT INTO node_closure").WillReturnResult(sqlmock.NewResult(1, 1))

	req, _ := http.NewRequest("POST", "/workflow/node/relationship", bytes.NewBuffer(relationshipJSON))
	resw := httptest.NewRecorder()

	handler.AddRelationship(resw, req)

	assert.Equal(t, http.StatusCreated, resw.Code)
}

func TestWorkflowHandler_CreateTask(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	handler := WorkflowHandler{DB: db}
	task := workflow.Task{
		Title:      "Test Task",
		Type:       "API",
		HttpMethod: "POST",
		Action:     "http://example.com/api",
		Params:     "{}",
		MaxRetries: 3,
	}
	taskJSON, _ := json.Marshal(task)

	mock.ExpectQuery("INSERT INTO tasks").
		WithArgs("Test Task", "API", "POST", "http://example.com/api", "{}", 3).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New().String()))

	req, _ := http.NewRequest("POST", "/workflow/task", bytes.NewBuffer(taskJSON))
	resw := httptest.NewRecorder()

	handler.CreateTask(resw, req)

	assert.Equal(t, http.StatusCreated, resw.Code)
}

func TestWorkflowHandler_AddTaskToNode(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	handler := WorkflowHandler{DB: db}
	nodeTask := workflow.NodeTask{
		NodeID:     uuid.New(),
		TaskID:     uuid.New(),
		TaskOrder:  1,
	}
	nodeTaskJSON, _ := json.Marshal(nodeTask)

	mock.ExpectExec("INSERT INTO node_tasks").WillReturnResult(sqlmock.NewResult(1, 1))

	req, _ := http.NewRequest("POST", "/workflow/node/task", bytes.NewBuffer(nodeTaskJSON))
	resw := httptest.NewRecorder()

	handler.AddTaskToNode(resw, req)

	assert.Equal(t, http.StatusCreated, resw.Code)
}
