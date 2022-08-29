package benchmarks

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/dnum"
	apd "github.com/cockroachdb/apd/v3"
	"github.com/ericlagergren/decimal"
	ssdec "github.com/shopspring/decimal"
	"gopkg.in/inf.v0"
)

const pi = "3.14159265358979323846264338327950288419716939937510582097494459230" +
	"78164062862089986280348253421170679821480865132823066470938446095505822317" +
	"25359408128481117450284102701938521105559644622948954930381964428810975665" +
	"933446128475648233786783165271201909145648566923460348610454326648213394"

type testFunc func(prec int) string

// TestPiBenchmarks tests the correctness of the Pi benchmarks. It only tests
// the benchmarks that can be calculated out to a specific precision.
func TestPiBenchmarks(t *testing.T) {
	for _, tc := range [...]struct {
		name string
		fn   testFunc
	}{
		{"ericlagergren/decimal (Go)", func(prec int) string {
			return PiDecimal_Go(prec).String()
		}},
		{"ericlagergren/decimal (GDA)", func(prec int) string {
			return PiDecimal_GDA(prec).String()
		}},
		{"cockroachdb/apd", func(prec int) string {
			return PiAPD(uint32(prec)).String()
		}},
		{"shopspring/decimal", func(prec int) string {
			return PiShopSpring(int32(prec)).String()
		}},
		{"go-inf/inf", func(prec int) string {
			return PiInf(prec).String()
		}},
	} {
		var ctx decimal.Context
		for _, prec := range [...]int{9, 19, 38, 100} {
			ctx := ctx
			t.Run(fmt.Sprintf("%s/%d", tc.name, prec), func(t *testing.T) {
				ctx.Precision = prec

				str := tc.fn(prec)
				name := tc.name

				var x decimal.Big
				if _, ok := ctx.SetString(&x, str); !ok {
					t.Fatalf("%s (%d): bad input: %q", name, prec, str)
				}

				var act decimal.Big
				ctx.SetString(&act, pi)
				if act.Cmp(&x) != 0 {
					t.Fatalf(`%s (%d): bad output:
want: %q
got : %q
`, name, prec, &act, &x)
				}
			})
		}
	}
}

func BenchmarkPi(b *testing.B) {
	for _, p := range [...]int{9, 19, 38, 100} {
		for _, pkg := range [...]struct {
			pkg string
			fn  func(prec int)
		}{
			{"ericlagergren (Go)", benchmarkPi_decimal_Go},
			{"ericlagergren (GDA)", benchmarkPi_decimal_GDA},
			{"cockroachdb/apd", benchmarkPi_apd},
			{"shopspring", benchmarkPi_shopspring},
			{"apmckinlay", benchmarkPi_dnum},
			{"go-inf", benchmarkPi_inf},
			{"float64", benchmarkPi_float64},
		} {
			pkg := pkg
			b.Run(fmt.Sprintf("foo=%s", pkg.pkg), func(b *testing.B) {
				p := p
				b.Run(fmt.Sprintf("prec=%d", p), func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						pkg.fn(p)
					}
				})

			})
		}
	}
}

var gdec *decimal.Big

func benchmarkPi_decimal_Go(prec int)  { gdec = PiDecimal_Go(int(prec)) }
func benchmarkPi_decimal_GDA(prec int) { gdec = PiDecimal_GDA(int(prec)) }

var gapd *apd.Decimal

func benchmarkPi_apd(prec int) { gapd = PiAPD(uint32(prec)) }

var gssdec ssdec.Decimal

func benchmarkPi_shopspring(prec int) { gssdec = PiShopSpring(int32(prec)) }

var ginf *inf.Dec

func benchmarkPi_inf(prec int) { ginf = PiInf(int(prec)) }

var gdnum dnum.Dnum

func benchmarkPi_dnum(_ int) { gdnum = PiDnum() }

var gf float64

func benchmarkPi_float64(_ int) { gf = PiFloat64() }

func RandNumber(prec int) string {
	var num strings.Builder

	radix := rand.Intn(prec-1) + 1

	for i := 0; i < prec+1; i++ {
		if i == radix {
			num.WriteString(".")
		} else {
			num.WriteString(strconv.Itoa(rand.Intn(9) + 1))
		}
	}

	return num.String()

	// return "2.2"

}

func BenchmarkMultiplication(b *testing.B) {
	for _, p := range [...]int{9, 19, 38, 100} {
		x := RandNumber(p)
		y := RandNumber(p)
		fmt.Println(x)
		fmt.Println(y)
		for _, pkg := range [...]struct {
			pkg string
			fn  func(x string, y string, prec int, b *testing.B) string
		}{
			{"ericlagergren (Go)", MulDecimal_Go},
			{"ericlagergren (GDA)", MulDecimal_GDA},
			{"cockroachdb/apd", MulAPD},
			{"shopspring", MulShopSpring},
			{"go-inf", MulInf},
		} {
			fmt.Println(pkg.fn(x, y, p, b))
		}
	}
}

// func benchmarkMul_decimal_Go(x string, y string, prec int, b *testing.B) {
// 	gdec = MulDecimal_Go(x, y, prec, b)
// }
// func benchmarkMul_decimal_GDA(x string, y string, prec int, b *testing.B) {
// 	gdec = MulDecimal_GDA(x, y, prec, b)
// }

// func benchmarkMul_apd(x string, y string, prec int, b *testing.B) { gapd = MulAPD(x, y, prec, b) }

// func benchmarkMul_shopspring(x string, y string, prec int, b *testing.B) {
// 	gssdec = MulShopSpring(x, y, prec, b)
// }

// func benchmarkMul_inf(x string, y string, prec int, b *testing.B) { ginf = MulInf(x, y, int(prec), b) }

// func benchmarkMul_dnum(_ int) { gdnum = PiDnum() }

// func benchmarkMul_float64(_ int) { gf = PiFloat64() }
