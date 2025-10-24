package models

import (
	"context"
	"strconv"
	"time"

	"github.com/ahmetkarakayaoffical/ent/certificate"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/filters"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/partials"
	scnorionplus_ent "github.com/open-uem/ent"
)

func (m *Model) GetCertificateByUID(uid string) (*scnorionplus_ent.Certificate, error) {
	return m.Client.Certificate.Query().Where(certificate.UID(uid)).Only(context.Background())
}

func (m *Model) GetCertificateBySerial(serial string) (*scnorionplus_ent.Certificate, error) {
	serialNumber, err := strconv.ParseInt(serial, 10, 64)
	if err != nil {
		return nil, err
	}

	return m.Client.Certificate.Query().Where(certificate.ID(serialNumber)).Only(context.Background())
}

func (m *Model) RevokeCertificate(cert *scnorionplus_ent.Certificate, info string, reason int) error {
	return m.Client.Revocation.Create().SetID(cert.ID).SetExpiry(cert.Expiry).SetRevoked(time.Now()).SetReason(reason).SetInfo(info).Exec(context.Background())
}

func (m *Model) DeleteCertificate(serial int64) error {
	return m.Client.Certificate.DeleteOneID(serial).Exec(context.Background())
}

func (m *Model) CountAllCertificates(f filters.CertificateFilter) (int, error) {
	query := m.Client.Certificate.Query()

	// Apply filters
	applyCertificateFilters(query, f)

	count, err := query.Count(context.Background())
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *Model) CountCertificatesAboutToexpire() (int, error) {
	// Certificates that expires in two months
	return m.Client.Certificate.Query().Where(certificate.ExpiryLT(time.Now().AddDate(0, 2, 0))).Count(context.Background())
}

func (m *Model) GetCertificatesByPage(p partials.PaginationAndSort, f filters.CertificateFilter) ([]*scnorionplus_ent.Certificate, error) {
	query := m.Client.Certificate.Query().Limit(p.PageSize).Offset((p.CurrentPage - 1) * p.PageSize)

	// Apply filters
	applyCertificateFilters(query, f)

	switch p.SortBy {
	case "serial":
		if p.SortOrder == "asc" {
			query = query.Order(scnorionplus_ent.Asc(certificate.FieldID))
		} else {
			query = query.Order(scnorionplus_ent.Desc(certificate.FieldID))
		}
	case "type":
		if p.SortOrder == "asc" {
			query = query.Order(scnorionplus_ent.Asc(certificate.FieldType))
		} else {
			query = query.Order(scnorionplus_ent.Desc(certificate.FieldType))
		}
	case "description":
		if p.SortOrder == "asc" {
			query = query.Order(scnorionplus_ent.Asc(certificate.FieldDescription))
		} else {
			query = query.Order(scnorionplus_ent.Desc(certificate.FieldDescription))
		}
	case "expiry":
		if p.SortOrder == "asc" {
			query = query.Order(scnorionplus_ent.Asc(certificate.FieldExpiry))
		} else {
			query = query.Order(scnorionplus_ent.Desc(certificate.FieldExpiry))
		}
	case "username":
		if p.SortOrder == "asc" {
			query = query.Order(scnorionplus_ent.Asc(certificate.FieldUID))
		} else {
			query = query.Order(scnorionplus_ent.Desc(certificate.FieldUID))
		}
	default:
		query = query.Order(scnorionplus_ent.Desc(certificate.FieldID))
	}

	return query.All(context.Background())
}

func (m *Model) GetCertificatesTypes() ([]string, error) {
	return m.Client.Certificate.Query().Unique(true).Select(certificate.FieldType).Strings(context.Background())
}

func applyCertificateFilters(query *scnorionplus_ent.CertificateQuery, f filters.CertificateFilter) {

	if len(f.TypeOptions) > 0 {
		selectedTypes := []certificate.Type{}
		for _, option := range f.TypeOptions {
			selectedTypes = append(selectedTypes, certificate.Type(option))
		}

		query.Where(certificate.TypeIn(selectedTypes...))
	}

	if len(f.Description) > 0 {
		query.Where(certificate.DescriptionContainsFold(f.Description))
	}

	if len(f.ExpiryFrom) > 0 {
		from, err := time.Parse("2006-01-02", f.ExpiryFrom)
		if err == nil {
			query.Where(certificate.ExpiryGTE(from))
		}
	}

	if len(f.ExpiryTo) > 0 {
		to, err := time.Parse("2006-01-02", f.ExpiryTo)
		if err == nil {
			query.Where(certificate.ExpiryLTE(to))
		}
	}

	if len(f.Username) > 0 {
		query.Where(certificate.UIDContainsFold(f.Username))
	}
}
