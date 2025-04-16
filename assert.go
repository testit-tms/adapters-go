package tms

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func True(t *testing.T, value bool, msgAndArgs ...interface{}) {
	success := assert.True(t, value, msgAndArgs...)
	if !success {
		setTestMessage("Should be true")
		t.Fail()
	}
}

func Exactly(t *testing.T, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	success := assert.Exactly(t, expected, actual, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Types expected to match exactly\n\t%s != %s", expected, actual))
		t.Fail()
	}
}

func Same(t *testing.T, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	success := assert.Same(t, expected, actual, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Not same: \n"+
			"expected: %p %#v\n"+
			"actual  : %p %#v", expected, expected, actual, actual))
		t.Fail()
	}
}

func NotSame(t *testing.T, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	success := assert.NotSame(t, expected, actual, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Expected and actual point to the same object: %p %#v",
			expected, expected))
		t.Fail()
	}
}

func Equal(t *testing.T, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	success := assert.Equal(t, expected, actual, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Not equal: \n"+
			"expected: %s\n"+
			"actual  : %s", expected, actual))
		t.Fail()
	}
}

func NotEqual(t *testing.T, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	success := assert.NotEqual(t, expected, actual, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Should not be: %#v\n", actual))
		t.Fail()
	}
}

func EqualValues(t *testing.T, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	success := assert.EqualValues(t, expected, actual, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Not equal: \n"+
			"expected: %s\n"+
			"actual  : %s", expected, actual))
		t.Fail()
	}
}

func NotEqualValues(t *testing.T, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	success := assert.NotEqualValues(t, expected, actual, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Should not be: %#v\n", actual))
		t.Fail()
	}
}

func Error(t *testing.T, err error, msgAndArgs ...interface{}) {
	success := assert.Error(t, err, msgAndArgs...)
	if !success {
		setTestMessage("An error is expected but got nil.")
		t.Fail()
	}
}

func NoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	success := assert.NoError(t, err, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Received unexpected error:\n%+v", err))
		t.Fail()
	}
}

func EqualError(t *testing.T, theError error, errString string, msgAndArgs ...interface{}) {
	var actualString string

	if theError != nil {
		actualString = theError.Error()
	}
	success := assert.EqualError(t, theError, errString, msgAndArgs...)

	if !success {
		setTestMessage(fmt.Sprintf("Error message not equal:\n"+
			"expected: %q\n"+
			"actual  : %q", errString, actualString))
		t.Fail()
	}
}

func ErrorIs(t *testing.T, err error, target error, msgAndArgs ...interface{}) {
	var (
		actualString string
		targetString string
	)

	if target != nil {
		targetString = target.Error()
	}

	if err != nil {
		actualString = err.Error()
	}
	success := assert.ErrorIs(t, err, target, msgAndArgs...)

	if !success {
		setTestMessage(fmt.Sprintf("Target error should be in err chain:\n"+
			"expected: %s\n"+
			"in chain: %s", targetString, actualString))
		t.Fail()
	}
}

func ErrorAs(t *testing.T, err error, target interface{}, msgAndArgs ...interface{}) {
	var (
		errorString string
	)

	if err != nil {
		errorString = err.Error()
	}
	success := assert.ErrorAs(t, err, target, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Should be in error chain:\n"+
			"expected: %s\n"+
			"in chain: %s", target, errorString))
		t.Fail()
	}
}

func Nil(t *testing.T, object interface{}, msgAndArgs ...interface{}) {
	success := assert.Nil(t, object, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Expected nil, but got: %#v", object))
		t.Fail()
	}
}

func NotNil(t *testing.T, object interface{}, msgAndArgs ...interface{}) {
	success := assert.NotNil(t, object, msgAndArgs...)
	if !success {
		setTestMessage("Expected value not to be nil.")
		t.Fail()
	}
}

func Len(t *testing.T, object interface{}, length int, msgAndArgs ...interface{}) {
	success := assert.Len(t, object, length, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("\"%s\" should have %d item(s)", object, length))
		t.Fail()
	}
}

func Contains(t *testing.T, s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	success := assert.Contains(t, s, contains, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("%#v does not contain %#v", s, contains))
		t.Fail()
	}
}

func NotContains(t *testing.T, s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	success := assert.NotContains(t, s, contains, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("%#v should not contain %#v", s, contains))
		t.Fail()
	}
}

func Greater(t *testing.T, e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	success := assert.Greater(t, e1, e2, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("\"%v\" is not greater than \"%v\"", e1, e2))
		t.Fail()
	}
}

func GreaterOrEqual(t *testing.T, e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	success := assert.GreaterOrEqual(t, e1, e2, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("\"%v\" is not greater than or equal to \"%v\"", e1, e2))
		t.Fail()
	}
}

func Less(t *testing.T, e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	success := assert.Less(t, e1, e2, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("\"%v\" is not greater than \"%v\"", e1, e2))
		t.Fail()
	}
}

func LessOrEqual(t *testing.T, e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	success := assert.LessOrEqual(t, e1, e2, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("\"%v\" is not less than or equal to \"%v\"", e1, e2))
		t.Fail()
	}
}

func Implements(t *testing.T, interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) {
	success := assert.Implements(t, interfaceObject, object, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("%T must implement %v", object, interfaceObject))
		t.Fail()
	}
}

func Empty(t *testing.T, object interface{}, msgAndArgs ...interface{}) {
	success := assert.Empty(t, object, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Should be empty, but was %v", object))
		t.Fail()
	}
}

func NotEmpty(t *testing.T, object interface{}, msgAndArgs ...interface{}) {
	success := assert.NotEmpty(t, object, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Should NOT be empty, but was %v", object))
		t.Fail()
	}
}

func WithinDuration(t *testing.T, expected, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) {
	success := assert.WithinDuration(t, expected, actual, delta, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Max difference between %v and %v allowed is %v", expected, actual, delta))
		t.Fail()
	}
}

func JSONEq(t *testing.T, expected, actual string, msgAndArgs ...interface{}) {
	success := assert.JSONEq(t, expected, actual, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("JSONs are not equal:\n"+
			"expected: %s\n"+
			"actual  : %s", expected, actual))
		t.Fail()
	}
}

func Subset(t *testing.T, list, subset interface{}, msgAndArgs ...interface{}) {
	success := assert.Subset(t, list, subset, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("%v is not a subset of %v", subset, list))
		t.Fail()
	}
}

func NotSubset(t *testing.T, list, subset interface{}, msgAndArgs ...interface{}) {
	success := assert.NotSubset(t, list, subset, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("%v should not be a subset of %v", subset, list))
		t.Fail()
	}
}

func IsType(t *testing.T, expectedType interface{}, object interface{}, msgAndArgs ...interface{}) {
	success := assert.IsType(t, expectedType, object, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Object expected to be of type %v, but was %v", reflect.TypeOf(expectedType), reflect.TypeOf(object)))
		t.Fail()
	}
}

func False(t *testing.T, value bool, msgAndArgs ...interface{}) {
	success := assert.False(t, value, msgAndArgs...)
	if !success {
		setTestMessage("Should be false")
		t.Fail()
	}
}

func Regexp(t *testing.T, rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	success := assert.Regexp(t, rx, str, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Expect \"%v\" to match \"%v\"", str, rx))
		t.Fail()
	}
}

func ElementsMatch(t *testing.T, listA interface{}, listB interface{}, msgAndArgs ...interface{}) {
	success := assert.ElementsMatch(t, listA, listB, msgAndArgs...)
	if !success {
		setTestMessage("Elements do not match")
		t.Fail()
	}
}

func DirExists(t *testing.T, path string, msgAndArgs ...interface{}) {
	success := assert.DirExists(t, path, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Directory \"%v\" does not exist", path))
		t.Fail()
	}
}

func Condition(t *testing.T, condition assert.Comparison, msgAndArgs ...interface{}) {
	success := assert.Condition(t, condition, msgAndArgs...)
	if !success {
		setTestMessage("Condition failed!")
		t.Fail()
	}
}

func Zero(t *testing.T, i interface{}, msgAndArgs ...interface{}) {
	success := assert.Zero(t, i, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Should be zero, but was %v", i))
		t.Fail()
	}
}

func NotZero(t *testing.T, i interface{}, msgAndArgs ...interface{}) {
	success := assert.NotZero(t, i, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Should not be zero, but was %v", i))
		t.Fail()
	}
}

func InDelta(t *testing.T, expected, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	success := assert.InDelta(t, expected, actual, delta, msgAndArgs...)
	if !success {
		setTestMessage(fmt.Sprintf("Difference between %v and %v should be less than %v", expected, actual, delta))
		t.Fail()
	}
}

func setTestMessage(msg string) {
	manipulateOnObjectFromCtx(
		testResultKey,
		func(test interface{}) {
			testStatusDetails := test.(hasErrorFields)
			testStatusDetails.addTrace(getTrace())
			testStatusDetails.addMessage(msg)
			testStatusDetails.addStatus(failed)
		})
	manipulateOnObjectFromCtx(
		nodeKey,
		func(node interface{}) {
			n := node.(hasStatus)
			n.addStatus(failed)
		})
}

func getTrace() string {
	traces := make([]string, 0)
	for _, v := range assert.CallerInfo() {
		if strings.Contains(v, "pkg/tms") || strings.Contains(v, "jtolds/gls") {
			continue
		}
		traces = append(traces, v)
	}

	return strings.Join(traces, "\n")
}
