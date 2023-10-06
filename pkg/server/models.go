package server

type Metadata struct {
}

type Option struct {
	ID   uint32 `gorm:"primaryKey,autoIncrement"`
	Name string
	Type string
}

type Action struct {
	ID   uint32 `gorm:"primaryKey,autoIncrement"`
	Name string
}

type BonusAction struct {
	ID   uint32 `gorm:"primaryKey,autoIncrement"`
	Name string
}

type Passive struct {
	ID   uint32 `gorm:"primaryKey,autoIncrement"`
	Name string
}

type Class struct {
	ID   uint32 `gorm:"primaryKey,autoIncrement"`
	Name string
}

type ClassFeature struct {
	ID      uint32 `gorm:"primaryKey,autoIncrement"`
	Name    string
	Class   string
	Level   uint32
	Options string
}

type Race struct {
	ID      uint32 `gorm:"primaryKey,autoIncrement"`
	Name    string
	Options string
}

type Feat struct {
	ID      uint32 `gorm:"primaryKey,autoIncrement"`
	Name    string
	Options string
}

type Item struct {
	ID      uint32 `gorm:"primaryKey,autoIncrement"`
	Name    string
	Options string
}

type ClassLevel struct {
	Class string
	Level uint32
}

type Character struct {
	ID    uint32 `gorm:"primaryKey,autoIncrement"`
	Name  string
	Class string
	Race  string
	Feats string
	Items string
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
