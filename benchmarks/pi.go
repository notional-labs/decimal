package benchmarks

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/dnum"
	apd "github.com/cockroachdb/apd/v3"
	"github.com/ericlagergren/decimal"
	"github.com/stretchr/testify/require"

	ssdec "github.com/shopspring/decimal"
	"gopkg.in/inf.v0"
)

func adjustPrecision(prec int) int {
	return int(math.Ceil(float64(prec) * 1.1))
}

var (
	infEight     = inf.NewDec(8, 0)
	infThirtyTwo = inf.NewDec(32, 0)
)

// PiInf calculates π to the desired precision using gopkg.in/inf.v0.
func PiInf(prec int) *inf.Dec {
	var (
		lasts = inf.NewDec(0, 0)
		t     = inf.NewDec(3, 0)
		s     = inf.NewDec(3, 0)
		n     = inf.NewDec(1, 0)
		na    = inf.NewDec(0, 0)
		d     = inf.NewDec(0, 0)
		da    = inf.NewDec(24, 0)

		work = adjustPrecision(prec)
	)

	for s.Cmp(lasts) != 0 {
		lasts.Set(s)
		n.Add(n, na)
		na.Add(na, infEight)
		d.Add(d, da)
		da.Add(da, infThirtyTwo)
		t.Mul(t, n)
		t.QuoRound(t, d, inf.Scale(work), inf.RoundHalfUp)
		s.Add(s, t)
	}
	// -1 because inf's precision == digits after radix
	return s.Round(s, inf.Scale(prec-1), inf.RoundHalfUp)
}

var (
	ssdecEight     = ssdec.New(8, 0)
	ssdecThirtyTwo = ssdec.New(32, 0)
)

// PiShopSpring calculates π to the desired precision using
// github.com/shopspring/decimal.
func PiShopSpring(prec int32) ssdec.Decimal {
	var (
		lasts = ssdec.New(0, 0)
		t     = ssdec.New(3, 0)
		s     = ssdec.New(3, 0)
		n     = ssdec.New(1, 0)
		na    = ssdec.New(0, 0)
		d     = ssdec.New(0, 0)
		da    = ssdec.New(24, 0)

		work = int32(adjustPrecision(int(prec)))
	)

	for s.Cmp(lasts) != 0 {
		lasts = s
		n = n.Add(na)
		na = na.Add(ssdecEight)
		d = d.Add(da)
		da = da.Add(ssdecThirtyTwo)
		t = t.Mul(n)
		t = t.DivRound(d, work)
		s = s.Add(t)
	}
	// -1 because shopSpring's prec == digits after radix
	return s.Round(prec - 1)
}

var (
	dnumEight     = dnum.New(+1, 8, 0)
	dnumThirtyTwo = dnum.New(+1, 32, 0)
)

// PiDnum calculates π to its maximum precision of 16 digits using
// github.com/apmckinlay/gsuneido/util/dnum.
func PiDnum() dnum.Dnum {
	var (
		lasts = dnum.New(+1, 0, 0)
		t     = dnum.New(+1, 3, 0)
		s     = dnum.New(+1, 3, 0)
		n     = dnum.New(+1, 1, 0)
		na    = dnum.New(+1, 0, 0)
		d     = dnum.New(+1, 0, 0)
		da    = dnum.New(+1, 24, 0)
	)

	for dnum.Compare(s, lasts) != 0 {
		lasts = s
		n = dnum.Add(n, na)
		na = dnum.Add(na, dnumEight)
		d = dnum.Add(d, da)
		da = dnum.Add(da, dnumThirtyTwo)
		t = dnum.Mul(t, n)
		t = dnum.Div(t, d)
		s = dnum.Add(s, t)
	}
	return s
}

// PiFloat calculates π to its maximum precision of 19 digits using Go's native
// float64.
func PiFloat64() float64 {
	var (
		lasts = 0.0
		t     = 3.0
		s     = 3.0
		n     = 1.0
		na    = 0.0
		d     = 0.0
		da    = 24.0
	)

	for s != lasts {
		lasts = s
		n += na
		na += 8
		d += da
		da += 32
		t = (t * n) / d
		s = t
	}
	return s
}

var (
	eight     = decimal.New(8, 0)
	thirtyTwo = decimal.New(32, 0)
)

// PiDecimal_Go calculates π to the desired precision using
// github.com/ericlagergren/decimal with the operating mode set to Go.
func PiDecimal_Go(prec int) *decimal.Big {
	var (
		ctx = decimal.Context{
			Precision:     adjustPrecision(prec),
			OperatingMode: decimal.Go,
		}

		lasts = new(decimal.Big)
		t     = decimal.New(3, 0)
		s     = decimal.New(3, 0)
		n     = decimal.New(1, 0)
		na    = new(decimal.Big)
		d     = new(decimal.Big)
		da    = decimal.New(24, 0)
		eps   = decimal.New(1, prec)
	)

	for {
		lasts.Copy(s)
		ctx.Add(n, n, na)
		ctx.Add(na, na, eight)
		ctx.Add(d, d, da)
		ctx.Add(da, da, thirtyTwo)
		ctx.Mul(t, t, n)
		ctx.Quo(t, t, d)
		ctx.Add(s, s, t)
		if ctx.Sub(lasts, s, lasts).CmpAbs(eps) < 0 {
			return s.Round(prec)
		}
	}
}

// PiDecimal_GDA calculates π to the desired precision using
// github.com/ericlagergren/decimal with the operating mode set to GDA.
func PiDecimal_GDA(prec int) *decimal.Big {
	var (
		ctx = decimal.Context{
			Precision:     adjustPrecision(prec),
			OperatingMode: decimal.GDA,
		}

		lasts = new(decimal.Big)
		t     = decimal.New(3, 0)
		s     = decimal.New(3, 0)
		n     = decimal.New(1, 0)
		na    = new(decimal.Big)
		d     = new(decimal.Big)
		da    = decimal.New(24, 0)
	)

	for s.Cmp(lasts) != 0 {
		lasts.Copy(s)
		ctx.Add(n, n, na)
		ctx.Add(na, na, eight)
		ctx.Add(d, d, da)
		ctx.Add(da, da, thirtyTwo)
		ctx.Mul(t, t, n)
		ctx.Quo(t, t, d)
		ctx.Add(s, s, t)
	}
	return s.Round(prec)
}

var (
	apdEight     = apd.New(8, 0)
	apdThirtyTwo = apd.New(32, 0)
)

// PiAPD calculates π to the desired precision using github.com/cockroachdb/apd.
func PiAPD(prec uint32) *apd.Decimal {
	var (
		ctx   = apd.BaseContext.WithPrecision(uint32(adjustPrecision(int(prec))))
		lasts = apd.New(0, 0)
		t     = apd.New(3, 0)
		s     = apd.New(3, 0)
		n     = apd.New(1, 0)
		na    = apd.New(0, 0)
		d     = apd.New(0, 0)
		da    = apd.New(24, 0)
	)

	for s.Cmp(lasts) != 0 {
		lasts.Set(s)
		ctx.Add(n, n, na)
		ctx.Add(na, na, apdEight)
		ctx.Add(d, d, da)
		ctx.Add(da, da, apdThirtyTwo)
		ctx.Mul(t, t, n)
		ctx.Quo(t, t, d)
		ctx.Add(s, s, t)
	}
	ctx.Precision = prec
	ctx.Round(s, s)
	return s
}

func StringToInf(x string) *inf.Dec {

	x_splited := strings.Split(x, ".")

	x_natural := x_splited[0]

	// a_natural_bigint, _ := new(big.Int).SetString(a_natural, 0)

	x_fractional := x_splited[1]

	scale := len(x_fractional)

	x_uncaled, success := new(big.Int).SetString(x_natural+x_fractional, 0)

	if !success {
		panic("fail to create inf.Dec from string")
	}

	return inf.NewDecBig(x_uncaled, inf.Scale(scale))
}

func MulInf(x string, y string, prec int, b *testing.B) string {

	x_dec := StringToInf(x)

	y_dec := StringToInf(y)

	fmt.Println(x_dec.Scale())
	fmt.Println(x_dec.UnscaledBig())

	z_dec := new(inf.Dec)

	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "go-inf"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// t.QuoRound(t, d, inf.Scale(work), inf.RoundHalfUp)
			z_dec.Mul(x_dec, y_dec)
		}
	})

	return z_dec.String()
}

func MulShopSpring(x string, y string, prec int, b *testing.B) string {
	dec_x, _ := ssdec.NewFromString(x)
	dec_y, _ := ssdec.NewFromString(y)

	// -1 because shopSpring's prec == digits after radix
	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "shopspring"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec_x.Mul(dec_y)
		}
	})

	return dec_x.Mul(dec_y).String()
}

func MulDecimal_Go(x string, y string, prec int, b *testing.B) string {
	var (
		ctx = decimal.Context{
			Precision:     prec,
			OperatingMode: decimal.Go,
		}
	)

	x_dec, _ := new(decimal.Big).SetString(x)

	y_dec, _ := new(decimal.Big).SetString(y)

	z_dec := new(decimal.Big)

	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "ericlagergren (Go)"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx.Mul(z_dec, x_dec, y_dec)
			ctx.Round(z_dec)
		}
	})

	return z_dec.String()
}

func MulDecimal_GDA(x string, y string, prec int, b *testing.B) string {
	var (
		ctx = decimal.Context{
			Precision:     prec,
			OperatingMode: decimal.GDA,
		}
	)

	x_dec, _ := new(decimal.Big).SetString(x)

	y_dec, _ := new(decimal.Big).SetString(y)

	z_dec := new(decimal.Big)

	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "ericlagergren (GDA)"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx.Mul(z_dec, x_dec, y_dec)
		}
	})

	return z_dec.String()
}

// PiAPD calculates π to the desired precision using github.com/cockroachdb/apd.
func MulAPD(x string, y string, prec int, b *testing.B) string {
	var (
		ctx = apd.BaseContext.WithPrecision(uint32(prec))
	)

	x_dec, _, err := ctx.NewFromString(x)
	require.NoError(b, err)

	y_dec, _, err := ctx.NewFromString(y)
	require.NoError(b, err)

	z_dec := new(apd.Decimal)

	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "cockroachdb/apd"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx.Mul(z_dec, x_dec, y_dec)
		}
	})

	return z_dec.String()
}

func QuoInf(x string, y string, prec int, b *testing.B) string {

	x_dec := StringToInf(x)

	y_dec := StringToInf(y)

	z_dec := new(inf.Dec)

	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "go-inf"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// t.QuoRound(t, d, inf.Scale(work), inf.RoundHalfUp)
			z_dec.QuoRound(x_dec, y_dec, inf.Scale(16), inf.RoundHalfUp)
		}
	})

	return z_dec.String()
}

func QuoShopSpring(x string, y string, prec int, b *testing.B) string {
	dec_x, _ := ssdec.NewFromString(x)
	dec_y, _ := ssdec.NewFromString(y)

	// -1 because shopSpring's prec == digits after radix
	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "shopspring"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec_x.Div(dec_y)
		}
	})

	return dec_x.Div(dec_y).String()
}

func QuoDecimal_Go(x string, y string, prec int, b *testing.B) string {
	var (
		ctx = decimal.Context{
			Precision:     prec,
			OperatingMode: decimal.Go,
		}
	)

	x_dec, _ := new(decimal.Big).SetString(x)

	y_dec, _ := new(decimal.Big).SetString(y)

	z_dec := new(decimal.Big)

	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "ericlagergren (Go)"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx.Quo(z_dec, x_dec, y_dec)
			// ctx.Round(z_dec)
		}
	})

	return z_dec.String()
}

func QuoDecimal_GDA(x string, y string, prec int, b *testing.B) string {
	var (
		ctx = decimal.Context{
			Precision:     prec,
			OperatingMode: decimal.GDA,
		}
	)

	x_dec, _ := new(decimal.Big).SetString(x)

	y_dec, _ := new(decimal.Big).SetString(y)

	z_dec := new(decimal.Big)

	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "ericlagergren (GDA)"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx.Quo(z_dec, x_dec, y_dec)
		}
	})

	return z_dec.String()
}

// PiAPD calculates π to the desired precision using github.com/cockroachdb/apd.
func QuoAPD(x string, y string, prec int, b *testing.B) string {
	var (
		ctx = apd.BaseContext.WithPrecision(uint32(prec))
	)

	x_dec, _, err := ctx.NewFromString(x)
	require.NoError(b, err)

	y_dec, _, err := ctx.NewFromString(y)
	require.NoError(b, err)

	z_dec := new(apd.Decimal)

	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "cockroachdb/apd"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx.Quo(z_dec, x_dec, y_dec)
		}
	})

	return z_dec.String()
}

func AddInf(x string, y string, prec int, b *testing.B) string {

	x_dec := StringToInf(x)

	y_dec := StringToInf(y)

	fmt.Println(x_dec.Scale())
	fmt.Println(x_dec.UnscaledBig())

	z_dec := new(inf.Dec)

	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "go-inf"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// t.QuoRound(t, d, inf.Scale(work), inf.RoundHalfUp)
			z_dec.Add(x_dec, y_dec)
		}
	})

	return z_dec.String()
}

func AddShopSpring(x string, y string, prec int, b *testing.B) string {
	dec_x, _ := ssdec.NewFromString(x)
	dec_y, _ := ssdec.NewFromString(y)

	// -1 because shopSpring's prec == digits after radix
	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "shopspring"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec_x.Add(dec_y)
		}
	})

	return dec_x.Add(dec_y).String()
}

func AddDecimal_Go(x string, y string, prec int, b *testing.B) string {
	var (
		ctx = decimal.Context{
			Precision:     prec,
			OperatingMode: decimal.Go,
		}
	)

	x_dec, _ := new(decimal.Big).SetString(x)

	y_dec, _ := new(decimal.Big).SetString(y)

	z_dec := new(decimal.Big)

	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "ericlagergren (Go)"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx.Add(z_dec, x_dec, y_dec)
			// ctx.Round(z_dec)
		}
	})

	return z_dec.String()
}

func AddDecimal_GDA(x string, y string, prec int, b *testing.B) string {
	var (
		ctx = decimal.Context{
			Precision:     prec,
			OperatingMode: decimal.GDA,
		}
	)

	x_dec, _ := new(decimal.Big).SetString(x)

	y_dec, _ := new(decimal.Big).SetString(y)

	z_dec := new(decimal.Big)

	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "ericlagergren (GDA)"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx.Add(z_dec, x_dec, y_dec)
		}
	})

	return z_dec.String()
}

// PiAPD calculates π to the desired precision using github.com/cockroachdb/apd.
func AddAPD(x string, y string, prec int, b *testing.B) string {
	var (
		ctx = apd.BaseContext.WithPrecision(uint32(prec))
	)

	x_dec, _, err := ctx.NewFromString(x)
	require.NoError(b, err)

	y_dec, _, err := ctx.NewFromString(y)
	require.NoError(b, err)

	z_dec := new(apd.Decimal)

	b.Run(fmt.Sprintf("prec=%d/pkg=%s", prec, "cockroachdb/apd"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx.Add(z_dec, x_dec, y_dec)
		}
	})

	return z_dec.String()
}

func RoundInf(x string, b *testing.B) string {

	x_dec := StringToInf(x)

	b.Run(fmt.Sprintf("pkg=%s", "go-inf"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x_dec.Round(x_dec, inf.Scale(18), inf.RoundHalfUp)
		}
	})

	return x_dec.String()
}

func RoundShopSpring(x string, b *testing.B) string {
	dec_x, _ := ssdec.NewFromString(x)

	// -1 because shopSpring's prec == digits after radix
	b.Run(fmt.Sprintf("pkg=%s", "shopspring"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec_x.Round(18)
		}
	})

	return dec_x.String()
}

func RoundDecimal_Go(x string, b *testing.B) string {
	var (
		ctx = decimal.Context{
			Precision:     28,
			OperatingMode: decimal.Go,
		}
	)

	x_dec, _ := new(decimal.Big).SetString(x)

	b.Run(fmt.Sprintf("pkg=%s", "ericlagergren (Go)"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx.Round(x_dec)
			// ctx.Round(z_dec)
		}
	})

	return x_dec.String()
}

func RoundDecimal_GDA(x string, b *testing.B) string {
	var (
		ctx = decimal.Context{
			Precision:     28,
			OperatingMode: decimal.GDA,
		}
	)

	x_dec, _ := new(decimal.Big).SetString(x)

	b.Run(fmt.Sprintf("pkg=%s", "ericlagergren (GDA)"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx.Round(x_dec)
		}
	})

	return x_dec.String()
}

// PiAPD calculates π to the desired precision using github.com/cockroachdb/apd.
func RoundAPD(x string, b *testing.B) string {
	var (
		ctx = apd.BaseContext.WithPrecision(uint32(28))
	)

	x_dec, _, err := ctx.NewFromString(x)
	require.NoError(b, err)

	b.Run(fmt.Sprintf("pkg=%s", "cockroachdb/apd"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx.Round(x_dec, x_dec)
		}
	})

	return x_dec.String()
}
