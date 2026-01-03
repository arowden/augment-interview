package ownership

import "github.com/google/uuid"

// CapTableView is a read model representing a paginated view of a fund's cap table.
type CapTableView struct {
	FundID     uuid.UUID
	Entries    []*Entry
	TotalCount int
	Limit      int
	Offset     int
}

// TotalUnits calculates the sum of all units across entries in this view.
// Note: This only reflects the units in the current page of entries.
func (c *CapTableView) TotalUnits() int {
	total := 0
	for _, e := range c.Entries {
		total += e.Units
	}
	return total
}

// FindOwner searches for an owner by name in the current page of entries.
// Returns the entry if found, or nil if the owner is not in this page.
// Owner name matching is case-sensitive.
func (c *CapTableView) FindOwner(ownerName string) *Entry {
	for _, e := range c.Entries {
		if e.OwnerName == ownerName {
			return e
		}
	}
	return nil
}
