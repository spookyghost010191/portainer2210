package registries

import (
	"net/http"

	portainer "github.com/portainer/portainer/api"
	httperrors "github.com/portainer/portainer/api/http/errors"
	httperror "github.com/portainer/portainer/pkg/libhttp/error"
	"github.com/portainer/portainer/pkg/libhttp/request"
	"github.com/portainer/portainer/pkg/libhttp/response"
)

// @id RegistryInspect
// @summary Inspect a registry
// @description Retrieve details about a registry.
// @description **Access policy**: restricted
// @tags registries
// @security ApiKeyAuth
// @security jwt
// @produce json
// @param id path int true "Registry identifier"
// @success 200 {object} portainer.Registry "Success"
// @failure 400 "Invalid request"
// @failure 403 "Permission denied to access registry"
// @failure 404 "Registry not found"
// @failure 500 "Server error"
// @router /registries/{id} [get]
func (handler *Handler) registryInspect(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	registryID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("Invalid registry identifier route variable", err)
	}

	registry, err := handler.DataStore.Registry().Read(portainer.RegistryID(registryID))
	if handler.DataStore.IsErrObjectNotFound(err) {
		return httperror.NotFound("Unable to find a registry with the specified identifier inside the database", err)
	} else if err != nil {
		return httperror.InternalServerError("Unable to find a registry with the specified identifier inside the database", err)
	}

	hasAccess, isAdmin, err := handler.userHasRegistryAccess(r, registry)
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve info from request context", err)
	}
	if !hasAccess {
		return httperror.Forbidden("Access denied to resource", httperrors.ErrResourceAccessDenied)
	}

	hideFields(registry, !isAdmin)
	return response.JSON(w, registry)
}
