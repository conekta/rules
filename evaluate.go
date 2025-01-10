package rules

import "github.com/conekta/Conekta-Golang-Rules-Engine/parser"

func Evaluate(rule string, items map[string]interface{}, opts ...parser.EvaluatorConfigOption) (bool, error) {
	ev, err := parser.NewEvaluator(rule, opts...)
	if err != nil {
		return false, err
	}
	res, err := ev.Process(items)
	if err != nil {
		return false, err
	}
	if err = ev.LastDebugErr(); err != nil {
		return res, err
	}
	return res, nil
}
