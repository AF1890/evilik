package domain

// Info représente une information pratique
type Info struct {
	Title   string
	Content string
	Icon    string // optionnel, pour le style
}

// PracticalInfo contient toutes les informations pratiques
type PracticalInfo struct {
	Venue         Info
	Access        Info
	Parking       Info
	Accommodation Info
	DressCode     Info
	MapURL        string
}

// GetDefaultPracticalInfo retourne les infos pratiques par défaut
func GetDefaultPracticalInfo() *PracticalInfo {
	return &PracticalInfo{
		Venue: Info{
			Title: "Le lieu",
			Content: `Le mariage civil aura lieu à la mairie de Cély.
La cérémonie laïque se tiendra ensuite chez nous, au 8 rue du Bois Beaudoin, à Cély-en-Bière.
La soirée se poursuivra à la Bergerie de Villiers-en-Bière.
Adresse : rue de la Bascule, 77190 Villiers-en-Bière.`,
			Icon: "📍",
		},
	}
}
