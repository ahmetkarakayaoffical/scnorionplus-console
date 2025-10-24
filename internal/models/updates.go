package models

import (
	"context"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	ent "github.com/ahmetkarakayaoffical/ent"
	"github.com/ahmetkarakayaoffical/ent/agent"
	"github.com/ahmetkarakayaoffical/ent/site"
	"github.com/ahmetkarakayaoffical/ent/systemupdate"
	"github.com/ahmetkarakayaoffical/ent/tenant"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/filters"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/partials"
)

type SystemUpdate struct {
	ID                 string
	Nickname           string
	OS                 string
	SystemUpdateStatus string    `sql:"system_update_status"`
	LastInstall        time.Time `sql:"last_install"`
	LastSearch         time.Time `sql:"last_search"`
	PendingUpdates     bool      `sql:"pending_updates"`
	SiteID             int
}

func mainUpdatesQuery(s *sql.Selector, p partials.PaginationAndSort) {
	// Info from agents waiting for admission won't be shown
	s.Select(sql.As(agent.FieldID, "ID"), agent.FieldNickname, agent.FieldOs, systemupdate.FieldSystemUpdateStatus, systemupdate.FieldLastInstall, systemupdate.FieldLastSearch, systemupdate.FieldPendingUpdates).
		LeftJoin(sql.Table(systemupdate.Table)).
		On(agent.FieldID, systemupdate.OwnerColumn).
		Where(sql.And(sql.NEQ(agent.FieldAgentStatus, agent.AgentStatusWaitingForAdmission)))
	if p.PageSize != 0 {
		s.Limit(p.PageSize).Offset((p.CurrentPage - 1) * p.PageSize)
	}
}

func (m *Model) CountAllSystemUpdates(f filters.SystemUpdatesFilter, c *partials.CommonInfo) (int, error) {
	var query *ent.AgentQuery

	siteID, err := strconv.Atoi(c.SiteID)
	if err != nil {
		return 0, err
	}
	tenantID, err := strconv.Atoi(c.TenantID)
	if err != nil {
		return 0, err
	}

	if siteID == -1 {
		query = m.Client.Agent.Query().
			Where(agent.AgentStatusNEQ(agent.AgentStatusWaitingForAdmission)).
			Where(agent.HasSiteWith(site.HasTenantWith(tenant.ID(tenantID))))
	} else {
		query = m.Client.Agent.Query().
			Where(agent.AgentStatusNEQ(agent.AgentStatusWaitingForAdmission)).
			Where(agent.HasSiteWith(site.ID(siteID), site.HasTenantWith(tenant.ID(tenantID))))
	}

	applySystemUpdatesFilters(query, f)

	return query.Count(context.Background())
}

func (m *Model) GetSystemUpdatesByPage(p partials.PaginationAndSort, f filters.SystemUpdatesFilter, c *partials.CommonInfo) ([]SystemUpdate, error) {
	var query *ent.AgentQuery
	var systemUpdates []SystemUpdate
	var err error

	siteID, err := strconv.Atoi(c.SiteID)
	if err != nil {
		return nil, err
	}
	tenantID, err := strconv.Atoi(c.TenantID)
	if err != nil {
		return nil, err
	}

	if siteID == -1 {
		query = m.Client.Agent.Query().Where(agent.HasSiteWith(site.HasTenantWith(tenant.ID(tenantID))))
	} else {
		query = m.Client.Agent.Query().Where(agent.HasSiteWith(site.ID(siteID), site.HasTenantWith(tenant.ID(tenantID))))
	}

	applySystemUpdatesFilters(query, f)

	// Default sort
	if p.SortBy == "" {
		p.SortBy = "nickname"
		p.SortOrder = "desc"
	}

	switch p.SortBy {
	case "nickname":
		if p.SortOrder == "asc" {
			err = query.Modify(func(s *sql.Selector) {
				mainUpdatesQuery(s, p)
				s.OrderBy(sql.Asc(agent.FieldNickname))
			}).Scan(context.Background(), &systemUpdates)
		} else {
			err = query.Modify(func(s *sql.Selector) {
				mainUpdatesQuery(s, p)
				s.OrderBy(sql.Desc(agent.FieldNickname))
			}).Scan(context.Background(), &systemUpdates)
		}
	case "agentOS":
		if p.SortOrder == "asc" {
			err = query.Modify(func(s *sql.Selector) {
				mainUpdatesQuery(s, p)
				s.OrderBy(sql.Asc(agent.FieldOs))
			}).Scan(context.Background(), &systemUpdates)
		} else {
			err = query.Modify(func(s *sql.Selector) {
				mainUpdatesQuery(s, p)
				s.OrderBy(sql.Desc(agent.FieldOs))
			}).Scan(context.Background(), &systemUpdates)
		}
	case "updateStatus":
		if p.SortOrder == "asc" {
			err = query.Modify(func(s *sql.Selector) {
				mainUpdatesQuery(s, p)
				s.OrderBy(sql.Asc(systemupdate.FieldSystemUpdateStatus))
			}).Scan(context.Background(), &systemUpdates)
		} else {
			err = query.Modify(func(s *sql.Selector) {
				mainUpdatesQuery(s, p)
				s.OrderBy(sql.Desc(systemupdate.FieldSystemUpdateStatus))
			}).Scan(context.Background(), &systemUpdates)
		}
	case "lastSearch":
		if p.SortOrder == "asc" {
			err = query.Modify(func(s *sql.Selector) {
				mainUpdatesQuery(s, p)
				s.OrderBy(sql.Asc(systemupdate.FieldLastSearch))
			}).Scan(context.Background(), &systemUpdates)
		} else {
			err = query.Modify(func(s *sql.Selector) {
				mainUpdatesQuery(s, p)
				s.OrderBy(sql.Desc(systemupdate.FieldLastSearch))
			}).Scan(context.Background(), &systemUpdates)
		}
	case "lastInstall":
		if p.SortOrder == "asc" {
			err = query.Modify(func(s *sql.Selector) {
				mainUpdatesQuery(s, p)
				s.OrderBy(sql.Asc(systemupdate.FieldLastInstall))
			}).Scan(context.Background(), &systemUpdates)
		} else {
			err = query.Modify(func(s *sql.Selector) {
				mainUpdatesQuery(s, p)
				s.OrderBy(sql.Desc(systemupdate.FieldLastInstall))
			}).Scan(context.Background(), &systemUpdates)
		}
	case "pendingUpdates":
		if p.SortOrder == "asc" {
			err = query.Modify(func(s *sql.Selector) {
				mainUpdatesQuery(s, p)
				s.OrderBy(sql.Asc(systemupdate.FieldPendingUpdates))
			}).Scan(context.Background(), &systemUpdates)
		} else {
			err = query.Modify(func(s *sql.Selector) {
				mainUpdatesQuery(s, p)
				s.OrderBy(sql.Desc(systemupdate.FieldPendingUpdates))
			}).Scan(context.Background(), &systemUpdates)
		}
	}

	if err != nil {
		return nil, err
	}

	// Add site ids
	sortedAgentIDs := []string{}
	for _, computer := range systemUpdates {
		sortedAgentIDs = append(sortedAgentIDs, computer.ID)
	}
	agents, err := m.Client.Agent.Query().WithSite().Where(agent.IDIn(sortedAgentIDs...)).All(context.Background())
	if err != nil {
		return nil, err
	}

	// Add site id to each computer in order
	for i, computer := range systemUpdates {
		for _, agent := range agents {
			if computer.ID == agent.ID {
				if len(agent.Edges.Site) == 1 {
					systemUpdates[i].SiteID = agent.Edges.Site[0].ID
				} else {
					systemUpdates[i].SiteID = -1
				}
				break
			}
		}
	}

	return systemUpdates, nil
}

func applySystemUpdatesFilters(query *ent.AgentQuery, f filters.SystemUpdatesFilter) {
	if len(f.Nickname) > 0 {
		query.Where(agent.NicknameContainsFold(f.Nickname))
	}

	if len(f.AgentOSVersions) > 0 {
		query.Where(agent.OsIn(f.AgentOSVersions...))
	}

	if len(f.UpdateStatus) > 0 {
		query.Where(agent.HasSystemupdateWith(systemupdate.SystemUpdateStatusIn(f.UpdateStatus...)))
	}

	if len(f.LastSearchFrom) > 0 {
		from, err := time.Parse("2006-01-02", f.LastSearchFrom)
		if err == nil {
			query.Where(agent.HasSystemupdateWith(systemupdate.LastSearchGTE(from)))
		}
	}

	if len(f.LastSearchTo) > 0 {
		to, err := time.Parse("2006-01-02", f.LastSearchTo)
		if err == nil {
			query.Where(agent.HasSystemupdateWith(systemupdate.LastSearchLTE(to)))
		}
	}

	if len(f.LastInstallFrom) > 0 {
		from, err := time.Parse("2006-01-02", f.LastInstallFrom)
		if err == nil {
			query.Where(agent.HasSystemupdateWith(systemupdate.LastInstallGTE(from)))
		}
	}

	if len(f.LastInstallTo) > 0 {
		to, err := time.Parse("2006-01-02", f.LastInstallTo)
		if err == nil {
			query.Where(agent.HasSystemupdateWith(systemupdate.LastInstallLTE(to)))
		}
	}

	if len(f.PendingUpdateOptions) > 0 {
		if len(f.PendingUpdateOptions) == 1 && f.PendingUpdateOptions[0] == "Yes" {
			query.Where(agent.HasSystemupdateWith(systemupdate.PendingUpdates(true)))
		}

		if len(f.PendingUpdateOptions) == 1 && f.PendingUpdateOptions[0] == "No" {
			query.Where(agent.HasSystemupdateWith(systemupdate.PendingUpdates(false)))
		}
	}
}
