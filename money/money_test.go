package money

import (
	// "github.com/lib/pq"
	"testing"
)

var valueTests = []struct {
	str          string
	expected     int64
	str_expected string
}{
	{"$10.05", 1005, "$10.05"},
	{"$10", 1000, "$10.00"},
	{"10.05", 1005, "$10.05"},
	{"10", 1000, "$10.00"},
	{"123423423", 12342342300, "$123423423.00"},
	{"$123423423", 12342342300, "$123423423.00"},
	{"$123423423.00", 12342342300, "$123423423.00"},
	{"$123423423.05", 12342342305, "$123423423.05"},
	{"123,423,423", 12342342300, "$123423423.00"},
	{"1234234.23", 123423423, "$1234234.23"},
	{"$1234234.23", 123423423, "$1234234.23"},
	{"1,234,234.23", 123423423, "$1234234.23"},
	{"-$10.05", -1005, "-$10.05"},
	{"-$10", -1000, "-$10.00"},
	{"-10.05", -1005, "-$10.05"},
	{"-10", -1000, "-$10.00"},
	{"-123423423", -12342342300, "-$123423423.00"},
	{"-$123423423", -12342342300, "-$123423423.00"},
	{"-$123423423.00", -12342342300, "-$123423423.00"},
	{"-$123423423.05", -12342342305, "-$123423423.05"},
	{"-123,423,423", -12342342300, "-$123423423.00"},
	{"-1234234.23", -123423423, "-$1234234.23"},
	{"-$1234234.23", -123423423, "-$1234234.23"},
	{"-1,234,234.23", -123423423, "-$1234234.23"},
}

var sumTests = []struct {
	str          string
	num          int64
	add          int64
	expected     int64
	str_expected string
}{
	{"$0.00", 0, 5, 5, "$0.05"},
	{"$0.00", 0, -5, -5, "-$0.05"},
	{"$0.00", 0, 0, 0, "$0.00"},
	{"$5.34", 534, 100, 634, "$6.34"},
	{"$5.34", 534, -100, 434, "$4.34"},
	{"$5.34", 534, -600, -66, "-$0.66"},
	{"-$5.34", -534, 100, -434, "-$4.34"},
	{"-$5.34", -534, -100, -634, "-$6.34"},
	{"-$5.34", -534, -600, -1134, "-$11.34"},
}

func TestDefault(t *testing.T) {
	for _, v := range valueTests {
		m := NewMoney(v.expected)

		if m.Int64() != v.expected {
			t.Errorf("Used NewMoney with %d, but int value reports as %d\n", v.expected, m.Int64())
		}

		if m.String() != v.str_expected {
			t.Errorf("Used NewMoney with %d, expected %s but string is %s\n", v.expected, v.str_expected, m.String())
		}
	}
}

func TestAssign(t *testing.T) {
	m := NewMoney(0)

	for _, v := range valueTests {
		m.SetString(v.str)
		if m.Int64() != v.expected {
			t.Errorf("Gave string %s and expected int %d, instead got int %d\n", v.str, v.expected, m.Int64())
		}

		if m.String() != v.str_expected {
			t.Errorf("Returned string should have been %s, but was %s\n", v.str_expected, m.String())
		}

		// Now try assigning the int's directly:
		m.SetInt(v.expected)
		if m.Int64() != v.expected {
			t.Errorf("Gave int %d, but got int %d\n", v.expected, m.Int64())
		}

		if m.String() != v.str_expected {
			t.Errorf("Returned string should have been %s, but was %s\n", v.str_expected, m.String())
		}
	}
}

func TestScanMoney(t *testing.T) {
	var m Money
	for _, v := range valueTests {
		err := m.Scan([]uint8(v.str_expected))

		if err != nil {
			t.Errorf("Error scanning: %s\n", err)
		}

		if m.Int64() != v.expected {
			t.Errorf("Scanned string %s and expected int %d, instead got int %d\n", v.str, v.expected, m.Int64())
		}

		if m.String() != v.str_expected {
			t.Errorf("Returned string should have been %s, but was %s\n", v.str_expected, m.String())
		}

		val, err := m.Value()

		if err != nil {
			t.Errorf("Money value error: %s\n", err)
		}

		if val != v.str_expected {
			t.Errorf("Money Value string should have been %s, but was %s\n", v.str_expected, val)
		}
	}

	for _, v := range valueTests {
		var nm NullMoney
		err := nm.Scan([]uint8(v.str_expected))

		if err != nil {
			t.Errorf("Error scanning: %s\n", err)
		}

		if !nm.Valid {
			t.Errorf("NullMoney should be valid but returned invalid\n")
		}

		if nm.Money.Int64() != v.expected {
			t.Errorf("Scanned string %s and expected int %d, instead got int %d\n", v.str, v.expected, nm.Money.Int64())
		}

		if nm.Money.String() != v.str_expected {
			t.Errorf("Returned string should have been %s, but was %s\n", v.str_expected, nm.Money.String())
		}

		val, err := nm.Value()

		if err != nil {
			t.Errorf("NullMoney value error: %s\n", err)
		}

		if val != v.str_expected {
			t.Errorf("NullMoney Value string should have been %s, but was %s\n", v.str_expected, val)
		}
	}
}

func TestScanNil(t *testing.T) {
	var nm NullMoney
	nm.Scan(nil)
	if nm.Valid {
		t.Errorf("NullMoney valid when should be returning invalid\n")
	}
}

func TestSum(t *testing.T) {
	for _, v := range sumTests {
		m := NewMoney(v.num)

		if m.String() != v.str {
			t.Errorf("m.String(): %s, v.str: %s, initialised with v.num: %d\n", m.String(), v.str, v.num)
		}

		if m.Int64() != v.num {
			t.Errorf("m.Int64: %d, v.num: %d\n", m.Int64(), v.num)
		}

		m.Add(v.add)
		if m.String() != v.str_expected {
			t.Errorf("m.String(): %s, v.str_expected: %s, added v.add: %d\n", m.String(), v.str_expected, v.add)
		}

		if m.Int64() != v.expected {
			t.Errorf("m.Int64: %d, v.expected: %d\n", m.Int64(), v.expected)
		}

	}
}
