package main

var raceData = []Characteristic{
	{Name: "Default", Type: "Race", Options: "Grapple|Shove|Dash|Disengage|Dodge|Use Object|Equip Item|Escape Grapple|Hide|Search|Ready Action"},
	{Name: "Dwarf", Type: "Race", Options: "Darkvision|Dwarven Resilience|Stonecunning|Dwarven Combat Training|Tool Proficiency"},
	{Name: "Elf", Type: "Race", Options: "Darkvision|Keen Senses|Fey Ancestry|Trance"},
	{Name: "Halfling", Type: "Race", Options: "Lucky|Brave|Halfling Nimbleness"},
	{Name: "Human", Type: "Race", Options: "None"},
	{Name: "Dragonborn", Type: "Race", Options: "None"},
	{Name: "Gnome", Type: "Race", Options: "Darkvision|Gnome Cunning"},
	{Name: "Half-Elf", Type: "Race", Options: "Darkvision|Fey Ancestry"},
	{Name: "Half-Orc", Type: "Race", Options: "Darkvision|Relentless Endurance|Savage Attacks"},
	{Name: "Tiefling", Type: "Race", Options: "Darkvision|Hellish Resistance|Infernal Legacy"},
}

var featData = []Characteristic{
	{Name: "Ability Score Improvement", Type: "Feat", Options: "Ability Score Improvement"},
	{Name: "Alert", Type: "Feat", Options: "Alert: Bonus Initiative|Alert: Surprise Immunity|Alert: Prevent Hidden Advantage"},
	{Name: "Athlete", Type: "Feat", Options: "Physical Training|Athlete: Kinesthetic Awareness"},
	{Name: "Actor", Type: "Feat", Options: "Charisma Training|Actor: Mimicry"},
	{Name: "Charger", Type: "Feat", Options: "Charger: Bull Rush"},
	{Name: "Crossbow Expert", Type: "Feat", Options: "Crossbow Expert: Crossbow Mastery|Crossbow Expert: Quick Load"},
	{Name: "Defensive Duelist", Type: "Feat", Options: "Defensive Duelist"},
	{Name: "Dual Wielder", Type: "Feat", Options: "Dual Wielder"},
	{Name: "Dungeon Delver", Type: "Feat", Options: "Dungeon Delver"},
	{Name: "Durable", Type: "Feat", Options: "Constitution Training"},
	{Name: "Elemental Adept", Type: "Feat", Options: "Elemental Adept"},
	{Name: "Grappler", Type: "Feat", Options: "Grappler: Grapple Mastery|Grappler: Pin Target"},
	{Name: "Great Weapon Master", Type: "Feat", Options: "Great Weapon Master: Bonus Attack|Great Weapon Master: Fierce Attack"},
	{Name: "Healer", Type: "Feat", Options: "Healer: Stabilizer|Healer: Greater Healer's Kit"},
	{Name: "Heavily Armored", Type: "Feat", Options: "Physical Training|Heavy Armor Proficiency"},
	{Name: "Heavy Armor Master", Type: "Feat", Options: "Strength Training|Heavily Armor Master: Physical Resistance"},
	{Name: "Inspiring Leader", Type: "Feat", Options: "Inspiring Leader"},
	{Name: "Keen Mind", Type: "Feat", Options: "Intelligence Training|Keen Mind: Eidetic Memory"},
	{Name: "Lightly Armored", Type: "Feat", Options: "Physical Training|Light Armor Proficiency"},
	{Name: "Linguist", Type: "Feat", Options: "Intelligence Training|Linguist: Polyglot|Linguist: Cipher Expert"},
	{Name: "Lucky", Type: "Feat", Options: "Lucky"},
	{Name: "Mage Slayer", Type: "Feat", Options: "Mage Slayer: Punish Spellcasting|Mage Slayer: Disruptor"},
	{Name: "Magic Initiate", Type: "Feat", Options: "Magic Initiate"},
	{Name: "Martial Adept", Type: "Feat", Options: "Martial Adept"},
	{Name: "Medium Armor Master", Type: "Feat", Options: "Medium Armor Master"},
	{Name: "Mobile", Type: "Feat", Options: "Mobile: Increased Speed|Mobile: Precision Movement"},
	{Name: "Moderately Armored", Type: "Feat", Options: "Physical Training|Medium Armor Proficiency|Shield Proficiency"},
	{Name: "Mounted Combatant", Type: "Feat", Options: "Mounted Combatant"},
	{Name: "Observant", Type: "Feat", Options: "Mental Training|Observant: Lip Reading|Observant: Improved Passive Perception/Investigation"},
	{Name: "Polearm Master", Type: "Feat", Options: "Polearm Master: Pommel Strike|Polearm Master: Pre-emptive Opportunity Attack"},
	{Name: "Resilient", Type: "Feat", Options: "Resilient"},
	{Name: "Ritual Caster", Type: "Feat", Options: "Ritual Caster: Bonus Ritual Spells|Ritual Caster: Transcribe Rituals"},
	{Name: "Savage Attacker", Type: "Feat", Options: "Savage Attacker"},
	{Name: "Sentinel", Type: "Feat", Options: "Sentinel: Enhanced Opportunity Attack|Sentinel: Guardian"},
	{Name: "Sharpshooter", Type: "Feat", Options: "Sharpshooter: Precise Aim|Sharpshooter: Overdraw"},
	{Name: "Shield Master", Type: "Feat", Options: "Shield Master: Shove|Shield Master: Improved Dexterity Saving Throws|Shield Master: Nullify Damage"},
	{Name: "Skilled", Type: "Feat", Options: "Skilled: Bonus Proficiencies"},
	{Name: "Skulker", Type: "Feat", Options: "Skulker"},
	{Name: "Spell Sniper", Type: "Feat", Options: "Spell Sniper: Precise Aim|Spell Sniper: Bonus Cantrip"},
	{Name: "Tavern Brawler", Type: "Feat", Options: "Physical Training|Unarmed Proficiency|Improvised Weapon Proficiency|Tavern Brawler: Improved Unarmed Damage|Tavern Brawler: Quick Grapple"},
	{Name: "Tough", Type: "Feat", Options: "Tough: Bonus Health|Tough: Healthy"},
	{Name: "War Caster", Type: "Feat", Options: "War Caster: Battlemage|War Caster: Opportunity Attack Spell"},
	{Name: "Weapon Master", Type: "Feat", Options: "Physical Training|Weapon Master: Bonus Proficiencies"},
}

var characteristicsData = [][]Characteristic{
	raceData,
	featData,
}
