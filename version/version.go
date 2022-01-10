package version

var (
	Version = "v1.0.0"
)

func SemVersion() string {
	return "semver:" + Version
}
