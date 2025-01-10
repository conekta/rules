package parser

type EvaluatorConfigOption func(*EvaluatorConfig)

type EvaluatorConfig struct {
	NilToZeroValue bool
}

func WithNilToZeroValue() EvaluatorConfigOption {
	return func(e *EvaluatorConfig) {
		e.NilToZeroValue = true
	}
}
