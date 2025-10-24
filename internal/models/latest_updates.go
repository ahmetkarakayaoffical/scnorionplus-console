package models

import (
	"context"
	"strconv"

	ent "github.com/ahmetkarakayaoffical/ent"
	"github.com/ahmetkarakayaoffical/ent/agent"
	"github.com/ahmetkarakayaoffical/ent/site"
	"github.com/ahmetkarakayaoffical/ent/tenant"
	"github.com/ahmetkarakayaoffical/ent/update"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/partials"
)

func (m *Model) CountLatestUpdates(agentId string, c *partials.CommonInfo) (int, error) {
	siteID, err := strconv.Atoi(c.SiteID)
	if err != nil {
		return 0, err
	}
	tenantID, err := strconv.Atoi(c.TenantID)
	if err != nil {
		return 0, err
	}

	if siteID == -1 {
		return m.Client.Update.Query().Where(update.HasOwnerWith(agent.ID(agentId), agent.HasSiteWith(site.HasTenantWith(tenant.ID(tenantID))))).Count(context.Background())
	} else {
		return m.Client.Update.Query().Where(update.HasOwnerWith(agent.ID(agentId), agent.HasSiteWith(site.ID(siteID), site.HasTenantWith(tenant.ID(tenantID))))).Count(context.Background())
	}
}

func (m *Model) GetLatestUpdates(agentId string, p partials.PaginationAndSort, c *partials.CommonInfo) ([]*ent.Update, error) {
	var query *ent.UpdateQuery

	siteID, err := strconv.Atoi(c.SiteID)
	if err != nil {
		return nil, err
	}
	tenantID, err := strconv.Atoi(c.TenantID)
	if err != nil {
		return nil, err
	}

	if siteID == -1 {
		query = m.Client.Update.Query().Where(update.HasOwnerWith(agent.ID(agentId), agent.HasSiteWith(site.HasTenantWith(tenant.ID(tenantID)))))
	} else {
		query = m.Client.Update.Query().Where(update.HasOwnerWith(agent.ID(agentId), agent.HasSiteWith(site.ID(siteID), site.HasTenantWith(tenant.ID(tenantID)))))
	}

	switch p.SortBy {
	case "title":
		if p.SortOrder == "asc" {
			query = query.Order(ent.Asc(update.FieldTitle))
		} else {
			query = query.Order(ent.Desc(update.FieldTitle))
		}
	case "date":
		if p.SortOrder == "asc" {
			query = query.Order(ent.Asc(update.FieldDate))
		} else {
			query = query.Order(ent.Desc(update.FieldDate))
		}
	default:
		query = query.Order(ent.Desc(update.FieldDate))
	}

	updates, err := query.Limit(p.PageSize).Offset((p.CurrentPage - 1) * p.PageSize).All(context.Background())
	if err != nil {
		return nil, err
	}

	return updates, nil
}
