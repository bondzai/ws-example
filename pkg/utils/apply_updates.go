package utils

import (
	"errors"
	"reflect"
	"strings"
)

// findMatchingField searches the destination struct for a field that matches the given source field name.
// It first compares field names case‑insensitively, then looks at JSON tags. If no match is found,
// and if the source field name ends with "id" (but isn’t exactly "id"), it heuristically returns the field "Id".
func findMatchingField(destElem reflect.Value, srcFieldName string) (reflect.Value, bool) {
	destType := destElem.Type()
	lowerSrc := strings.ToLower(srcFieldName)
	// Look for a field with matching name or JSON tag.
	for i := 0; i < destElem.NumField(); i++ {
		field := destType.Field(i)
		if strings.ToLower(field.Name) == lowerSrc {
			return destElem.Field(i), true
		}
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			tagName := strings.Split(jsonTag, ",")[0]
			if strings.ToLower(tagName) == lowerSrc {
				return destElem.Field(i), true
			}
		}
	}
	// Heuristic: if srcFieldName ends with "id" and is not exactly "id", try the "Id" field.
	if lowerSrc != "id" && strings.HasSuffix(lowerSrc, "id") {
		if f := destElem.FieldByName("Id"); f.IsValid() && f.CanSet() {
			return f, true
		}
	}
	return reflect.Value{}, false
}

// replaceSliceField creates a new slice for the destination by converting each element
// from the source slice to the destination element type using ApplyUpdates recursively.
func replaceSliceField(destField, srcField reflect.Value) error {
	srcSlice := srcField.Elem()
	// If the slice types are directly assignable, use them.
	if srcSlice.Type().AssignableTo(destField.Type()) {
		destField.Set(srcSlice)
		return nil
	}
	// Otherwise, create a new slice of the destination type.
	newSlice := reflect.MakeSlice(destField.Type(), srcSlice.Len(), srcSlice.Len())
	destElemType := destField.Type().Elem()
	for i := 0; i < srcSlice.Len(); i++ {
		srcElem := srcSlice.Index(i)
		// Allocate a new destination element.
		destElem := reflect.New(destElemType).Elem()
		var srcElemPtr reflect.Value
		if srcElem.Kind() != reflect.Ptr && srcElem.CanAddr() {
			srcElemPtr = srcElem.Addr()
		} else {
			srcElemPtr = srcElem
		}
		// Recursively convert/update the element.
		if err := ApplyUpdates(destElem.Addr().Interface(), srcElemPtr.Interface()); err != nil {
			return err
		}
		newSlice.Index(i).Set(destElem)
	}
	destField.Set(newSlice)
	return nil
}

// ApplyUpdates updates non‑nil fields from src into dest.
// Both dest and src must be non‑nil pointers to structs.
// Scalar fields are assigned directly; nested structs are updated recursively;
// and slice fields are replaced by building a new slice from the src slice.
func ApplyUpdates(dest, src interface{}) error {
	destVal := reflect.ValueOf(dest)
	srcVal := reflect.ValueOf(src)

	// Validate that dest and src are non‑nil pointers to structs.
	if destVal.Kind() != reflect.Ptr || destVal.IsNil() || destVal.Elem().Kind() != reflect.Struct {
		return errors.New("destination must be a non-nil pointer to a struct")
	}
	if srcVal.Kind() != reflect.Ptr || srcVal.IsNil() || srcVal.Elem().Kind() != reflect.Struct {
		return errors.New("source must be a non-nil pointer to a struct")
	}

	destElem := destVal.Elem()
	srcElem := srcVal.Elem()
	srcType := srcElem.Type()

	// Iterate over each exported field in src.
	for i := 0; i < srcElem.NumField(); i++ {
		fieldInfo := srcType.Field(i)
		// Skip unexported fields.
		if fieldInfo.PkgPath != "" {
			continue
		}
		srcField := srcElem.Field(i)
		// Dynamically find matching field in dest.
		destField, found := findMatchingField(destElem, fieldInfo.Name)
		if !found || !destField.IsValid() || !destField.CanSet() {
			continue
		}
		// Skip if src field is nil.
		if srcField.Kind() != reflect.Ptr || srcField.IsNil() {
			continue
		}
		// Determine the underlying kind.
		switch srcField.Elem().Kind() {
		case reflect.Slice:
			if err := replaceSliceField(destField, srcField); err != nil {
				return err
			}
		case reflect.Struct:
			// For nested structs, update recursively.
			if destField.Kind() == reflect.Ptr {
				if destField.IsNil() {
					destField.Set(reflect.New(destField.Type().Elem()))
				}
				if err := ApplyUpdates(destField.Interface(), srcField.Interface()); err != nil {
					return err
				}
			} else if destField.Kind() == reflect.Struct {
				if err := ApplyUpdates(destField.Addr().Interface(), srcField.Interface()); err != nil {
					return err
				}
			} else {
				destField.Set(srcField.Elem())
			}
		default:
			// For scalar fields, assign the value.
			if destField.Kind() == reflect.Ptr {
				destField.Set(srcField)
			} else {
				destField.Set(srcField.Elem())
			}
		}
	}
	return nil
}
