package stacks

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	portainer "github.com/portainer/portainer/api"
	httperrors "github.com/portainer/portainer/api/http/errors"
	"github.com/portainer/portainer/api/http/security"
	"github.com/portainer/portainer/api/internal/stackutils"
	"github.com/portainer/portainer/api/stacks"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	bolterrors "github.com/portainer/portainer/api/bolt/errors"
)

// @id StackStart
// @summary Starts a stopped Stack
// @description Starts a stopped Stack.
// @description **Access policy**: restricted
// @tags stacks
// @security jwt
// @param id path int true "Stack identifier"
// @success 200 {object} portainer.Stack "Success"
// @failure 400 "Invalid request"
// @failure 403 "Permission denied"
// @failure 404 "Not found"
// @failure 500 "Server error"
// @router /stacks/{id}/start [post]
func (handler *Handler) stackStart(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	stackID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid stack identifier route variable", err}
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	}

	stack, err := handler.DataStore.Stack().Stack(portainer.StackID(stackID))
	if err == bolterrors.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a stack with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a stack with the specified identifier inside the database", err}
	}

	endpoint, err := handler.DataStore.Endpoint().Endpoint(stack.EndpointID)
	if err == bolterrors.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an endpoint with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an endpoint with the specified identifier inside the database", err}
	}

	err = handler.requestBouncer.AuthorizedEndpointOperation(r, endpoint)
	if err != nil {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to access endpoint", err}
	}

	isUnique, err := handler.checkUniqueName(endpoint, stack.Name, stack.ID, stack.SwarmID != "")
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to check for name collision", err}
	}
	if !isUnique {
		errorMessage := fmt.Sprintf("A stack with the name '%s' is already running", stack.Name)
		return &httperror.HandlerError{http.StatusConflict, errorMessage, errors.New(errorMessage)}
	}

	resourceControl, err := handler.DataStore.ResourceControl().ResourceControlByResourceIDAndType(stackutils.ResourceControlID(stack.EndpointID, stack.Name), portainer.StackResourceControl)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve a resource control associated to the stack", err}
	}

	access, err := handler.userCanAccessStack(securityContext, endpoint.ID, resourceControl)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to verify user authorizations to validate stack access", err}
	}
	if !access {
		return &httperror.HandlerError{http.StatusForbidden, "Access denied to resource", httperrors.ErrResourceAccessDenied}
	}

	if stack.Status == portainer.StackStatusActive {
		return &httperror.HandlerError{http.StatusBadRequest, "Stack is already active", errors.New("Stack is already active")}
	}

	if stack.AutoUpdate != nil && stack.AutoUpdate.Interval != "" {
		if stack.AutoUpdate.JobID != "" {
			if err := handler.Scheduler.StopJob(stack.AutoUpdate.JobID); err != nil {
				log.Printf("[WARN] could not stop the job for the stack %v", stack.ID)
			}
		}

		d, err := time.ParseDuration(stack.AutoUpdate.Interval)
		if err != nil {
			return &httperror.HandlerError{StatusCode: http.StatusBadRequest, Message: "Unable to parse stack's auto update interval", Err: err}
		}
		jobID := handler.Scheduler.StartJobEvery(d, func() {
			if err := stacks.RedeployWhenChanged(stack.ID, handler.StackDeployer, handler.DataStore, handler.GitService); err != nil {
				log.Printf("[ERROR] %s\n", err)
			}
		})

		stack.AutoUpdate.JobID = jobID
	}

	err = handler.startStack(stack, endpoint)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to start stack", err}
	}

	stack.Status = portainer.StackStatusActive
	err = handler.DataStore.Stack().UpdateStack(stack.ID, stack)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to update stack status", err}
	}

	return response.JSON(w, stack)
}

func (handler *Handler) startStack(stack *portainer.Stack, endpoint *portainer.Endpoint) error {
	switch stack.Type {
	case portainer.DockerComposeStack:
		return handler.ComposeStackManager.Up(stack, endpoint)
	case portainer.DockerSwarmStack:
		return handler.SwarmStackManager.Deploy(stack, true, endpoint)
	}
	return nil
}
