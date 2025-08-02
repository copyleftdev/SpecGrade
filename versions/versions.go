package versions

// VersionToSchemaURL maps OpenAPI versions to their schema URLs
var VersionToSchemaURL = map[string]string{
	"3.0.0": "https://spec.openapis.org/oas/3.0/schema/2019-04-02",
	"3.1.0": "https://spec.openapis.org/oas/3.1/schema/2022-10-07",
}

// IsValidVersion checks if the given version is supported
func IsValidVersion(version string) bool {
	_, exists := VersionToSchemaURL[version]
	return exists
}

// GetSchemaURL returns the schema URL for a given version
func GetSchemaURL(version string) (string, bool) {
	url, exists := VersionToSchemaURL[version]
	return url, exists
}
