package linter

import (
	"golang.org/x/tools/go/packages"
)

const (
	PresetBugs        = "bugs"        // Related to bugs detection.
	PresetComment     = "comment"     // Related to comments analysis.
	PresetComplexity  = "complexity"  // Related to code complexity analysis.
	PresetError       = "error"       // Related to error handling analysis.
	PresetFormatting  = "format"      // Related to code formatting.
	PresetImport      = "import"      // Related to imports analysis.
	PresetMetaLinter  = "metalinter"  // Related to linter that contains multiple rules or multiple linters.
	PresetModule      = "module"      // Related to Go modules analysis.
	PresetPerformance = "performance" // Related to performance.
	PresetSQL         = "sql"         // Related to SQL.
	PresetStyle       = "style"       // Related to coding style.
	PresetTest        = "test"        // Related to the analysis of the code of the tests.
	PresetUnused      = "unused"      // Related to the detection of unused code.
)

type Deprecation struct {
	Since       string
	Message     string
	Replacement string
}

type Config struct {
	Linter           Linter
	EnabledByDefault bool

	LoadMode packages.LoadMode

	InPresets        []string
	AlternativeNames []string

	OriginalURL     string // URL of original (not forked) repo, needed for autogenerated README
	CanAutoFix      bool
	IsSlow          bool
	DoesChangeTypes bool

	Since       string
	Deprecation *Deprecation
}

func (lc *Config) ConsiderSlow() *Config {
	lc.IsSlow = true
	return lc
}

func (lc *Config) IsSlowLinter() bool {
	return lc.IsSlow
}

func (lc *Config) WithLoadFiles() *Config {
	lc.LoadMode |= packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles
	return lc
}

func (lc *Config) WithLoadForGoAnalysis() *Config {
	lc = lc.WithLoadFiles()
	lc.LoadMode |= packages.NeedImports | packages.NeedDeps | packages.NeedExportsFile | packages.NeedTypesSizes
	lc.IsSlow = true
	return lc
}

func (lc *Config) WithPresets(presets ...string) *Config {
	lc.InPresets = presets
	return lc
}

func (lc *Config) WithURL(url string) *Config {
	lc.OriginalURL = url
	return lc
}

func (lc *Config) WithAlternativeNames(names ...string) *Config {
	lc.AlternativeNames = names
	return lc
}

func (lc *Config) WithAutoFix() *Config {
	lc.CanAutoFix = true
	return lc
}

func (lc *Config) WithChangeTypes() *Config {
	lc.DoesChangeTypes = true
	return lc
}

func (lc *Config) WithSince(version string) *Config {
	lc.Since = version
	return lc
}

func (lc *Config) Deprecated(message, version, replacement string) *Config {
	lc.Deprecation = &Deprecation{
		Since:       version,
		Message:     message,
		Replacement: replacement,
	}
	return lc
}

func (lc *Config) IsDeprecated() bool {
	return lc.Deprecation != nil
}

func (lc *Config) AllNames() []string {
	return append([]string{lc.Name()}, lc.AlternativeNames...)
}

func (lc *Config) Name() string {
	return lc.Linter.Name()
}

func NewConfig(linter Linter) *Config {
	lc := &Config{
		Linter: linter,
	}
	return lc.WithLoadFiles()
}
