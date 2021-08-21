package config

type BaseConfig interface {
	Validate() []error
}
