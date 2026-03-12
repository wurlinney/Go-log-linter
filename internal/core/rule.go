package core

// Rule описывает контракт проверки одного лога.
type Rule interface {
	ID() string
	Check(entry LogEntry) []Violation
}
