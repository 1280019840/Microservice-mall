package money

import (
	"errors"

	pb "checkoutservice/proto"
)

const (
	nanosMin = -999999999
	nanosMax = +999999999
	nanosMod = 1000000000
)

var (
	ErrInvalidValue        = errors.New("指定的货币是无效的")
	ErrMismatchingCurrency = errors.New("没有货币编码")
)

// 是否有效
func IsValid(m *pb.Money) bool {
	return signMatches(m) && validNanos(m.GetNanos())
}

func signMatches(m *pb.Money) bool {
	return m.GetNanos() == 0 || m.GetUnits() == 0 || (m.GetNanos() < 0) == (m.GetUnits() < 0)
}

func validNanos(nanos int32) bool { return nanosMin <= nanos && nanos <= nanosMax }

// 零
func IsZero(m *pb.Money) bool { return m.GetUnits() == 0 && m.GetNanos() == 0 }

// 正的
func IsPositive(m *pb.Money) bool {
	return IsValid(m) && m.GetUnits() > 0 || (m.GetUnits() == 0 && m.GetNanos() > 0)
}

// 负的
func IsNegative(m *pb.Money) bool {
	return IsValid(m) && m.GetUnits() < 0 || (m.GetUnits() == 0 && m.GetNanos() < 0)
}

// 是否相同
func AreSameCurrency(l, r *pb.Money) bool {
	return l.GetCurrencyCode() == r.GetCurrencyCode() && l.GetCurrencyCode() != ""
}

// 是否相等
func AreEquals(l, r *pb.Money) bool {
	return l.GetCurrencyCode() == r.GetCurrencyCode() &&
		l.GetUnits() == r.GetUnits() && l.GetNanos() == r.GetNanos()
}

// 变成负的
func Negate(m *pb.Money) *pb.Money {
	return &pb.Money{
		Units:        -m.GetUnits(),
		Nanos:        -m.GetNanos(),
		CurrencyCode: m.GetCurrencyCode()}
}

// Must
func Must(v *pb.Money, err error) *pb.Money {
	if err != nil {
		panic(err)
	}
	return v
}

// sum
func Sum(l, r *pb.Money) (*pb.Money, error) {
	if !IsValid(l) || !IsValid(r) {
		return &pb.Money{}, ErrInvalidValue
	} else if l.GetCurrencyCode() != r.GetCurrencyCode() {
		return &pb.Money{}, ErrMismatchingCurrency
	}
	units := l.GetUnits() + r.GetUnits()
	nanos := l.GetNanos() + r.GetNanos()

	if (units == 0 && nanos == 0) || (units >= 0 && nanos >= 0) || (units < 0 && nanos <= 0) {
		// 相同 sign <units, nanos>
		units += int64(nanos / nanosMod)
		nanos = nanos % nanosMod
	} else {
		// 不同 sign.
		if units > 0 {
			units--
			nanos += nanosMod
		} else {
			units++
			nanos -= nanosMod
		}
	}

	return &pb.Money{
		Units:        units,
		Nanos:        nanos,
		CurrencyCode: l.GetCurrencyCode()}, nil
}

func MultiplySlow(m *pb.Money, n uint32) *pb.Money {
	out := m
	for n > 1 {
		out = Must(Sum(out, m))
		n--
	}
	return out
}
