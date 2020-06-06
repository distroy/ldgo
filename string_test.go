/*
 * Copyright (C) distroy
 */

package ldgo

import (
	"strconv"
	"testing"
)

func TestStrMapReplace(t *testing.T) {
	s0 := "a}{{a}{abc}{b}{}.{{a}}"
	m := map[string]string{
		"a":   "A",
		"b":   "B",
		"env": "test",
	}
	s1 := StrMapReplace(s0, m, "{", "}")
	t.Logf("str: %s", s0)
	t.Logf("res: %s", s1)
	t.Logf("path: %s", StrMapReplace("./config_{env}.ini", m, "{", "}"))
}

func testConvToFloat(t testing.TB, s string) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		t.Errorf("str to float fail. str:%s, err:%s", s, err.Error())
		return
	}
	t.Logf("str to float succ. str:%s, float:%f", s, f)
}

func testConvToInt(t testing.TB, s string, base int) {
	i, err := strconv.ParseInt(s, base, 64)
	if err != nil {
		t.Errorf("str to float fail. str:%s, base:%d, err:%s", s, base, err.Error())
		return
	}
	t.Logf("str to float succ. str:%s, base:%d, float:%d", s, base, i)
}

func TestConvNumber(t *testing.T) {
	testConvToFloat(t, ".1E3")
	testConvToFloat(t, ".1e-3")
	// testConvToFloat(t, ".1e")
	testConvToFloat(t, ".1e+1")
	// testConvToFloat(t, ".1e++1")
	testConvToInt(t, "0123", 10)
	testConvToInt(t, "0123", 8)
	// testConvToInt(t, "0x123", 16)
	testConvToInt(t, "0123", 10)
}

func TestStrIMapReplace(t *testing.T) {
	s0 := "a}{{a}{abc}{b}{}.{{a}}"
	m0 := map[string]string{
		"a": "A",
		"b": "B",
	}
	s1 := StrIMapReplace(s0, m0, "{", "}")
	t.Logf("str: %s", s0)
	t.Logf("res: %s", s1)

	m1 := map[string]int{
		"a": 1,
		"b": 2,
	}
	s2 := StrIMapReplace(s0, m1, "{", "}")
	t.Logf("str: %s", s0)
	t.Logf("res: %s", s2)
}
