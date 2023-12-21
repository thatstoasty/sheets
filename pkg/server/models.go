package server

type Option struct {
	Name        string `gorm:"primaryKey"`
	Type        string
	Description string
}

type ClassFeature struct {
	Name     string `gorm:"primaryKey"`
	Class    string `gorm:"primaryKey"`
	Type     string
	SubClass string `gorm:"primaryKey"`
	Level    uint32
	Options  string
}

type Characteristic struct {
	Name        string `gorm:"primaryKey"`
	Description string
	Options     string
}

type Item struct {
	Name        string `gorm:"primaryKey"`
	Category    string
	Type        string
	Description string
	Properties  string
	Options     string
}

type Weapon struct {
	Name        string `gorm:"primaryKey"`
	Type        string
	Description string
	Properties  string
	Options     string
}

type Spell struct {
	Name        string `gorm:"primaryKey"`
	Level       uint16
	Description string
	Options     string
}

type Character struct {
	Name           string `gorm:"primaryKey"`
	Class          string
	HP             uint8
	Proficiency    uint8
	Strength       uint8
	Dexterity      uint8
	Constitution   uint8
	Intelligence   uint8
	Wisdom         uint8
	Charisma       uint8
	Race           string
	Feats          string
	Items          string
	Helmet         string
	Cloak          string
	Armor          string
	Jewelery1      string
	Jewelery2      string
	Jewelery3      string
	Boots          string
	Gloves         string
	MainHandWeapon string
	OffHandWeapon  string
}

type Inventory struct {
	Character string `gorm:"primaryKey"`
	Item      string
}

type FeatureChoices struct {
	Character string `gorm:"primaryKey"`
	Feature   string `gorm:"primaryKey"`
	Choice    string
}

type Config struct {
	Name  string
	Class string
	Race  string
	Feats string
	Items string
}

type Race struct {
	Name string `gorm:"primaryKey"`
}
