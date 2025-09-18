package foreground

const (
	ESC = "\x1b["

	DEFAULT_COLOR  = ESC + "39m"

	BLACK          = ESC + "30m"

	RED            = ESC + "31m"
	GREEN          = ESC + "32m"
	YELLOW         = ESC + "33m"
	BLUE           = ESC + "34m"
	MAGENTA        = ESC + "35m"
	CYAN           = ESC + "36m"
	LIGHT_GRAY     = ESC + "37m"

	DARK_GRAY      = ESC + "90m"
	LIGHT_RED      = ESC + "91m"
	LIGHT_GREEN    = ESC + "92m"
	LIGHT_YELLOW   = ESC + "93m"
	LIGHT_BLUE     = ESC + "94m"
	LIGHT_MAGENTA  = ESC + "95m"
	LIGHT_CYAN     = ESC + "96m"

	WHITE          = ESC + "97m"
)
