package config

type Config struct {
	FileUploadLimit int `yaml:"file_upload_limit,omitempty"`

	Sql PgsqlConfig `yaml:"sql"`
}

type PgsqlConfig struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
	DB string `yaml:"db"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	SSL  bool   `yaml:"ssl"`
}
