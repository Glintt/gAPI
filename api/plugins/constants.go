package plugins

// Go plugin extension
var PLUGIN_EXTENSION = ".so"

// Plugins location
var PLUGINS_LOCATION = "./custom_plugins"

// Before request plugins
var BEFORE_REQUEST_PLUGINS_INNER_FOLDER_NAME = "BeforeRequest"

// Lookup for before request plugin
var BeforeRequestEntryPluginLookup = "BeforeRequestEntryPlugin"

var SERVICE_NAME = "/plugins"

var PLUGINS_TYPES = []string{BEFORE_REQUEST_PLUGINS_INNER_FOLDER_NAME}

const (
	PLUGINS_COLLECTION = "plugins"
)
