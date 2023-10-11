package server

type Metadata struct {
}

type Option struct {
	Name string `gorm:"primaryKey"`
	Type string
}

type Action struct {
	Name string `gorm:"primaryKey"`
}

type BonusAction struct {
	Name string `gorm:"primaryKey"`
}

type Passive struct {
	Name string `gorm:"primaryKey"`
}

type Class struct {
	Name string `gorm:"primaryKey"`
}

type ClassFeature struct {
	Name    string `gorm:"primaryKey"`
	Class   string
	Level   uint32
	Options string
}

type Race struct {
	Name    string `gorm:"primaryKey"`
	Options string
}

type Feat struct {
	Name    string `gorm:"primaryKey"`
	Options string
}

type Item struct {
	Name    string `gorm:"primaryKey"`
	Options string
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

// type Spell struct {
// 	ID  		uint32 `gorm:"primaryKey,autoIncrement"`
// 	Name	 	string
// 	Description string
// 	Level		uint32
// }
