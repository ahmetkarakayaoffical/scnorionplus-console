package models

import (
	"context"
	"strconv"

	"entgo.io/ent/dialect/sql"
	ent "github.com/ahmetkarakayaoffical/ent"
	"github.com/ahmetkarakayaoffical/ent/agent"
	"github.com/ahmetkarakayaoffical/ent/antivirus"
	"github.com/ahmetkarakayaoffical/ent/site"
	"github.com/ahmetkarakayaoffical/ent/tenant"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/filters"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/partials"
)

type Antivirus struct {
	ID        string
	Nickname  string
	OS        string
	Name      string
	IsActive  bool `sql:"is_active"`
	IsUpdated bool `sql:"is_updated"`
	SiteID    int
}

func mainAntivirusQuery(s *sql.Selector, p partials.PaginationAndSort) {
	// Info from agents waiting for admission won't be shown
	s.Select(sql.As(agent.FieldID, "ID"), agent.FieldNickname, agent.FieldOs, antivirus.FieldName, antivirus.FieldIsActive, antivirus.FieldIsUpdated).
		LeftJoin(sql.Table(antivirus.Table)).
		On(agent.FieldID, antivirus.OwnerColumn).
		Where(sql.And(sql.NEQ(agent.FieldAgentStatus, agent.AgentStatusWaitingForAdmission)))

	if p.PageSize != 0 {
		s.Limit(p.PageSize).Offset((p.CurrentPage - 1) * p.PageSize)
	}
}

func (m *Model) CountAllAntiviri(f filters.AntivirusFilter, c *partials.CommonInfo) (int, error) {
	var query *ent.AgentQuery
	siteID, err := strconv.Atoi(c.SiteID)
	if err != nil {
		return 0, err
	}
	tenantID, err := strconv.Atoi(c.TenantID)
	if err != nil {
		return 0, err
	}

	// Info from agents waiting for admission won't be shown

	if siteID == -1 {
		query = m.Client.Agent.Query().
			Where(agent.AgentStatusNEQ(agent.AgentStatusWaitingForAdmission)).
			Where(agent.HasSiteWith(site.HasTenantWith(tenant.ID(tenantID))))

	} else {
		query = m.Client.Agent.Query().
			Where(agent.AgentStatusNEQ(agent.AgentStatusWaitingForAdmission)).
			Where(agent.HasSiteWith(site.ID(siteID), site.HasTenantWith(tenant.ID(tenantID))))
	}

	applyAntiviriFilters(query, f)

	return query.Count(context.Background())
}

func (m *Model) GetAntiviriByPage(p partials.PaginationAndSort, f filters.AntivirusFilter, c *partials.CommonInfo) ([]Antivirus, error) {
	var query *ent.AgentQuery
	var antiviri []Antivirus
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
		query = m.Client.Agent.Query().
			Where(agent.HasSiteWith(site.HasTenantWith(tenant.ID(tenantID))))
	} else {
		query = m.Client.Agent.Query().
			Where(agent.HasSiteWith(site.ID(siteID), site.HasTenantWith(tenant.ID(tenantID))))
	}

	applyAntiviriFilters(query, f)

	// Default sort
	if p.SortBy == "" {
		p.SortBy = "nickname"
		p.SortOrder = "desc"
	}

	switch p.SortBy {
	case "nickname":
		if p.SortOrder == "asc" {
			err = query.Modify(func(s *sql.Selector) {
				mainAntivirusQuery(s, p)
				s.OrderBy(sql.Asc(agent.FieldNickname))
			}).Scan(context.Background(), &antiviri)
		} else {
			err = query.Modify(func(s *sql.Selector) {
				mainAntivirusQuery(s, p)
				s.OrderBy(sql.Desc(agent.FieldNickname))
			}).Scan(context.Background(), &antiviri)
		}
	case "agentOS":
		if p.SortOrder == "asc" {
			err = query.Modify(func(s *sql.Selector) {
				mainAntivirusQuery(s, p)
				s.OrderBy(sql.Asc(agent.FieldOs))
			}).Scan(context.Background(), &antiviri)
		} else {
			err = query.Modify(func(s *sql.Selector) {
				mainAntivirusQuery(s, p)
				s.OrderBy(sql.Desc(agent.FieldOs))
			}).Scan(context.Background(), &antiviri)
		}
	case "antivirusName":
		if p.SortOrder == "asc" {
			err = query.Modify(func(s *sql.Selector) {
				mainAntivirusQuery(s, p)
				s.OrderBy(sql.Asc(antivirus.FieldName))
			}).Scan(context.Background(), &antiviri)
		} else {
			err = query.Modify(func(s *sql.Selector) {
				mainAntivirusQuery(s, p)
				s.OrderBy(sql.Desc(antivirus.FieldName))
			}).Scan(context.Background(), &antiviri)
		}
	case "antivirusEnabled":
		if p.SortOrder == "asc" {
			err = query.Modify(func(s *sql.Selector) {
				mainAntivirusQuery(s, p)
				s.OrderBy(sql.Asc(antivirus.FieldIsActive))
			}).Scan(context.Background(), &antiviri)
		} else {
			err = query.Modify(func(s *sql.Selector) {
				mainAntivirusQuery(s, p)
				s.OrderBy(sql.Desc(antivirus.FieldIsActive))
			}).Scan(context.Background(), &antiviri)
		}
	case "antivirusUpdated":
		if p.SortOrder == "asc" {
			err = query.Modify(func(s *sql.Selector) {
				mainAntivirusQuery(s, p)
				s.OrderBy(sql.Asc(antivirus.FieldIsUpdated))
			}).Scan(context.Background(), &antiviri)
		} else {
			err = query.Modify(func(s *sql.Selector) {
				mainAntivirusQuery(s, p)
				s.OrderBy(sql.Desc(antivirus.FieldIsUpdated))
			}).Scan(context.Background(), &antiviri)
		}
	}

	if err != nil {
		return nil, err
	}

	// Add site ids
	sortedAgentIDs := []string{}
	for _, computer := range antiviri {
		sortedAgentIDs = append(sortedAgentIDs, computer.ID)
	}
	agents, err := m.Client.Agent.Query().WithSite().Where(agent.IDIn(sortedAgentIDs...)).All(context.Background())
	if err != nil {
		return nil, err
	}

	// Add site id to each computer in order
	for i, computer := range antiviri {
		for _, agent := range agents {
			if computer.ID == agent.ID {
				if len(agent.Edges.Site) == 1 {
					antiviri[i].SiteID = agent.Edges.Site[0].ID
				} else {
					antiviri[i].SiteID = -1
				}
				break
			}
		}
	}

	return antiviri, nil
}

func (m *Model) GetDetectedAntiviri(c *partials.CommonInfo) ([]string, error) {
	siteID, err := strconv.Atoi(c.SiteID)
	if err != nil {
		return nil, err
	}
	tenantID, err := strconv.Atoi(c.TenantID)
	if err != nil {
		return nil, err
	}

	if siteID == -1 {
		return m.Client.Antivirus.Query().Unique(true).
			Where(antivirus.HasOwnerWith(agent.AgentStatusNEQ(agent.AgentStatusWaitingForAdmission),
				agent.HasSiteWith(site.HasTenantWith(tenant.ID(tenantID))))).
			Select(antivirus.FieldName).Strings(context.Background())
	} else {
		return m.Client.Antivirus.Query().Unique(true).
			Where(antivirus.HasOwnerWith(agent.AgentStatusNEQ(agent.AgentStatusWaitingForAdmission),
				agent.HasSiteWith(site.ID(siteID), site.HasTenantWith(tenant.ID(tenantID))))).
			Select(antivirus.FieldName).Strings(context.Background())
	}
}

func applyAntiviriFilters(query *ent.AgentQuery, f filters.AntivirusFilter) {
	if len(f.Nickname) > 0 {
		query.Where(agent.NicknameContainsFold(f.Nickname))
	}

	if len(f.AgentOSVersions) > 0 {
		query.Where(agent.OsIn(f.AgentOSVersions...))
	}

	if len(f.AntivirusNameOptions) > 0 {
		query.Where(agent.HasAntivirusWith(antivirus.NameIn(f.AntivirusNameOptions...)))
	}

	if len(f.AntivirusEnabledOptions) > 0 {
		if len(f.AntivirusEnabledOptions) == 1 && f.AntivirusEnabledOptions[0] == "Enabled" {
			query.Where(agent.HasAntivirusWith(antivirus.IsActive(true)))
		}

		if len(f.AntivirusEnabledOptions) == 1 && f.AntivirusEnabledOptions[0] == "Disabled" {
			query.Where(agent.HasAntivirusWith(antivirus.IsActive(false)))
		}
	}

	if len(f.AntivirusUpdatedOptions) > 0 {
		if len(f.AntivirusUpdatedOptions) == 1 && f.AntivirusUpdatedOptions[0] == "UpdatedYes" {
			query.Where(agent.HasAntivirusWith(antivirus.IsUpdated(true)))
		}

		if len(f.AntivirusUpdatedOptions) == 1 && f.AntivirusUpdatedOptions[0] == "UpdatedNo" {
			query.Where(agent.HasAntivirusWith(antivirus.IsUpdated(false)))
		}
	}
}
