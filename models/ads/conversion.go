package ads

import (
	"context"
	"fmt"
	"time"

	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/xlog"
	"github.com/sirupsen/logrus"
)

// ConversionByID try to insert conversion based on ID
func ConversionByID(ctx context.Context, impid int64, acid string, t time.Time) error {
	var err error
	imp, err := FindImpressionByID(impid, t)
	if err == nil {
		return addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid)
	}
	imp, err = FindImpressionByID(impid, t.AddDate(0, 0, -1))
	if err == nil {
		return addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid)
	}
	imp, err = FindImpFromClickByImpID(impid)
	if err == nil {
		return addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid)
	}

	xlog.GetWithError(ctx, err).WithFields(logrus.Fields{
		"action_id": acid,
		"imp_id":    fmt.Sprint(impid),
	}).Error("failed to add conversion")

	return err
}

// ConversionByRH add conversion based on reserved hash
func ConversionByRH(ctx context.Context, rh string, acid string, t time.Time) error {
	var err error
	imp, err := FindImpressionByRH(rh, t)
	if err == nil {
		return addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid)
	}
	imp, err = FindImpressionByRH(rh, t.AddDate(0, 0, -1))
	if err == nil {
		return addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid)
	}
	imp, err = FindImpFromClickByRH(rh)
	if err == nil {
		return addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid)
	}

	xlog.GetWithError(ctx, err).WithFields(logrus.Fields{
		"action_id": acid,
		"imp_id":    rh,
	}).Error("failed to add conversion")
	return err

}

// addConversion insert to conversion table
func addConversion(wID, appID, wpID, caID, adID, copID, cpID, slotID, impID int64, acid string) error {
	return entities.AddConversion(wID, appID, wpID, caID, adID, copID, cpID, slotID, impID, acid)
}
