package validation

const (
	MinNameLength = 1
	MaxNameLength = 255
)

const (
	MinUnits = 1
	MaxUnits = 2_147_483_647
)

const (
	DefaultLimit = 100
	MaxLimit = 1000
)

const (
	PercentageMultiplier = 100.0
)

type ListParams struct {
	Limit  int
	Offset int
}

func (p ListParams) Normalize() ListParams {
	if p.Limit <= 0 {
		p.Limit = DefaultLimit
	}
	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	return p
}
