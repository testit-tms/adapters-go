package htmlutils

import (
	"os"
	"reflect"
	"regexp"
	"strings"
)

const (
	NoEscapeHtmlEnvVar = "NO_ESCAPE_HTML"
)

var (
	// Regex pattern to detect HTML tags
	htmlTagPattern = regexp.MustCompile(`<\S.*?(?:>|\/>)`)
)

// EscapeHtmlTags escapes HTML tags to prevent XSS attacks.
// First checks if the string contains HTML tags using regex pattern.
// Only performs escaping if HTML tags are detected.
// Escapes all < as \< and > as \> avoiding double escaping.
// Can be disabled by setting NO_ESCAPE_HTML environment variable to "true"
func EscapeHtmlTags(text string) string {
	if text == "" {
		return text
	}

	// Check if escaping is disabled via environment variable
	if strings.EqualFold(os.Getenv(NoEscapeHtmlEnvVar), "true") {
		return text
	}

	// First check if the string contains HTML tags
	if !htmlTagPattern.MatchString(text) {
		return text // No HTML tags found, return original string
	}

	// Simple approach: replace all < and > with escaped versions
	// Check if already escaped to avoid double escaping
	result := text

	// Replace unescaped < with \<
	result = strings.ReplaceAll(result, "<", "&lt;")

	// Replace unescaped > with \>
	result = strings.ReplaceAll(result, ">", "&gt;")

	return result
}

// EscapeHtmlInObject escapes HTML tags in all string fields and properties of an object using reflection
// Also processes slice/array fields: if slice of objects - processes each object,
// if slice of strings - escapes each string
// Can be disabled by setting NO_ESCAPE_HTML environment variable to "true"
func EscapeHtmlInObject(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}

	// Check if escaping is disabled via environment variable
	if strings.EqualFold(os.Getenv(NoEscapeHtmlEnvVar), "true") {
		return obj
	}

	processValue(reflect.ValueOf(obj))
	return obj
}

// EscapeHtmlInObjectSlice escapes HTML tags in all string fields of objects in a slice using reflection
// Can be disabled by setting NO_ESCAPE_HTML environment variable to "true"
func EscapeHtmlInObjectSlice(slice interface{}) interface{} {
	if slice == nil {
		return nil
	}

	// Check if escaping is disabled via environment variable
	if strings.EqualFold(os.Getenv(NoEscapeHtmlEnvVar), "true") {
		return slice
	}

	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice && sliceValue.Kind() != reflect.Array {
		return slice
	}

	for i := 0; i < sliceValue.Len(); i++ {
		processValue(sliceValue.Index(i))
	}

	return slice
}

func processValue(value reflect.Value) {
	// Handle pointers
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return
		}
		processValue(value.Elem())
		return
	}

	// Handle interfaces
	if value.Kind() == reflect.Interface {
		if value.IsNil() {
			return
		}
		processValue(value.Elem())
		return
	}

	switch value.Kind() {
	case reflect.String:
		if value.CanSet() {
			escaped := EscapeHtmlTags(value.String())
			value.SetString(escaped)
		}

	case reflect.Struct:
		processStruct(value)

	case reflect.Slice, reflect.Array:
		processSlice(value)

	case reflect.Map:
		processMap(value)
	}
}

func processStruct(structValue reflect.Value) {
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		// Skip if field is not settable
		if !field.CanSet() {
			continue
		}

		processValue(field)
	}
}

func processSlice(sliceValue reflect.Value) {
	if sliceValue.IsNil() {
		return
	}

	for i := 0; i < sliceValue.Len(); i++ {
		processValue(sliceValue.Index(i))
	}
}

func processMap(mapValue reflect.Value) {
	if mapValue.IsNil() {
		return
	}

	for _, key := range mapValue.MapKeys() {
		value := mapValue.MapIndex(key)

		// For maps, we can only process if the value type is settable
		// We create a new value, process it, and set it back
		if value.Kind() == reflect.String {
			escaped := EscapeHtmlTags(value.String())
			mapValue.SetMapIndex(key, reflect.ValueOf(escaped))
		} else if !isSimpleType(value.Type()) {
			// For complex types, process in place if possible
			processValue(value)
		}
	}
}

// isSimpleType checks if a type is a simple type that doesn't need HTML escaping
func isSimpleType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true
	case reflect.String:
		return false // String needs processing
	case reflect.Ptr:
		return isSimpleType(t.Elem())
	default:
		// Check for common standard library types
		typeName := t.String()
		simpleTypes := []string{
			"time.Time", "time.Duration", "url.URL", "uuid.UUID",
			"json.Number", "big.Int", "big.Float", "big.Rat",
		}

		for _, simpleType := range simpleTypes {
			if typeName == simpleType {
				return true
			}
		}

		return false
	}
}
