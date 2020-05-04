package enums

import (
	"testing"

	"github.com/marianolg/enums"
)

type testFunc func(*testing.T)

func testItPanics(wrapped func(), notPanicMsg string) testFunc {
	return func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf(notPanicMsg)
			}
		}()
		wrapped()
	}
}

func TestCreateEmptyEnumPanics(t *testing.T) {
	testItPanics(func() { enums.New() }, "Empty enum creation should panic")(t)
}

func TestCreateEnumOfNonHashableValuesPanics(t *testing.T) {
	testItPanics(func() { enums.New([]int{1}) },
		"Enum of non-hashable values should panic")(t)
}

func TestCreateEnumWithDifferentTypesPanics(t *testing.T) {
	testItPanics(func() { enums.New(1, "2") },
		"Enum with diffrent types of values should panic")(t)
}

func TestCreateEnumWithDupedValuePanics(t *testing.T) {
	testItPanics(func() { enums.New(1, 1) }, "Enum with duped values should panic")(t)
}

func TestCreateEnum(t *testing.T) {
	enums.New(1, 2, 3)
}

func TestCheckingValidityOfUnsupportedTypePanics(t *testing.T) {
	enum := enums.New(1, 2, 3)
	testItPanics(func() { enum.IsValid("A") },
		"Checking on a string on a int enum should panic")(t)
	testItPanics(func() { enum.IsAnyValid("A", 1) },
		"Checking on a string on a int enum should panic")(t)
	testItPanics(func() { enum.AreAllValid(1, "A") },
		"Checking on a string on a int enum should panic")(t)
}

func TestIsValid(t *testing.T) {
	e := enums.New(1, 2, 3)

	if !e.IsValid(1) {
		t.Errorf("1 should be considered valid")
	}

	if e.IsValid(4) {
		t.Errorf("4 should NOT be considered valid")
	}
}

func TestIsAnyValid(t *testing.T) {
	e := enums.New(1, 2, 3)

	if e.IsAnyValid() {
		t.Errorf("IsAnyValid of no values should return false")
	}

	if !e.IsAnyValid(1, 2) {
		t.Errorf("IsAnyValid should be true for [1 2] (both values should be valid)")
	}

	if !e.IsAnyValid(1, 4) {
		t.Errorf("IsAnyValid should be true for [1 4] (1 should be valid)")
	}

	if e.IsAnyValid(4, 5) {
		t.Errorf("IsAnyValid should be false for [4 5] (neither value should be valid)")
	}
}

func TestAreAllValid(t *testing.T) {
	e := enums.New(1, 2, 3)

	if !e.AreAllValid() {
		t.Errorf("AreAllValid of no values should return true")
	}

	if !e.AreAllValid(1, 2) {
		t.Errorf("AreAllValid should be true for [1 2] (both values should be valid)")
	}

	if e.AreAllValid(1, 4) {
		t.Errorf("AreAllValid should be false for [1 4] (4 should be invalid)")
	}

	if e.AreAllValid(4, 5) {
		t.Errorf("AreAllValid should be false for [4 5] (both values should be invalid)")
	}
}

func TestStringEnum(t *testing.T) {
	enum := enums.New("A", "B", "C")

	if !enum.IsValid("A") {
		t.Errorf("\"A\" should be considered valid")
	}
	if enum.IsValid("D") {
		t.Errorf("\"D\" should NOT be considered valid")
	}
}

func TestConvertEnum(t *testing.T) {
	enumNoConvert := enums.New(1, 2, 3)
	enumConvert := enums.NewConvert(1, 2, 3)

	testItPanics(func() { enumNoConvert.IsValid(1.0) },
		"Checking for validity of a float in an no-convert int Enum should panic")
	if !enumConvert.IsValid(1.0) {
		t.Errorf("1.0 should be considered valid")
	}
}

func TestSetTypedMethodsWithWrongFuncTypesPanics(t *testing.T) {
	enum := enums.New(1, 2, 3)
	enumConvert := enums.NewConvert(1, 2, 3)

	testItPanics(func() {
		var f func(int) bool
		enum.SetTypedIsValid(f)
	}, "Should panic - not a pointer")

	testItPanics(func() {
		var f int
		enum.SetTypedIsValid(&f)
	}, "Should panic - pointer to non-func type")

	testItPanics(func() {
		var f func() bool
		enum.SetTypedIsValid(&f)
	}, "Should panic - wrong in args count")

	testItPanics(func() {
		var f func(int, int) bool
		enum.SetTypedIsValid(&f)
	}, "Should panic - wrong in args count")

	testItPanics(func() {
		var f func(string) bool
		enum.SetTypedIsValid(&f)
	}, "Should panic - wrong in arg type")

	testItPanics(func() {
		var f func(float32) bool
		enum.SetTypedIsValid(&f)
	}, "Should panic - wrong in arg type")

	testItPanics(func() {
		var f func(string) bool
		enumConvert.SetTypedIsValid(&f)
	}, "Should panic - wrong in arg type")

	// Should NOT panic because float32 is convertible to int
	var f func(float64) bool
	enumConvert.SetTypedIsValid(&f)

	testItPanics(func() {
		var f func(...int) bool
		enum.SetTypedIsValid(&f)
	}, "Should panic - variadic in arg type of IsValid")

	testItPanics(func() {
		var f func(int) bool
		enum.SetTypedIsAnyValid(&f)
	}, "Should panic - non-variadic in arg type of IsAnyValid")

	testItPanics(func() {
		var f func(...int) bool
		enum.SetTypedIsAnyValid(&f)
	}, "Should panic - non-variadic in arg type of AreAllValid")

	testItPanics(func() {
		var f func(int)
		enum.SetTypedIsValid(&f)
	}, "Should panic - wrong out args count")

	testItPanics(func() {
		var f func(int) (bool, bool)
		enum.SetTypedIsValid(&f)
	}, "Should panic - wrong out args count")

	testItPanics(func() {
		var f func(int) int
		enum.SetTypedIsValid(&f)
	}, "Should panic - wrong out arg type")
}

func TestTypedIsValid(t *testing.T) {
	enum := enums.New(1, 2, 3)

	var intEnumIsValid func(int) bool
	enum.SetTypedIsValid(&intEnumIsValid)

	if !intEnumIsValid(1) {
		t.Errorf("1 should be considered valid")
	}

	if intEnumIsValid(4) {
		t.Errorf("4 should NOT be considered valid")
	}
}

func TestTypedIsAnyValid(t *testing.T) {
	enum := enums.New(1, 2, 3)

	var intEnumIsAnyValid func(...int) bool
	enum.SetTypedIsAnyValid(&intEnumIsAnyValid)

	if !intEnumIsAnyValid(1, 2) {
		t.Errorf("IsAnyValid should be true for [1 2] (both values should be valid)")
	}

	if !intEnumIsAnyValid(1, 4) {
		t.Errorf("IsAnyValid should be true for [1 4] (1 should be valid)")
	}

	if intEnumIsAnyValid(4, 5) {
		t.Errorf("IsAnyValid should be false for [4 5] (neither value should be valid)")
	}
}

func TestTypedAreAllValid(t *testing.T) {
	enum := enums.New(1, 2, 3)

	var intEnumAreAllValid func(...int) bool
	enum.SetTypedAreAllValid(&intEnumAreAllValid)

	if !intEnumAreAllValid(1, 2) {
		t.Errorf("AreAllValid should be true for [1 2] (both values should be valid)")
	}

	if intEnumAreAllValid(1, 4) {
		t.Errorf("AreAllValid should be false for [1 4] (4 should be invalid)")
	}

	if intEnumAreAllValid(4, 5) {
		t.Errorf("AreAllValid should be false for [4 5] (both values should be invalid)")
	}
}
