package schema

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func DecimalToNumeric(d decimal.Decimal) pgtype.Numeric {
	return pgtype.Numeric{
		Int:              d.BigInt(),
		Exp:              d.Exponent(),
		NaN:              false,
		InfinityModifier: pgtype.Finite,
		Valid:            true,
	}
}
