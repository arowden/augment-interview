package ownership

import "github.com/google/uuid"

type CapTableView struct {
	FundID     uuid.UUID
	Entries    []*Entry
	TotalCount int
	Limit      int
	Offset     int
}

func (c *CapTableView) TotalUnits() int {
	total := 0
	for _, e := range c.Entries {
		total += e.Units
	}
	return total
}

func (c *CapTableView) FindOwner(ownerName string) *Entry {
	for _, e := range c.Entries {
		if e.OwnerName == ownerName {
			return e
		}
	}
	return nil
}
