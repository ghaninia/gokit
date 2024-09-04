package response

import (
	"testing"
	"github.com/ghaninia/gokit/translation"

	"github.com/go-playground/validator/v10"
)

func TestNewValidationTranslator(t *testing.T) {
	trans := translation.NewTranslation(translation.Config{})
	got := newValidationTranslator(trans)
	if got.Translation != trans {
		t.Errorf("newValidationTranslator() = %v, want %v", got.Translation, trans)
	}
}

func TestValidation_translate(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name string
		v    validation
		args args
	}{
		{
			name: "Test case 1",
			v: validation{
				Translation: translation.NewTranslation(translation.Config{}),
			},
			args: args{
				err: nil,
			},
		},
		{
			name: "Test case 2",
			v: validation{
				Translation: translation.NewTranslation(translation.Config{}),
			},
			args: args{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v.translate(tt.args.err)
			if got != nil {
				t.Errorf("validation.translate() = %v, want %v", got, nil)
			}
		})
	}
}

func TestValidation_validators(t *testing.T) {
	request := struct {
		Name string `validate:"required"`
	}{
		Name: "",
	}

	validate := validator.New()
	err := validate.Struct(request)

	errs := newValidationTranslator(translation.NewTranslation(translation.Config{})).translate(err)

	if len(errs) == 0 {
		t.Errorf("validation.validators() = %v, want %v", errs, "not empty")
	}

	if errs["Name"][0] != "validation.required" {
		t.Errorf("validation.validators() = %v, want %v", errs["Name"][0], "Name is a required field")
	}
}

func TestValidation_validators2(t *testing.T) {
	request := struct {
		Name string `validate:"required"`
	}{
		Name: "John",
	}

	validate := validator.New()
	err := validate.Struct(request)

	errs := newValidationTranslator(translation.NewTranslation(translation.Config{})).translate(err)

	if len(errs) != 0 {
		t.Errorf("validation.validators() = %v, want %v", errs, "empty")
	}
}
