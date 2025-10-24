package models

import (
	"context"
	"strconv"

	"github.com/ahmetkarakayaoffical/ent/agent"
	"github.com/ahmetkarakayaoffical/ent/printer"
	"github.com/ahmetkarakayaoffical/ent/site"
	"github.com/ahmetkarakayaoffical/ent/tenant"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/partials"
)

func (m *Model) CountDifferentPrinters(c *partials.CommonInfo) (int, error) {
	siteID, err := strconv.Atoi(c.SiteID)
	if err != nil {
		return 0, err
	}
	tenantID, err := strconv.Atoi(c.TenantID)
	if err != nil {
		return 0, err
	}

	if siteID == -1 {
		return m.Client.Printer.Query().Where(printer.HasOwnerWith(agent.HasSiteWith(site.HasTenantWith(tenant.ID(tenantID))))).Select(printer.FieldName).Unique(true).Count(context.Background())
	} else {
		return m.Client.Printer.Query().Where(printer.HasOwnerWith(agent.HasSiteWith(site.ID(siteID), site.HasTenantWith(tenant.ID(tenantID))))).Select(printer.FieldName).Unique(true).Count(context.Background())
	}
}
