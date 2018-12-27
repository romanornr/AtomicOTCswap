package bcoins

import (
	"errors"
	"math"
	"strconv"
)

// AmountUnit describes a method of converting an Amount to something
// other than the base unit of a viacoin.  The value of the AmountUnit
// is the exponent component of the decadic multiple to convert from
// an amount in viacoin to an amount counted in units.
type AmountUnit int

// These constants define various units used when describing a viacoin
// monetary amount.
const (
	AmountMegaBTC  AmountUnit = 6
	AmountKiloBTC  AmountUnit = 3
	AmountBTC      AmountUnit = 0
	AmountMilliBTC AmountUnit = -3
	AmountMicroBTC AmountUnit = -6
	AmountSatoshi  AmountUnit = -8

	// SatoshiPerBitcent is the number of satoshi in one bitcoin cent.
	SatoshiPerBitcent = 1e6

	// SatoshiPerBitcoin is the number of satoshi in one bitcoin (1 BTC).
	SatoshiPerBitcoin = 1e8

	// MaxSatoshi is the maximum transaction amount allowed in satoshi.
	MaxSatoshi = 21e6 * SatoshiPerBitcoin

)

// String returns the unit as a string.  For recognized units, the SI
// prefix is used, or "Satoshi" for the base unit.  For all unrecognized
// units, "1eN BTC" is returned, where N is the AmountUnit.
func (u AmountUnit) String() string {
	switch u {
	case AmountMegaBTC:
		return "MVIA"
	case AmountKiloBTC:
		return "kVIA"
	case AmountBTC:
		return "VIA"
	case AmountMilliBTC:
		return "mVIA"
	case AmountMicroBTC:
		return "μVIA"
	case AmountSatoshi:
		return "Satoshi"
	default:
		return "1e" + strconv.FormatInt(int64(u), 10) + " VIA"
	}
}

// Amount represents the base viacoin monetary unit (colloquially referred
// to as a `Satoshi').  A single Amount is equal to 1e-8 of a viacoin.
type Amount int64

// round converts a floating point number, which may or may not be representable
// as an integer, to the Amount integer type by rounding to the nearest integer.
// This is performed by adding or subtracting 0.5 depending on the sign, and
// relying on integer truncation to round the value to the nearest Amount.
func round(f float64) Amount {
	if f < 0 {
		return Amount(f - 0.5)
	}
	return Amount(f + 0.5)
}

// NewAmount creates an Amount from a floating point value representing
// some value in viacoin.  NewAmount errors if f is NaN or +-Infinity, but
// does not check that the amount is within the total amount of viacoin
// producible as f may not refer to an amount at a single moment in time.
//
// NewAmount is for specifically for converting BTC to Satoshi.
// For creating a new Amount with an int64 value which denotes a quantity of Satoshi,
// do a simple type conversion from type int64 to Amount.
// See GoDoc for example: http://godoc.org/github.com/viacoin/viautil#example-Amount
func NewAmount(f float64) (Amount, error) {
	// The amount is only considered invalid if it cannot be represented
	// as an integer type.  This may happen if f is NaN or +-Infinity.
	switch {
	case math.IsNaN(f):
		fallthrough
	case math.IsInf(f, 1):
		fallthrough
	case math.IsInf(f, -1):
		return 0, errors.New("invalid viacoin amount")
	}

	return round(f * SatoshiPerBitcoin), nil
}

// ToUnit converts a monetary amount counted in viacoin base units to a
// floating point value representing an amount of viacoin.
func (a Amount) ToUnit(u AmountUnit) float64 {
	return float64(a) / math.Pow10(int(u+8))
}

// ToBTC is the equivalent of calling ToUnit with AmountBTC.
func (a Amount) ToBTC() float64 {
	return a.ToUnit(AmountBTC)
}

// Format formats a monetary amount counted in viacoin base units as a
// string for a given unit.  The conversion will succeed for any unit,
// however, known units will be formated with an appended label describing
// the units with SI notation, or "Satoshi" for the base unit.
func (a Amount) Format(u AmountUnit) string {
	units := " " + u.String()
	return strconv.FormatFloat(a.ToUnit(u), 'f', -int(u+8), 64) + units
}

// String is the equivalent of calling Format with AmountBTC.
func (a Amount) String() string {
	return a.Format(AmountBTC)
}

// MulF64 multiplies an Amount by a floating point value.  While this is not
// an operation that must typically be done by a full node or wallet, it is
// useful for services that build on top of viacoin (for example, calculating
// a fee by multiplying by a percentage).
func (a Amount) MulF64(f float64) Amount {
	return round(float64(a) * f)
}