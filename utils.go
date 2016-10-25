package opay

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func checkTimeout(deadline time.Time) (timeout time.Duration, errTimeout error) {
	// No timeout
	if deadline.IsZero() {
		return
	}

	timeout = deadline.Sub(time.Now())

	// Timeout, cancel order.
	if timeout <= 0 {
		return timeout, ErrTimeout
	}

	// No timeout
	return
}

type Floater struct {
	numOfDecimalPlaces int
	accuracy           float64
	zeroString         string
	format             string
}

func NewFloater(numOfDecimalPlaces int) *Floater {
	if numOfDecimalPlaces < 0 || numOfDecimalPlaces > 14 {
		panic("the range of Floater.numOfDecimalPlaces must be between 0 and 14.")
	}
	if numOfDecimalPlaces == 0 {
		return &Floater{
			numOfDecimalPlaces: 0,
			accuracy:           0,
			zeroString:         "0",
			format:             "%0.0f",
		}
	}
	accuracyString := "0." + strings.Repeat("0", numOfDecimalPlaces-1) + "1"
	accuracy, _ := strconv.ParseFloat(accuracyString, 64)
	zeroString := accuracyString[:len(accuracyString)-1] + "0"
	return &Floater{
		numOfDecimalPlaces: numOfDecimalPlaces,
		accuracy:           accuracy,
		zeroString:         zeroString,
		format:             "%0." + strconv.Itoa(numOfDecimalPlaces) + "f",
	}
}

func (this *Floater) NumOfDecimalPlaces() int {
	return this.numOfDecimalPlaces
}

func (this *Floater) Accuracy() float64 {
	return this.accuracy
}

func (this *Floater) Format() string {
	return this.format
}

func (this *Floater) Ftoa(f float64) string {
	return fmt.Sprintf(this.format, f)
}

func (this *Floater) Atof(s string, bitSize int) (float64, error) {
	f, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return f, err
	}
	return strconv.ParseFloat(fmt.Sprintf(this.format, f), bitSize)
}

func (this *Floater) Ftof(f float64) float64 {
	f, _ = strconv.ParseFloat(fmt.Sprintf(this.format, f), 64)
	return f
}

func (this *Floater) Atoa(s string, bitSize int) (string, error) {
	f, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return s, err
	}
	return fmt.Sprintf(this.format, f), nil
}

func (this *Floater) Equal(a, b float64) bool {
	return this.IsZero(a - b)
}

func (this *Floater) Greater(a, b float64) bool {
	return math.Max(a, b) == a && !this.IsZero(a-b)
}

func (this *Floater) GreaterOrEqual(a, b float64) bool {
	return math.Max(a, b) == a || this.IsZero(a-b)
}

func (this *Floater) Smaller(a, b float64) bool {
	return math.Min(a, b) == a && !this.IsZero(a-b)
}

func (this *Floater) SmallerOrEqual(a, b float64) bool {
	return math.Min(a, b) == a || this.IsZero(a-b)
}

func (this *Floater) IsZero(a float64) bool {
	return this.Ftoa(a) <= this.zeroString
}
