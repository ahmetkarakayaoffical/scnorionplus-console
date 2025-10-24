package handlers

import (
	"fmt"
	"log"

	scnorionplus_nats "github.com/ahmetkarakayaoffical/nats"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/filters"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/partials"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/security_views"
	"github.com/invopop/ctxi18n/i18n"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ListAntivirusStatus(c echo.Context) error {
	var err error

	commonInfo, err := h.GetCommonInfo(c)
	if err != nil {
		return err
	}

	p := partials.NewPaginationAndSort()
	p.GetPaginationAndSortParams(c.FormValue("page"), c.FormValue("pageSize"), c.FormValue("sortBy"), c.FormValue("sortOrder"), c.FormValue("currentSortBy"))

	// Get filters values
	f, availableOSes, detectedAntiviri, err := h.GetAntiviriFilters(c)
	if err != nil {
		return RenderError(c, partials.ErrorMessage(err.Error(), false))
	}

	antiviri, err := h.Model.GetAntiviriByPage(p, *f, commonInfo)
	if err != nil {
		return RenderError(c, partials.ErrorMessage(err.Error(), false))
	}

	p.NItems, err = h.Model.CountAllAntiviri(*f, commonInfo)
	if err != nil {
		return RenderError(c, partials.ErrorMessage(err.Error(), false))
	}

	refreshTime, err := h.Model.GetDefaultRefreshTime()
	if err != nil {
		log.Println("[ERROR]: could not get refresh time from database")
		refreshTime = 5
	}

	return RenderView(c, security_views.SecurityIndex("| Security", security_views.Antivirus(c, p, *f, antiviri, detectedAntiviri, availableOSes, refreshTime, commonInfo), commonInfo))
}

func (h *Handler) ListSecurityUpdatesStatus(c echo.Context) error {
	var err error

	commonInfo, err := h.GetCommonInfo(c)
	if err != nil {
		return err
	}

	p := partials.NewPaginationAndSort()
	p.GetPaginationAndSortParams(c.FormValue("page"), c.FormValue("pageSize"), c.FormValue("sortBy"), c.FormValue("sortOrder"), c.FormValue("currentSortBy"))

	// Get filters values
	f, availableOSes, availableUpdateStatus, err := h.GetSystemUpdatesFilters(c)
	if err != nil {
		return RenderError(c, partials.ErrorMessage(err.Error(), false))
	}

	systemUpdates, err := h.Model.GetSystemUpdatesByPage(p, *f, commonInfo)
	if err != nil {
		return RenderError(c, partials.ErrorMessage(err.Error(), false))
	}

	p.NItems, err = h.Model.CountAllSystemUpdates(*f, commonInfo)
	if err != nil {
		return RenderError(c, partials.ErrorMessage(err.Error(), false))
	}

	refreshTime, err := h.Model.GetDefaultRefreshTime()
	if err != nil {
		log.Println("[ERROR]: could not get refresh time from database")
		refreshTime = 5
	}

	return RenderView(c, security_views.SecurityIndex("| Security", security_views.SecurityUpdates(c, p, *f, systemUpdates, availableOSes, availableUpdateStatus, refreshTime, commonInfo), commonInfo))
}

func (h *Handler) ListLatestUpdates(c echo.Context) error {
	var err error

	commonInfo, err := h.GetCommonInfo(c)
	if err != nil {
		return err
	}

	agentId := c.Param("uuid")
	if agentId == "" {
		return RenderError(c, partials.ErrorMessage(i18n.T(c.Request().Context(), "agents.no_empty_id"), false))
	}

	agent, err := h.Model.GetAgentById(agentId, commonInfo)
	if err != nil {
		return RenderError(c, partials.ErrorMessage(i18n.T(c.Request().Context(), "agents.could_not_get_agent"), false))
	}

	p := partials.NewPaginationAndSort()
	p.GetPaginationAndSortParams(c.FormValue("page"), c.FormValue("pageSize"), c.FormValue("sortBy"), c.FormValue("sortOrder"), c.FormValue("currentSortBy"))

	if p.SortBy == "" {
		p.SortBy = "name"
		p.SortOrder = "asc"
	}

	p.NItems, err = h.Model.CountLatestUpdates(agentId, commonInfo)
	if err != nil {
		return RenderError(c, partials.ErrorMessage(err.Error(), false))
	}

	updates, err := h.Model.GetLatestUpdates(agentId, p, commonInfo)
	if err != nil {
		return RenderView(c, security_views.SecurityIndex("| Security", partials.Error(c, err.Error(), "Security", partials.GetNavigationUrl(commonInfo, "/security"), commonInfo), commonInfo))
	}

	if c.Request().Method == "POST" {
		return RenderView(c, security_views.LatestUpdates(c, p, agent, updates, commonInfo))
	}

	return RenderView(c, security_views.SecurityIndex("| Security", security_views.LatestUpdates(c, p, agent, updates, commonInfo), commonInfo))
}

func (h *Handler) GetAntiviriFilters(c echo.Context) (*filters.AntivirusFilter, []string, []string, error) {
	// Get filters values
	f := filters.AntivirusFilter{}

	commonInfo, err := h.GetCommonInfo(c)
	if err != nil {
		return nil, nil, nil, err
	}

	f.Nickname = c.FormValue("filterByNickname")

	availableOSes, err := h.Model.GetAgentsUsedOSes(commonInfo)
	if err != nil {
		return nil, nil, nil, err
	}
	filteredAgentOSes := []string{}
	for index := range availableOSes {
		value := c.FormValue(fmt.Sprintf("filterByAgentOS%d", index))
		if value != "" {
			filteredAgentOSes = append(filteredAgentOSes, value)
		}
	}
	f.AgentOSVersions = filteredAgentOSes

	detectedAntiviri, err := h.Model.GetDetectedAntiviri(commonInfo)
	if err != nil {
		return nil, nil, nil, err
	}
	filteredAntiviri := []string{}
	for index := range detectedAntiviri {
		value := c.FormValue(fmt.Sprintf("filterByAntivirusName%d", index))
		if value != "" {
			filteredAntiviri = append(filteredAntiviri, value)
		}
	}
	f.AntivirusNameOptions = filteredAntiviri

	filteredEnableStatus := []string{}
	for index := range []string{"Enabled", "Disabled"} {
		value := c.FormValue(fmt.Sprintf("filterByAntivirusEnabled%d", index))
		if value != "" {
			filteredEnableStatus = append(filteredEnableStatus, value)
		}
	}
	f.AntivirusEnabledOptions = filteredEnableStatus

	filteredUpdateStatus := []string{}
	for index := range []string{"UpdatedYes", "UpdatedNo"} {
		value := c.FormValue(fmt.Sprintf("filterByAntivirusUpdated%d", index))
		if value != "" {
			filteredUpdateStatus = append(filteredUpdateStatus, value)
		}
	}
	f.AntivirusUpdatedOptions = filteredUpdateStatus

	return &f, availableOSes, detectedAntiviri, nil
}

func (h *Handler) GetSystemUpdatesFilters(c echo.Context) (*filters.SystemUpdatesFilter, []string, []string, error) {
	// Get filters values
	f := filters.SystemUpdatesFilter{}

	commonInfo, err := h.GetCommonInfo(c)
	if err != nil {
		return nil, nil, nil, err
	}

	f.Nickname = c.FormValue("filterByNickname")

	availableOSes, err := h.Model.GetAgentsUsedOSes(commonInfo)
	if err != nil {
		return nil, nil, nil, err
	}
	filteredAgentOSes := []string{}
	for index := range availableOSes {
		value := c.FormValue(fmt.Sprintf("filterByAgentOS%d", index))
		if value != "" {
			filteredAgentOSes = append(filteredAgentOSes, value)
		}
	}
	f.AgentOSVersions = filteredAgentOSes

	lastSearchFrom := c.FormValue("filterByLastSearchDateFrom")
	if lastSearchFrom != "" {
		f.LastSearchFrom = lastSearchFrom
	}
	lastSearchTo := c.FormValue("filterByLastSearchDateTo")
	if lastSearchTo != "" {
		f.LastSearchTo = lastSearchTo
	}

	lastInstallFrom := c.FormValue("filterByLastInstallDateFrom")
	if lastInstallFrom != "" {
		f.LastInstallFrom = lastInstallFrom
	}
	lastInstallTo := c.FormValue("filterByLastInstallDateTo")
	if lastInstallTo != "" {
		f.LastInstallTo = lastInstallTo
	}

	filteredPendingUpdates := []string{}
	for index := range []string{"Yes", "No"} {
		value := c.FormValue(fmt.Sprintf("filterByPendingUpdate%d", index))
		if value != "" {
			filteredPendingUpdates = append(filteredPendingUpdates, value)
		}
	}
	f.PendingUpdateOptions = filteredPendingUpdates

	availableUpdateStatus := scnorionplus_nats.SystemUpdatePossibleStatus()
	filteredUpdateStatus := []string{}
	for index := range availableUpdateStatus {
		value := c.FormValue(fmt.Sprintf("filterByUpdateStatus%d", index))
		if value != "" {
			filteredUpdateStatus = append(filteredUpdateStatus, value)
		}
	}
	f.UpdateStatus = filteredUpdateStatus

	return &f, availableOSes, availableUpdateStatus, err
}
