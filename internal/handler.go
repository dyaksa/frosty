package internal

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kadzany/frosty/workflow"

	"github.com/gorilla/mux"
)

type WorkflowHandler struct {
	DB *sql.DB
}

func (wh *WorkflowHandler) CreateNode(resw http.ResponseWriter, req *http.Request) {
	node := workflow.Node{}
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&node); err != nil {
		responseError(resw, http.StatusBadRequest, "Invalid request payload")
	}
	defer req.Body.Close()

	id, err := workflow.CreateNode(wh.DB, node.Title, node.Type, node.Description)

	if err != nil {
		responseError(resw, http.StatusInternalServerError, err.Error())
		return
	}

	responseJson(resw, http.StatusCreated, id)
}

func (wh *WorkflowHandler) GetNode(resw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])

	if err != nil {
		responseError(resw, http.StatusBadRequest, "Invalid Node Id")
	}

	node, err := workflow.GetNode(wh.DB, id)
	if err != nil {
		responseError(resw, http.StatusInternalServerError, err.Error())
		return
	}

	responseJson(resw, http.StatusOK, node)
}

func (wh *WorkflowHandler) AddRelationship(resw http.ResponseWriter, req *http.Request) {
	var relationship workflow.NodeClosure
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&relationship); err != nil {
		responseError(resw, http.StatusBadRequest, "Invalid request payload")
	}
	defer req.Body.Close()

	err := workflow.AddRelationship(wh.DB, relationship.Ancestor, relationship.Descendant)
	if err != nil {
		responseError(resw, http.StatusInternalServerError, err.Error())
		return
	}
	responseJson(resw, http.StatusCreated, relationship)
}

func (wh *WorkflowHandler) ExecuteWorkflow(resw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])

	if err != nil {
		responseError(resw, http.StatusBadRequest, "Invalid Workflow Id")
		return
	}

	err = workflow.ExecuteWorkflow(wh.DB, id)
	if err != nil {
		responseError(resw, http.StatusInternalServerError, err.Error())
		return
	}

	responseJson(resw, http.StatusOK, nil)
}

func (wh *WorkflowHandler) CreateWorkflow(resw http.ResponseWriter, req *http.Request) {
	wf := workflow.Workflow{}
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&wf); err != nil {
		responseError(resw, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer req.Body.Close()

	id, err := workflow.CreateWorkflow(wh.DB, wf.Name, wf.Description, wf.StartingNodeID)
	if err != nil {
		responseError(resw, http.StatusInternalServerError, err.Error())
		return
	}

	responseJson(resw, http.StatusCreated, id)
}

func (wh *WorkflowHandler) CreateTask(resw http.ResponseWriter, req *http.Request) {
	task := workflow.Task{}
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&task); err != nil {
		responseError(resw, http.StatusBadRequest, "Invalid request payload")
	}
	defer req.Body.Close()

	id, err := workflow.CreateTask(wh.DB, task.Title, task.Type, task.HttpMethod, task.Action, task.Params, task.MaxRetries)
	if err != nil {
		responseError(resw, http.StatusInternalServerError, err.Error())
		return
	}

	responseJson(resw, http.StatusCreated, id)
}

func (wh *WorkflowHandler) AddTaskToNode(resw http.ResponseWriter, req *http.Request) {
	var nodeTask workflow.NodeTask
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&nodeTask); err != nil {
		responseError(resw, http.StatusBadRequest, "Invalid request payload")
	}
	defer req.Body.Close()

	err := workflow.AddTaskToNode(wh.DB, nodeTask.NodeID, nodeTask.TaskID, nodeTask.TaskOrder)
	if err != nil {
		responseError(resw, http.StatusInternalServerError, err.Error())
		return
	}

	responseJson(resw, http.StatusCreated, nodeTask)
}

func (wh *WorkflowHandler) CreateWorkflowExecution(resw http.ResponseWriter, req *http.Request) {
	var wfExec workflow.WorkflowExecution
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&wfExec); err != nil {
		responseError(resw, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer req.Body.Close()

	id, err := workflow.CreateWorkflowExecution(wh.DB, wfExec.WorkflowID, wfExec.ReferenceNumber)
	if err != nil {
		responseError(resw, http.StatusInternalServerError, err.Error())
		return
	}

	responseJson(resw, http.StatusCreated, id)
}

func (wh *WorkflowHandler) ExecuteWorkflowByExecutionID(resw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])

	if err != nil {
		responseError(resw, http.StatusBadRequest, "Invalid Workflow Execution Id")
		return
	}

	err = workflow.ExecuteWorkflowByExecutionID(wh.DB, id)
	if err != nil {
		responseError(resw, http.StatusInternalServerError, err.Error())
		return
	}

	responseJson(resw, http.StatusOK, nil)
}
