package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
)

type TasksValue struct {
	ID          int     `json:"valueId"`
	Value       *string `json:"value,omitempty"`
	OptionID    *int    `json:"optionId,omitempty"`
	ColumnID    int     `json:"columnId"`
	TaskID      int     `json:"taskId"`
	Action      string  `json:"action"`
	StringValue *string `json:"stringValue,omitempty"`
	NumberValue *int    `json:"numberValue,omitempty"`
	Type        string  `json:"type"`
}

func (tv *TasksValue) Validate() error {
	return validation.ValidateStruct(tv,
		validation.Field(&tv.Value, validation.By(validateValueField(tv))),
		validation.Field(&tv.OptionID, validation.By(validateOptionIDField(tv))),
		validation.Field(&tv.ColumnID, validation.Required, validation.Min(1)),
		validation.Field(&tv.TaskID, validation.Required, validation.Min(1)),
		validation.Field(&tv.Action, validation.Required, validation.In("create", "update", "delete", "get")),
		validation.Field(&tv.Type, validation.Required, validation.In("number", "string")),
	)
}

func (tv *TasksValue) ValidateUpdate() error {
	return validation.ValidateStruct(tv,
		validation.Field(&tv.Value, validation.By(validateValueField(tv))),
		validation.Field(&tv.OptionID, validation.By(validateOptionIDField(tv))),
		validation.Field(&tv.Type, validation.Required, validation.In("number", "string")),
	)
}

func (tv *TasksValue) ValidateGet() error {
	return validation.ValidateStruct(tv,
		validation.Field(&tv.ColumnID, validation.Required, validation.Min(1)),
		validation.Field(&tv.TaskID, validation.Required, validation.Min(1)),
		validation.Field(&tv.Action, validation.Required, validation.In("create", "update", "delete", "get")),
	)
}

// validateValueField ensures that the Value field is required if OptionID is nil.
func validateValueField(tv *TasksValue) validation.RuleFunc {
	return func(value interface{}) error {
		if tv.OptionID == nil && tv.Value == nil {
			return fmt.Errorf("Value is required when OptionID is not provided")
		}
		if tv.Value != nil && len(*tv.Value) == 0 {
			return fmt.Errorf("Value cannot be empty")
		}
		return nil
	}
}

// validateOptionIDField ensures that the OptionID field is valid if provided.
func validateOptionIDField(tv *TasksValue) validation.RuleFunc {
	return func(value interface{}) error {
		if tv.OptionID != nil && *tv.OptionID == 0 {
			return fmt.Errorf("OptionID must be greater than 0")
		}
		return nil
	}
}
