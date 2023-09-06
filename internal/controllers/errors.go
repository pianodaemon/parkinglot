package controllers

type ApiErrorCode int

const (
	Success ApiErrorCode = 0

	EndPointNoCarsYet      ApiErrorCode = 1001
	EndPointCarNotFound    ApiErrorCode = 1002
	EndPointCarNotCreated  ApiErrorCode = 1003
	EndPointFailedDeletion ApiErrorCode = 1004
)
