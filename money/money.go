package money

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// This can handle commas as well to a reasonable degree, but they get stripped out prior in SetString() anyway
var validMoney = regexp.MustCompile(`^\-?\(?\$?\s*\-?\s*\(?(((\d{1,3}((\,\d{3})*|\d*))?(\.\d{1,4})?)|((\d{1,3}((\,\d{3})*|\d*))(\.\d{0,4})?))\)?$`)

// Creates a new Money{} instance, with the amount initialised.  Avoid creating a Money instance directly, as it will not
// correctly initialise the raw values
func NewMoney(amount int64) Money {
	m := Money{}
	m.SetInt(amount)

	return m
}

type NullMoney struct {
	// NullMoney represents a Money instance that may be null. NullMoney implements the
	// sql.Scanner interface so it can be used as a scan destination, similar to
	// sql.NullString.
	Money *Money
	Valid bool
}

func (nm *NullMoney) SetMoney(m *Money) {
	nm.Money = m
	nm.Valid = true
}

// Scan implements the Scanner interface
func (nm *NullMoney) Scan(value interface{}) error {
	if value == nil {
		nm.Money, nm.Valid = nil, false
		return nil
	}

	m := NewMoney(0)

	nm.Money = &m

	err := nm.Money.Scan(value)

	if err != nil {
		nm.Money, nm.Valid = nil, false
		return err
	} else {
		nm.Valid = true
	}

	return nil
}

// Value implements the driver Valuer interface
func (nm NullMoney) Value() (driver.Value, error) {
	if !nm.Valid {
		return nil, nil
	}

	return nm.Money.Value()
}

type Money struct {
	i_val int64
	s_val string
}

// Scan implements the Scanner interface
func (m *Money) Scan(value interface{}) error {
	if value == nil {
		return errors.New("Value cannot be nil!  Use NullMoney instead")
	}

	raw, correct := value.([]uint8)

	if !correct {
		return errors.New(fmt.Sprintf("pq/money: Value not of type []uint8.  Type: %s, value: %+v", reflect.TypeOf(value), value))
	}

	err := m.SetString(string(raw))

	if err != nil {
		return err
	}

	return nil
}

// Value implements the driver Valuer interface
func (m Money) Value() (driver.Value, error) {
	return m.String(), nil
}

// Returns the value in a string suitable for use with postgresql, with no commas separating thousands
func (m Money) String() string {
	return m.s_val
}

// Returns an integer value for the amount, in cents.  So $1.20 = 120
func (m Money) Int64() int64 {
	return m.i_val
}

// Set the money instance to be a value as given by the string
func (m *Money) SetString(amount string) error {
	negative := strings.Contains(amount, "-")
	amount = strings.Replace(amount, ",", "", -1)
	if !validMoney.MatchString(amount) {
		return errors.New(fmt.Sprintf("String provided does not appear to be a money value.  Value: %s", amount))
	}

	str := strings.Split(amount, "$")
	f_amount, err := strconv.ParseFloat(str[len(str)-1], 64)

	if err != nil {
		return err
	}

	m.i_val = int64(f_amount * 100)

	if negative && m.i_val > 0 {
		// Sometimes ParseFloat recognises it is negative
		m.i_val = m.i_val * -1
	}
	m.s_val = m.stringFromInt(m.i_val)

	return nil
}

// Set the money instance to be an amount in cents
func (m *Money) SetInt(amount int64) {
	m.s_val = m.stringFromInt(amount)
	m.i_val = amount
}

// Add or subtract an integer amount of cents
func (m *Money) Add(amount int64) {
	m.SetInt(m.Int64() + amount)
}

func (m Money) stringFromInt(amount int64) string {
	var sign string
	if amount < 0 {
		sign = "-"
	}

	new_f := float64(amount) / 100

	if new_f < 0 {
		new_f = new_f * -1
	}
	return fmt.Sprintf("%s$%0.2f", sign, new_f)
}
