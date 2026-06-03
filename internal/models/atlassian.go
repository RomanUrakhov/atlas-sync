package models

type GoalStatus struct {
	Value string `json:"value"`
}

type GoalTargetDate struct {
	Label string `json:"label"`
}

type Goal struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Status     GoalStatus     `json:"status"`
	TargetDate GoalTargetDate `json:"targetDate"`
}

type GoalEdge struct {
	Node Goal `json:"node"`
}

type GoalSearch struct {
	Edges    []GoalEdge `json:"edges"`
	PageInfo PageInfo   `json:"pageInfo"`
}

type GoalsData struct {
	GoalSearch GoalSearch `json:"goals_search"`
}

type ProjectDescription struct {
	What string `json:"what"`
	Why  string `json:"why"`
}

type ProjectDueDate struct {
	Label string `json:"label"`
}

type Project struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description ProjectDescription `json:"description"`
	DueDate     ProjectDueDate     `json:"dueDate"`
	Archived    bool               `json:"archived"`
}

type ProjectEdge struct {
	Node Project `json:"node"`
}

type ProjectSearch struct {
	Edges    []ProjectEdge `json:"edges"`
	PageInfo PageInfo      `json:"pageInfo"`
}

type ProjectsData struct {
	ProjectSearch ProjectSearch `json:"projects_search"`
}

type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}
