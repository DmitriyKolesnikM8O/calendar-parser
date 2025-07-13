package connection

type ColorInfo struct {
	Name        string
	HexCode     string
	Description string
}

var Colors = map[string]ColorInfo{
	"6": {
		Name:        "red",
		HexCode:     "#E3450B",
		Description: "important events",
	},
	"2": {
		Name:        "green",
		HexCode:     "#17D427",
		Description: "sleep",
	},
	"": {
		Name:        "blue",
		HexCode:     "#1634DB",
		Description: "useful activities",
	},
	"3": {
		Name:        "violet",
		HexCode:     "#9B2AC9",
		Description: "useless activities",
	},
	"4": {
		Name:        "flamingo",
		HexCode:     "#DE8157",
		Description: "cooking and eating",
	},
	"5": {
		Name:        "yellow",
		HexCode:     "#EFD10F",
		Description: "time instead",
	},
	"8": {
		Name:        "grey",
		HexCode:     "#7D7877",
		Description: "trains",
	},
	"11": {
		Name:        "bright red",
		HexCode:     "#FF0000",
		Description: "another option for important events",
	},
	"7": {
		Name:        "bright blue",
		HexCode:     "#031D9C",
		Description: "anouther useful activities",
	},
}

func GetColorHex(nameColor string) string {
	for digit, _ := range Colors {
		if Colors[digit].Name == nameColor {
			return Colors[digit].HexCode
		}
	}

	return "#CCCCCC"
}

func GetColorDescription(nameColor string) string {
	for digit, _ := range Colors {
		if Colors[digit].Name == nameColor {
			return Colors[digit].Description
		}
	}

	return "Unknown"
}
