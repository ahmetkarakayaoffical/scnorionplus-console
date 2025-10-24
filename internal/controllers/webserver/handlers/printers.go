package handlers

import (
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/printers_views"
	"github.com/labstack/echo/v4"
)

func (h *Handler) NetworkPrinters(c echo.Context) error {
	commonInfo, err := h.GetCommonInfo(c)
	if err != nil {
		return err
	}

	return RenderView(c, printers_views.PrintersIndex("| Network Printers", printers_views.Printers(c, commonInfo), commonInfo))
}
