package safemap

import (
	"server/internal/model"
	"sync"
)

type TierSafeMap struct {
	sync.Map
}

func (t *TierSafeMap) SafeStore(value *model.TierList) {
	t.Store((string(value.Tier)), value)
}

func (t *TierSafeMap) SafeLoad(key string) (*model.TierList, bool) {
	v, ok := t.Load(key)
	if !ok {
		return nil, ok
	}
	rank, ok := v.(*model.TierList)
	if !ok {
		return nil, ok
	}
	return rank, ok
}

func (t *TierSafeMap) SafeDelete(key string) {
	t.Delete(key)
}

func (t *TierSafeMap) ToList() []*model.TierList {
	var list []*model.TierList
	t.Range(func(key, value interface{}) bool {
		tier, ok := value.(*model.TierList)
		if !ok {
			return false
		}
		list = append(list, tier)
		return true
	})

	return list
}
