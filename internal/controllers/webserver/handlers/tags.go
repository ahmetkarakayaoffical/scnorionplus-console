package handlers

import (
	"strconv"

	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/admin_views"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/partials"
	"github.com/labstack/echo/v4"
)

func (h *Handler) TagManager(c echo.Context) error {
	var err error

	commonInfo, err := h.GetCommonInfo(c)
	if err != nil {
		return err
	}

	p := partials.NewPaginationAndSort()
	p.GetPaginationAndSortParams(c.FormValue("page"), c.FormValue("pageSize"), c.FormValue("sortBy"), c.FormValue("sortOrder"), c.FormValue("currentSortBy"))

	if c.Request().Method == "POST" {
		tagId := c.FormValue("tagId")
		tag := c.FormValue("tag")
		description := c.FormValue("description")
		color := c.FormValue("color")

		if tag != "" && color != "" {
			if tagId == "" {
				if err := h.Model.NewTag(tag, description, color, commonInfo); err != nil {
					return RenderError(c, partials.ErrorMessage(err.Error(), false))
				}
			} else {
				id, err := strconv.Atoi(tagId)
				if err != nil {
					return RenderError(c, partials.ErrorMessage(err.Error(), false))
				}
				if err := h.Model.UpdateTag(id, tag, description, color, commonInfo); err != nil {
					return RenderError(c, partials.ErrorMessage(err.Error(), false))
				}
			}

		}
	}

	if c.Request().Method == "DELETE" {
		tagId := c.FormValue("tagId")
		if tagId == "" {
			return RenderError(c, partials.ErrorMessage("tag cannot be empty", false))
		}

		id, err := strconv.Atoi(tagId)
		if err != nil {
			return RenderError(c, partials.ErrorMessage(err.Error(), false))
		}

		if err := h.Model.DeleteTag(id, commonInfo); err != nil {
			return RenderError(c, partials.ErrorMessage(err.Error(), false))
		}
	}

	p.NItems, err = h.Model.CountAllTags(commonInfo)
	if err != nil {
		return RenderError(c, partials.ErrorMessage(err.Error(), false))
	}

	tags, err := h.Model.GetTagsByPage(p, commonInfo)
	if err != nil {
		return RenderError(c, partials.ErrorMessage(err.Error(), false))
	}

	agentsExists, err := h.Model.AgentsExists(commonInfo)
	if err != nil {
		return RenderError(c, partials.ErrorMessage(err.Error(), false))
	}

	serversExists, err := h.Model.ServersExists()
	if err != nil {
		return RenderError(c, partials.ErrorMessage(err.Error(), false))
	}

	return RenderView(c, admin_views.TagsIndex(" | Tags", admin_views.Tags(c, p, tags, agentsExists, serversExists, commonInfo, h.GetAdminTenantName(commonInfo)), commonInfo))
}
