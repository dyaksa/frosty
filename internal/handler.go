package internal

import (
	"database/sql"
	"encoding/json"
	"log"
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

	log.Println("Creating Node")
	log.Println(req.Body)

	if err := decoder.Decode(&node); err != nil {
		responseError(resw, http.StatusBadRequest, "Invalid request payload")
	}
	defer req.Body.Close()

	id, err := workflow.CreateNode(wh.DB, node.Title, node.Type, node.Description)

	log.Println("Created Node")
	log.Println(`id: `, id)
	log.Println(`node: `, node)
	log.Println(`err: `, err)

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

	log.Println("Adding Relationship")

	if err := decoder.Decode(&relationship); err != nil {
		responseError(resw, http.StatusBadRequest, "Invalid request payload")
	}
	defer req.Body.Close()

	log.Println(`Ancestor: `, relationship.Ancestor)
	log.Println(`Descendant: `, relationship.Descendant)

	err := workflow.AddRelationship(wh.DB, relationship.Ancestor, relationship.Descendant)
	if err != nil {
		responseError(resw, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("Added Relationship")

	responseJson(resw, http.StatusCreated, relationship)
}

func (wh *WorkflowHandler) GetDescendants(resw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])

	if err != nil {
		responseError(resw, http.StatusBadRequest, "Invalid Node Id")
	}

	descendants, err := workflow.GetDescendants(wh.DB, id)
	if err != nil {
		responseError(resw, http.StatusInternalServerError, err.Error())
		return
	}

	responseJson(resw, http.StatusOK, descendants)
}

func (wh *WorkflowHandler) ExecuteWorkflow(resw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])

	if err != nil {
		responseError(resw, http.StatusBadRequest, "Invalid Node Id")
	}

	err = workflow.ValidateWorkflow(wh.DB, id)
	if err != nil {
		responseError(resw, http.StatusInternalServerError, err.Error())
		return
	}

	err = workflow.ExecuteWorkflow(wh.DB, id, func(node workflow.Node) error {
		return nil
	})
	if err != nil {
		responseError(resw, http.StatusInternalServerError, err.Error())
		return
	}

	responseJson(resw, http.StatusOK, nil)
}

func (wh *WorkflowHandler) RollbackWorkflow(resw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])

	if err != nil {
		responseError(resw, http.StatusBadRequest, "Invalid Node Id")
	}

	err = workflow.RollbackWorkflow(wh.DB, id, "start", func(node workflow.Node) error {
		return nil
	})
	if err != nil {
		responseError(resw, http.StatusInternalServerError, err.Error())
		return
	}

	responseJson(resw, http.StatusOK, nil)
}
