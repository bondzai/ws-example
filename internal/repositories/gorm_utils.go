package repositories

import (
	"errors"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// UpdateAssociations inspects model's fields for slice associations (using GORM tags)
// and replaces the corresponding associations using db.Model(model).Association(fieldName).Replace(...).
// The model must be a pointer to a struct.
func UpdateAssociations(db *gorm.DB, model interface{}) error {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		return errors.New("model must be a non-nil pointer to a struct")
	}
	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		// Only process exported fields (PkgPath is empty for exported fields)
		if fieldType.PkgPath != "" {
			continue
		}
		// If the field is a slice, check if it represents an association.
		if field.Kind() == reflect.Slice {
			// Get the gorm tag for the field.
			gormTag := fieldType.Tag.Get("gorm")
			// If the tag includes "many2many" or "foreignkey", we consider it an association.
			if gormTag != "" && (strings.Contains(gormTag, "many2many") || strings.Contains(gormTag, "foreignkey")) {
				// Only update the association if the slice is non-nil.
				if !field.IsNil() {
					assocName := fieldType.Name // by default, use the struct field name as association name.
					// Replace the association with the provided slice value.
					if err := db.Model(model).Association(assocName).Replace(field.Interface()); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
