package decorators

const (
	ESC = "\x1b["

	RESET_ALL      = ESC + "0m"

    BOLD           = ESC + "1m"
    BOLD_OFF       = ESC + "21m"

    DIM            = ESC + "2m"
    DIM_OFF        = ESC + "22m"

    ITALICS        = ESC + "3m"
    ITALICS_OFF    = ESC + "23m"

    UNDERLINE      = ESC + "4m"
    UNDERLINE_OFF  = ESC + "24m"

    BLINK          = ESC + "5m"    // apparently 5 & 6 have the same effect
    BLINK_OFF      = ESC + "25m"

    REVERSE        = ESC + "7m"
    REVERSE_OFF    = ESC + "27m"

    HIDDEN         = ESC + "8m"
    HIDDEN_OFF     = ESC + "28m"

    STRIKE         = ESC + "9m"
    STRIKE_OFF     = ESC + "29m"
)

// another reason to have the names in CAPS is because Go is strange and public variables need to start with a capital letter.
