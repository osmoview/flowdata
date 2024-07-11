package flowdata

var DefaultLangKey = "en"

const (
	FieldTypeText     = "text"
	FieldTypeNumber   = "number"
	FieldTypeCheckbox = "checkbox"
	FieldTypeRadio    = "radio"
	FieldTypeButton   = "button"
	FieldTypeHidden   = "hidden"
	FieldTypeImage    = "image"
	FieldTypeSelect   = "select"
	FieldTypeURL      = "url"

	// sorry for this type, this should allows to add a pairs of fields, as key:value, to fill it to the map[string]float64
	FieldTypeMapNumbers = "map-numbers"
)

const (
	FieldLevelDefault = ""
	FieldLevelNewbie  = "newbie"
	FieldLevelMaster  = "master"
)

const (
	FieldHookCameras  = "cameras"
	FieldHookServices = "services"
	FieldHookDatabase = "database"
)

// Form it's just an slice of form fields
type Form []FormGroup

// FormGroup contains multiple fields and group it by some purpose
type FormGroup struct {
	// name of the form group
	Name string `yaml:"name" json:"name"`

	// map of localization labels, for example: {"en":"Threshold", "es":"Umbral", ...}, if it is null, then recommend use name as label
	Label L10n `yaml:"l10n,omitempty" json:"l10n,omitempty"`

	// list of fields of this form group
	Fields []FormField `yaml:"fields" json:"fields"`
}

// FormField describe how to required draw some input field on the frontend
type FormField struct {
	// field name, should be used as key for fill values in the data map
	Name string `yaml:"name" json:"name"`

	// type of field, like: text, number, etc
	Type string `yaml:"type,omitempty" json:"type,omitempty"`

	// default value
	Default any `yaml:"default,omitempty" json:"default,omitempty"`

	// map of localization labels, for example: {"en":"Threshold", "es":"Umbral", ...}, if it is null, then recommend use name as label
	Label L10n `yaml:"l10n,omitempty" json:"l10n,omitempty"`

	// help contains description of this field, it should be a some kind of hint/tooltip text
	Help L10n `yaml:"help,omitempty" json:"help,omitempty"`

	// level of field, to display or hide some fields on the frontend
	Level string `yaml:"level,omitempty" json:"level,omitempty"`

	// list of options for the field type select
	Options []FormFieldOption `yaml:"options,omitempty" json:"options,omitempty"`

	// some hook to execute after the form is built
	Hook string `yaml:"hook,omitempty" json:"-"`
}

type FormFieldOption struct {
	Label string `yaml:"label" json:"label"`
	Value any    `yaml:"value" json:"value"`
}

type L10n map[string]string

// Get should extract from the l10n map required value by key
// if not exist, try to get default language, if not then returns a key
func (m L10n) Get(key string) string {
	if v, ok := m[key]; ok {
		return v
	}

	if v, ok := m[DefaultLangKey]; ok {
		return v
	}

	return key
}
