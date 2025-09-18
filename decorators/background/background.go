package background

const (
	ESC = "\x1b["

	DEFAULT_COLOR  = ESC + "49m"
	BLACK          = ESC + "40m"
	RED            = ESC + "41m"
	GREEN          = ESC + "42m"

	YELLOW         = ESC + "43m"
	BLUE           = ESC + "44m"
	MAGENTA        = ESC + "45m"
	CYAN           = ESC + "46m"

	LIGHT_GRAY     = ESC + "47m"
	DARK_GRAY      = ESC + "100m"
	LIGHT_RED      = ESC + "101m"
	LIGHT_GREEN    = ESC + "102m"

	LIGHT_YELLOW   = ESC + "103m"
	LIGHT_BLUE     = ESC + "104m"
	LIGHT_MAGENTA  = ESC + "105m"
	LIGHT_CYAN     = ESC + "106m"

	WHITE          = ESC + "107m"
)
