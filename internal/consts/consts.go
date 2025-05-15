package consts

const (
	PORT           = "3169"     // Port to run the server on
	HTML_TMPL_PATH = "html-src" // Where we store our HTML templates
	STATIC_PATH    = "static"   // Where we store our static files (js, images, etc.)
	APP_TITLE      = "Yafti"
	// Time limit (in seconds) to check if a client is connected.
	// If it is not, shutdown the server.
	HEARTBEAT_SECONDS = 30
)
