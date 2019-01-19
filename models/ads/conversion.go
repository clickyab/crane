package ads

import (
	"context"
	"time"

	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/xlog"
	"github.com/sirupsen/logrus"
)

// ConversionByRH add conversion based on reserved hash
func ConversionByRH(ctx context.Context, rh string, acid string, t time.Time, ip string) (err error) {

	defer func() {
		xlog.GetWithError(ctx, err).WithFields(logrus.Fields{
			"action_id": acid,
			"imp_id":    rh,
		}).Error("failed to add conversion")
	}()

	imp, err := FindImpressionByRH(rh, t)
	if err == nil {
		err = addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid, rh, ip)
		return
	}
	imp, err = FindImpressionByRH(rh, t.AddDate(0, 0, -1))
	if err == nil {
		err = addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid, rh, ip)
		return
	}
	imp, err = FindImpFromClickByRH(rh)
	if err == nil {
		err = addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid, rh, ip)
		return
	}

	return
}

// addConversion insert to conversion table
func addConversion(wID, appID, wpID, caID, adID, copID, cpID, slotID, impID int64, acid, rh, ip string) error {
	return entities.AddConversion(wID, appID, wpID, caID, adID, copID, cpID, slotID, impID, acid, rh, ip)
}
