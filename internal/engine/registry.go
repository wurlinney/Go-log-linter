package engine

import (
	"github.com/wurlinney/go-log-linter/internal/config"
	"github.com/wurlinney/go-log-linter/internal/core"
	"github.com/wurlinney/go-log-linter/internal/extractors"
	"github.com/wurlinney/go-log-linter/internal/rules"
)

// NewDefaultEngine создаёт движок с конфигурацией по умолчанию.
func NewDefaultEngine(customRules []core.Rule) *Engine {
	return NewEngineWithConfig(config.Default(), customRules)
}

// NewEngineWithConfig создаёт движок на основе конфигурации.
func NewEngineWithConfig(cfg config.Config, customRules []core.Rule) *Engine {
	extrs := []extractors.Extractor{
		extractors.NewSlogExtractor(),
		extractors.NewZapExtractor(),
	}

	// Базовый набор правил с учётом конфигурации.
	var baseRules []core.Rule
	if cfg.Rules.Lowercase {
		baseRules = append(baseRules, rules.NewLowercaseRule())
	}
	if cfg.Rules.English {
		baseRules = append(baseRules, rules.NewEnglishRule())
	}
	if cfg.Rules.Symbols {
		baseRules = append(baseRules, rules.NewSymbolsRule())
	}
	if cfg.Rules.Sensitive {
		baseRules = append(baseRules, rules.NewSensitiveRuleWithKeywords(cfg.Sensitive.Keywords))
	}

	allRules := append(baseRules, customRules...)
	return NewEngine(extrs, allRules)
}
