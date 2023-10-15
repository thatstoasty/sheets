package server

type Option struct {
	Name        string `gorm:"primaryKey"`
	Type        string
	Description string
}

type Class struct {
	Name string `gorm:"primaryKey"`
}

type ClassFeature struct {
	Name     string `gorm:"primaryKey"`
	Class    string `gorm:"primaryKey"`
	SubClass string `gorm:"primaryKey"`
	Level    uint32
	Options  string
}

type Race struct {
	Name    string `gorm:"primaryKey"`
	Options string
}

type Feat struct {
	Name        string `gorm:"primaryKey"`
	Description string
	Options     string
}

type Item struct {
	Name        string `gorm:"primaryKey"`
	Description string
	Options     string
}

type Spell struct {
	Name        string `gorm:"primaryKey"`
	Level       uint32
	Description string
	Options     string
}

type Character struct {
	Name         string `gorm:"primaryKey"`
	Class        string
	HP           string
	Proficiency  string
	Strength     string
	Dexterity    string
	Constitution string
	Intelligence string
	Wisdom       string
	Charisma     string
	Race         string
	Feats        string
	Items        string
}

type Config struct {
	Name  string
	Class string
	Race  string
	Feats string
	Items string
}
