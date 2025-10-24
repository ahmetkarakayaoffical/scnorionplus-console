package handlers

import (
	"fmt"

	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/admin_views"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/partials"
	"github.com/invopop/ctxi18n/i18n"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ListSessions(c echo.Context, successMessage string) error {
	var err error

	commonInfo, err := h.GetCommonInfo(c)
	if err != nil {
		return err
	}

	errMessage := ""

	p := partials.NewPaginationAndSort()
	p.GetPaginationAndSortParams(c.FormValue("page"), c.FormValue("pageSize"), c.FormValue("sortBy"), c.FormValue("sortOrder"), c.FormValue("currentSortBy"))

	p.NItems, err = h.Model.CountAllSessions()
	if err != nil {
		errMessage = err.Error()
	}

	s, err := h.Model.GetSessionsByPage(p)
	if err != nil {
		successMessage = ""
		errMessage = err.Error()
	}

	agentsExists, err := h.Model.AgentsExists(commonInfo)
	if err != nil {
		return RenderError(c, partials.ErrorMessage(err.Error(), false))
	}

	serversExists, err := h.Model.ServersExists()
	if err != nil {
		return RenderError(c, partials.ErrorMessage(err.Error(), false))
	}

	return RenderView(c, admin_views.SessionsIndex(" | Sessions", admin_views.Sessions(c, p, s, successMessage, errMessage, h.SessionManager.Manager.Codec, agentsExists, serversExists, commonInfo), commonInfo))
}

func (h *Handler) SessionDelete(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		return RenderError(c, partials.ErrorMessage("no token was found in request", true))
	}

	return RenderConfirm(c, partials.ConfirmDelete(c, i18n.T(c.Request().Context(), "confirm.session_delete"), "/admin/sessions", fmt.Sprintf("/admin/sessions/%s", token)))
}

func (h *Handler) SessionConfirmDelete(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		return RenderError(c, partials.ErrorMessage("no token was found in request", true))
	}

	if err := h.Model.DeleteSession(token); err != nil {
		return RenderError(c, partials.ErrorMessage(err.Error(), true))
	}

	return h.ListSessions(c, i18n.T(c.Request().Context(), "success.session_delete"))
}
