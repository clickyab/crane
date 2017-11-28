package cyslot

import (
	"errors"
	"strconv"

	"fmt"

	"clickyab.com/gad/models"
	"github.com/clickyab/services/assert"
)

// GetWebSlotID return slot id
func GetWebSlotID(pubID string, wID int64, size int) (int64, error) {
	m := models.NewManager()
	fSlots, err := m.FetchWebAppSlots(pubID, wID, 0)
	if err == nil {
		return fSlots.ID, nil
	}
	//slot not found
	pubInt, err := strconv.ParseInt(pubID, 10, 64)
	if err != nil {
		return 0, errors.New("pub id not integer")
	}
	rSlot, err := m.InsertSlots(wID, 0, pubInt, size)
	if err != nil {
		return 0, errors.New("cant insert web slot")
	}
	return rSlot.ID, nil
}

// GetCommonSlotIDs return slot ids (native and vast)
func GetCommonSlotIDs(pubIDSize map[string]int, wID int64) (map[string]int64, error) {
	var answer = make(map[string]int64)
	var newSlots = make(map[string]int)
	var newSlotsIDs []int64
	var pubIDs []string
	m := models.NewManager()
	for i := range pubIDSize {
		pubIDs = append(pubIDs, i)
	}
	foundSlots := m.FetchNativeSlots(pubIDs, wID)
	for i := range pubIDSize {
		if _, ok := answer[i]; ok {
			continue
		}
		for j := range foundSlots {
			if fmt.Sprintf("%d", foundSlots[j].PublicID) == i {
				answer[i] = foundSlots[j].ID
				break
			}
		}
		if _, ok := answer[i]; !ok {
			s, err := strconv.ParseInt(i, 10, 0)
			if err == nil {
				newSlots[i] = pubIDSize[i]
				newSlotsIDs = append(newSlotsIDs, s)
			}
		}
	}
	if len(newSlots) > 0 {
		newInsertedSlot := InsertNewSlots(wID, newSlotsIDs, func() []int {
			var x []int
			for k := range newSlots {
				x = append(x, newSlots[k])

			}
			return x
		}(),
		)
		for i := range newInsertedSlot {
			answer[i] = newInsertedSlot[i]
		}
	}
	return answer, nil
}

// GetWebSlotID return slot id
func GetAppSlotID(pubID string, appID int64, size int) (int64, error) {
	m := models.NewManager()
	fSlots, err := m.FetchWebAppSlots(pubID, 0, appID)
	if err == nil {
		return fSlots.ID, nil
	}
	//slot not found
	pubInt, err := strconv.ParseInt(pubID, 10, 64)
	if err != nil {
		return 0, errors.New("pub id not integer")
	}
	rSlot, err := m.InsertSlots(0, appID, pubInt, size)
	if err != nil {
		return 0, errors.New("cant insert app slot")
	}
	return rSlot.ID, nil
}

func InsertNewSlots(wID int64, newSlots []int64, newSize []int) map[string]int64 {
	assert.True(len(newSlots) == len(newSize), "[BUG] slot public and count is not matched")
	result := make(map[string]int64)
	if len(newSlots) > 0 {
		for i := range newSlots {
			insertedSlots, err := models.NewManager().InsertSlots(wID, 0, newSlots[i], newSize[i])
			if err == nil {
				p := fmt.Sprintf("%d", insertedSlots.PublicID)
				result[p] = insertedSlots.ID
			}
		}
	}

	return result
}
