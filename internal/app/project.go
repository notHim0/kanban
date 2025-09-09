package app

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/notHim0/kanban/internal/models"
	"github.com/notHim0/kanban/internal/types"
	"github.com/notHim0/kanban/internal/utils"
)

//creates a project
func (app *App) CreateProject(w http.ResponseWriter, r *http.Request){
	var project models.Project
	
	err := json.NewDecoder(r.Body).Decode(&project)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	claims := r.Context().Value("claims").(*types.Claims)
	
	userId := claims.Id

	var id string

	err = app.DB.QueryRow(`INSERT INTO projects ("user_id", name, repo_url, site_url, description, dependencies, dev_dependencies, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, userId, project.Name, project.RepoUrl, project.SiteUrl, 
	project.Description, pq.Array(project.Dependencies), pq.Array(project.DevDependencies), project.Status).Scan(&id)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating projects")
		return
	}

	project.Id = id
	project.UserId = userId

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(project)
}

//fetches all projects of a user
func (app *App) GetProjects(w http.ResponseWriter, r *http.Request){
	claims := r.Context().Value("claims").(*types.Claims)
	userId, err:= strconv.Atoi(claims.Id)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return 
	}

	rows, err := app.DB.Query(`SELECT id, "user_id", name, repo_url, site_url, description,
	dependencies, dev_dependencies, status FROM projects WHERE "user_id"=$1`, userId)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching projects")
		return
	}

	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var project models.Project
		var dependencies, devDependencies [] string

		err = rows.Scan(&project.Id, &project.UserId, &project.Name, &project.RepoUrl, &project.SiteUrl, &project.Description,
		pq.Array(&dependencies), pq.Array(&devDependencies), &project.Status)
	
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error scanning project")
			return
		}

		project.Dependencies = dependencies
		project.DevDependencies = devDependencies
	
		projects = append(projects, project)
	}
	
	err = rows.Err()

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching projects")
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

//fetches a specific project of a user
func (app *App) GetProject(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	projectId := vars["id"]

	claims := r.Context().Value("claims").(*types.Claims)
	userId, err := strconv.Atoi(claims.Id)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	var project models.Project 

	var dependencies, devDependencies []string

	err = app.DB.QueryRow(`SELECT id, user_id, name, repo_url, site_url, description, dependencies, dev_dependencies, status FROM projects WHERE user_id=$1 AND id=$2`,
	userId, projectId).Scan(&project.Id, &project.UserId, &project.Name, &project.RepoUrl, &project.SiteUrl, &project.Description, pq.Array(&dependencies), pq.Array(&devDependencies), &project.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Project not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching projects")
		return
	}

	project.Dependencies = dependencies
	project.DevDependencies = devDependencies
	
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(project)

}

//Updates a project 
func (app *App) UpdateProject(w http.ResponseWriter, r *http.Request){
	var project models.Project 
	
	err := json.NewDecoder(r.Body).Decode(&project)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return 
	}

	vars := mux.Vars(r)
	projectId := vars["id"]

	claims := r.Context().Value("claims").(*types.Claims)
	userId, er := strconv.Atoi(claims.Id)

	if er != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return 
	}

	var storedUserId int

	err = app.DB.QueryRow(`SELECT user_id FROM projects WHERE id=$1`, projectId).Scan(&storedUserId)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Project not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching project")
		return
	}

	if storedUserId != userId {
		utils.RespondWithError(w, http.StatusForbidden, "You do not have to permission to update this project")
		return
	}

	_, err = app.DB.Exec(`UPDATE projects SET name=$1, repo_url=$2, site_url=$3, description=$4, dependencies=$5,
		dev_dependencies=$6, status=$7 WHERE id=$8 AND user_id=$9`, project.Name, project.RepoUrl, project.SiteUrl, project.Description,
		pq.Array(project.Dependencies), pq.Array(project.DevDependencies), project.Status, projectId, userId)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating project")
		return 
	}

	project.Id = projectId
	project.UserId = string(rune(userId))
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(project)
}

//deletes a project
func (app *App) DeleteProject(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	projectId := vars["id"]

	claims := r.Context().Value("claims").(*types.Claims)
	userId := claims.Id

	var storedUserId string

	err := app.DB.QueryRow(`SELECT user_id FROM projects WHERE id=$1`, projectId).Scan(&storedUserId)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Project not found")
			return 
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching project")
		return
	}

	if storedUserId != userId {
		utils.RespondWithError(w, http.StatusForbidden, "You do not have permission to delete this project")
		return
	}

	_, err = app.DB.Exec("DELETE FROM projects WHERE id=$1 AND user_id=$2", projectId, userId)
	
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting project")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}