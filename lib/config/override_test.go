package config

import (
	"fmt"
	"os"
	"testing"
)

type testConfig struct {
	TestBool    bool
	TestFloat32 float32
	TestFloat64 float64
	TestInt     int
	TestInt8    int8
	TestInt16   int16
	TestInt32   int32
	TestInt64   int64
	TestUint    uint
	TestUint8   uint8
	TestUint16  uint16
	TestUint32  uint32
	TestUint64  uint64
	TestString  string
}

func TestOverride(t *testing.T) {
	testConf := testConfig{}

	var envTests = []struct {
		key   string
		value string
	}{
		{"TESTBOOL", "true"},
		{"TESTFLOAT32", "32.32"},
		{"TESTFLOAT64", "64.64"},
		{"TESTINT", "4"},
		{"TESTINT8", "8"},
		{"TESTINT16", "16"},
		{"TESTINT32", "32"},
		{"TESTINT64", "64"},
		{"TESTUINT", "4"},
		{"TESTUINT8", "8"},
		{"TESTUINT16", "16"},
		{"TESTUINT32", "32"},
		{"TESTUINT64", "64"},
	}

	for pos := range envTests {
		if err := os.Setenv(fmt.Sprintf("TESTCONFIG_%s", envTests[pos].key), envTests[pos].value); err != nil {
			t.Fatalf("Setting TESTCONFIG_ failed: %s", err.Error())
		}
		defer os.Unsetenv(fmt.Sprintf("TESTCONFIG_%s", envTests[pos].key))
	}

	Override(&testConf, "", "_", &KVenv{})

	if testConf.TestBool != true {
		t.Fatalf("TestBool is not True")
	}

	if testConf.TestFloat32 != 32.32 {
		t.Fatalf("TestFloat32 does not equal 32.32")
	}

	if testConf.TestFloat64 != 64.64 {
		t.Fatalf("TestFloat64 does not equal 64.64")
	}

	if testConf.TestInt != 4 {
		t.Fatalf("TestInt does not equal 4")
	}

	if testConf.TestInt8 != 8 {
		t.Fatalf("TestInt8 does not equal 8")
	}

	if testConf.TestInt16 != 16 {
		t.Fatalf("TestInt16 does not equal 16")
	}

	if testConf.TestInt32 != 32 {
		t.Fatalf("TestInt32 does not equal 32")
	}

	if testConf.TestInt64 != 64 {
		t.Fatalf("TestInt64 does not equal 64")
	}

	if testConf.TestUint != 4 {
		t.Fatalf("TestUint does not equal 4")
	}

	if testConf.TestUint8 != 8 {
		t.Fatalf("TestUint8 does not equal 8")
	}

	if testConf.TestUint16 != 16 {
		t.Fatalf("TestUint16 does not equal 16")
	}

	if testConf.TestInt32 != 32 {
		t.Fatalf("TestInt32 does not equal 32")
	}

	if testConf.TestUint64 != 64 {
		t.Fatalf("TestUint64 does not equal 64")
	}
}

// EOF
