package config

type TomlConfig struct {
	Build           Build          `toml:"build" json:"build"`
	Log             Log            `toml:"log" json:"log"`
	Token           Token          `toml:"token" json:"token"`
	Http            HttpConfig     `toml:"http" json:"http"`
	Https           HttpsConfig    `toml:"https" json:"https"`
	Auto_cert       AutoCert       `toml:"auto_cert" json:"auto_cert"`
	Redis           Redis          `toml:"redis" json:"redis"`
	Spr             Spr            `tome:"spr" json:"spr"`
	General_counter GeneralCounter `tome:"general_counter" json:"general_counter"`
	Db              DB             `toml:"db" json:"db"`
	Elastic_search  ElasticSearch  `toml:"elastic_search" json:"elastic_search"`
	Geo_ip          GeoIp          `toml:"geo_ip" json:"geo_ip"`
	Level_db        LevelDB        `toml:"level_db" json:"level_db"`
	Smtp            SMTP           `toml:"smtp" json:"smtp"`
	Sqlite          Sqlite         `toml:"sqlite" json:"sqlite"`
}

type Build struct {
	Mode string `toml:"mode" json:"mode"`
}
type Log struct {
	Level string `toml:"level" json:"level"`
}

type Token struct {
	Salt string `toml:"salt" json:"salt"`
}

type HttpConfig struct {
	Enable     bool `toml:"enable" json:"enable"`
	Port       int  `toml:"port" json:"port"`
	Keep_alive bool `toml:"keep_alive" json:"keep_alive"`
}

type HttpsConfig struct {
	Enable     bool   `toml:"enable" json:"enable"`
	Port       int    `toml:"port" json:"port"`
	Keep_alive bool   `toml:"keep_alive" json:"keep_alive"`
	Crt_path   string `toml:"crt_path" json:"crt_path"`
	Key_path   string `toml:"key_path" json:"key_path"`
	Html_dir   string `toml:"html_dir" json:"html_dir"`
}

type AutoCert struct {
	Enable         bool   `toml:"enable" json:"enable"`
	Check_interval int    `toml:"check_interval" json:"check_interval"`
	Crt_path       string `toml:"crt_path" json:"crt_path"`
	Init_download  bool   `toml:"init_download" json:"init_download"`
	Key_path       string `toml:"key_path" json:"key_path"`
	Url            string `toml:"url" json:"url"`
}

type Spr struct {
	Enable bool `toml:"enable" json:"enable"`
}

type GeneralCounter struct {
	Enable       bool   `toml:"enable" json:"enable"`
	Project_name string `toml:"project_name" json:"project_name"`
}

type Redis struct {
	Enable   bool   `toml:"enable" json:"enable"`
	Use_tls  bool   `toml:"use_tls" json:"use_tls"`
	Host     string `toml:"host" json:"host"`
	Port     int    `toml:"port" json:"port"`
	Username string `toml:"username" json:"username"`
	Password string `toml:"password" json:"password"`
	Prefix   string `toml:"prefix" json:"prefix"`
}

type DB struct {
	Enable   bool   `toml:"enable" json:"enable"`
	Host     string `toml:"host" json:"host"`
	Port     int    `toml:"port" json:"port"`
	Name     string `toml:"name" json:"name"`
	Username string `toml:"username" json:"username"`
	Password string `toml:"password" json:"password"`
}

type ElasticSearch struct {
	Enable   bool   `toml:"enable" json:"enable"`
	Host     string `toml:"host" json:"host"`
	Username string `toml:"username" json:"username"`
	Password string `toml:"password" json:"password"`
}

type GeoIp struct {
	Enable          bool   `toml:"enable" json:"enable"`
	Update_key      string `toml:"update_key" json:"update_key"`
	Dataset_folder  string `toml:"dataset_folder" json:"dataset_folder"`
	Dataset_version string `toml:"dataset_version" json:"dataset_version"`
}

type LevelDB struct {
	Enable bool   `toml:"enable" json:"enable"`
	Path   string `toml:"path" json:"path"`
}

type SMTP struct {
	Enable     bool   `toml:"enable" json:"enable"`
	From_email string `toml:"from_email" json:"from_email"`
	Host       string `toml:"host" json:"host"`
	Port       int    `toml:"port" json:"port"`
	Password   string `toml:"password" json:"password"`
	Username   string `toml:"username" json:"username"`
}

type Sqlite struct {
	Enable bool   `toml:"enable" json:"enable"`
	Path   string `toml:"path" json:"path"`
}
