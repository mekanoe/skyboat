package util

import "testing"
import "os"

func TestStringEnv(t *testing.T) {
	os.Setenv("TEST::stringenv", "abc,def")
	val := Getenvdef("TEST::stringenv", "def,ghi")

	if val.Empty {
		t.Error("value is empty on a key that's definitely set.")
		t.Fail()
		return
	}

	str, err := val.String()
	if err != nil {
		t.Fatal(err)
		return
	}

	if str != "abc,def" {
		t.Errorf("expected `abc,def` got `%s`", str)
	}
}

func TestBoolEnv(t *testing.T) {
	os.Setenv("TEST::boolenv", "1")
	val := Getenvdef("TEST::boolenv", "0")

	if val.Empty {
		t.Error("value is empty on a key that's definitely set.")
		t.Fail()
		return
	}

	str, err := val.Bool()
	if err != nil {
		t.Fatal(err)
		return
	}

	if !str {
		t.Errorf("expected `abc,def` got `%v`", str)
	}

	val2, err := Getenvdef("TEST::bool2env", "0").Bool()
	if err != nil {
		t.Error(err)
	}
	if val2 {
		t.Error("should be false")
	}
}

func TestBytesEnv(t *testing.T) {
	os.Setenv("TEST::bytesenv", "abc,def")
	val := Getenvdef("TEST::bytesenv", "def,ghi")

	if val.Empty {
		t.Error("value is empty on a key that's definitely set.")
		t.Fail()
		return
	}

	str, err := val.Bytes()
	if err != nil {
		t.Fatal(err)
		return
	}

	if string(str) != "abc,def" {
		t.Errorf("expected `abc,def` got `%s`", str)
	}

	val2, _ := Getenvdef("TEST::bytes2env", []byte("def,ghi")).Bytes()
	if string(val2) != "def,ghi" {
		t.Error("bytes as bytes default problem")
	}

	val3, _ := Getenvdef("TEST::bytes3env", "def,ghi").Bytes()
	if string(val3) != "def,ghi" {
		t.Error("bytes as string default problem")
	}

}

func TestStringSliceEnv(t *testing.T) {
	os.Setenv("TEST::stringsliceenv", "abc,def")
	val := Getenvdef("TEST::stringsliceenv", "def,ghi")

	if val.Empty {
		t.Error("value is empty on a key that's definitely set.")
		t.Fail()
		return
	}

	sli, err := val.StringSlice()
	if err != nil {
		t.Fatal(err)
		return
	}

	if sli[0] != "abc" || sli[1] != "def" {
		str, _ := val.String()
		t.Errorf("expected `abc,def` got `%s`", str)
	}
}

func TestEmpty(t *testing.T) {
	val2 := Getenvdef("TEST::emptyenv", "def,ghi")

	if !val2.Empty {
		str, _ := val2.String()
		t.Errorf("expected empty, got `%s`", str)
	}
}

func TestIntEnv(t *testing.T) {
	os.Setenv("TEST::intenv", "2013")
	val := Getenvdef("TEST::intenv", 2013)

	if val.Empty {
		t.Error("value is empty on a key that's definitely set.")
		t.Fail()
		return
	}

	num, err := val.Int()
	if err != nil {
		t.Fatal(err)
		return
	}

	if num != 2013 {
		t.Errorf("expected `2013` got `%d`", num)
		t.Fail()
	}

	val2 := Getenvdef("TEST::int2env", 2014)
	num2, _ := val2.Int()
	if num2 != 2014 {
		t.Errorf("expected `2014` got `%d`", num)
	}
}
