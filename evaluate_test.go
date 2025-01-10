package rules

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/conekta/Conekta-Golang-Rules-Engine/parser"
)

func TestEvaluateAttrPathValue(t *testing.T) {
	res, err := Evaluate(`x.c.d > x.c.e`, map[string]interface{}{
		"x": map[string]interface{}{
			"c": map[string]interface{}{
				"d": 1.2,
				"e": 1.1,
			},
		},
	})
	require.NoError(t, err)
	require.True(t, res)
}

func TestEvaluateAttrPathValueNil(t *testing.T) {
	res, err := Evaluate(`x.c.d > x.c.e`, map[string]interface{}{
		"x": map[string]interface{}{
			"c": map[string]interface{}{
				"d": (*int)(nil),
				"e": 1,
			},
		},
	})
	var nestedError *parser.NestedError
	if errors.As(err, &nestedError) &&
		!errors.Is(nestedError.Err, parser.ErrEvalOperandMissing) {
		require.NoError(t, err)
	}
	require.False(t, res)
}

func TestEvaluateBasic(t *testing.T) {
	res, err := Evaluate(`x.c.d eq "abc"`, map[string]interface{}{
		"x": map[string]interface{}{
			"c": map[string]interface{}{
				"d": "abc",
			},
		},
	})
	require.NoError(t, err)
	require.True(t, res)
}

func TestEvaluateMissingField(t *testing.T) {
	res, err := Evaluate(`x.c.e eq "abc"`, map[string]interface{}{
		"x": map[string]interface{}{
			"c": map[string]interface{}{
				"d": "abc",
			},
		},
	})
	var nestedError *parser.NestedError
	if errors.As(err, &nestedError) &&
		!errors.Is(nestedError.Err, parser.ErrEvalOperandMissing) {
		require.NoError(t, err)
	}
	require.False(t, res)
}

func TestEvaluateMissingFieldWithNilToZeroValue(t *testing.T) {
	res, err := Evaluate(`x.c.e eq ""`, map[string]interface{}{
		"x": map[string]interface{}{
			"c": map[string]interface{}{
				"d": "abc",
			},
		},
	}, parser.WithNilToZeroValue())
	var nestedError *parser.NestedError
	if errors.As(err, &nestedError) &&
		!errors.Is(nestedError.Err, parser.ErrEvalOperandMissing) {
		require.NoError(t, err)
	}
	require.True(t, res)
}

func TestEvaluateMissingFieldEqEmpty(t *testing.T) {
	res, err := Evaluate(`x.c.e eq ""`, map[string]interface{}{
		"x": map[string]interface{}{
			"c": map[string]interface{}{
				"d": "abc",
			},
		},
	})
	var nestedError *parser.NestedError
	if errors.As(err, &nestedError) &&
		!errors.Is(nestedError.Err, parser.ErrEvalOperandMissing) {
		require.NoError(t, err)
	}
	require.False(t, res)
}

func TestEvaluateNilField(t *testing.T) {
	res, err := Evaluate(`x.c.e eq "abc"`, map[string]interface{}{
		"x": (*int)(nil),
		"e": nil,
	})
	var nestedError *parser.NestedError
	if errors.As(err, &nestedError) &&
		!errors.Is(nestedError.Err, parser.ErrEvalOperandMissing) {
		require.NoError(t, err)
	}
	require.False(t, res)
}

func TestSum(t *testing.T) {
	res, err := Evaluate(`SUM (x,y,z) eq 10`, map[string]interface{}{
		"x": 5,
		"y": 5,
		"z": 0,
	})
	require.NoError(t, err)
	require.True(t, res)
}

func TestMul(t *testing.T) {
	res, err := Evaluate(`MLP (x,y,a) > 24`, map[string]interface{}{
		"x": 5,
		"y": 5,
		"a": 5,
	})
	require.NoError(t, err)
	require.True(t, res)
}

func TestSubtractSuccessfully(t *testing.T) {
	t.Run("when subtract  (5-9) the result should be -7", func(t *testing.T) {
		res, err := Evaluate(`SUBTRACT (x.c,y) EQ -7`, map[string]interface{}{
			"x": map[string]interface{}{
				"c": 2,
			},
			"y": 9,
		})
		require.NoError(t, err)
		require.True(t, res)
	})
	t.Run("when subtract  (5-9) the result should be 7", func(t *testing.T) {
		res, err := Evaluate(`SUBTRACT (y,x.c.d) EQ 7`, map[string]interface{}{
			"x": map[string]interface{}{
				"c": map[string]interface{}{
					"d": 2,
				},
			},
			"y": 9,
		})
		require.NoError(t, err)
		require.True(t, res)
	})
	t.Run("when subtract  (-5- (-9)) the result should be 4", func(t *testing.T) {
		res, err := Evaluate(`SUBTRACT (x,y) EQ 4`, map[string]interface{}{
			"x": -5,
			"y": -9,
		})
		require.NoError(t, err)
		require.True(t, res)
	})
	t.Run("when subtract  (5-4) the result not should be 1", func(t *testing.T) {
		res, err := Evaluate(`SUBTRACT (x,y) EQ -1`, map[string]interface{}{
			"x": 5,
			"y": 4,
		})
		require.NoError(t, err)
		require.False(t, res)
	})
}

func TestSubtractUnSuccessfully(t *testing.T) {
	t.Run("when parameters are not enough should return an error", func(t *testing.T) {
		res, err := Evaluate(`SUBTRACT (y.c) EQ 1`, map[string]interface{}{
			"y": map[string]interface{}{
				"c": 4,
			},
		})
		require.Error(t, err)
		require.False(t, res)
	})
}

func TestDiv(t *testing.T) {
	t.Run("when division works successfully", func(t *testing.T) {
		res, err := Evaluate(`DIV (x,y) EQ 10`, map[string]interface{}{
			"x": 100,
			"y": 10,
		})
		require.NoError(t, err)
		require.True(t, res)
	})
	t.Run("when dividing by 0 should return 0", func(t *testing.T) {
		res, err := Evaluate(`DIV (x,y) EQ 0`, map[string]interface{}{
			"x": 100,
			"y": 0,
		})
		require.NoError(t, err)
		require.True(t, res)
	})
}

func TestCombiningMultipleOperators(t *testing.T) {
	t.Run("using 3 rules", func(t *testing.T) {
		res, err := Evaluate(`SUBTRACT (x,y) EQ 1 and SUBTRACT (x,y) EQ 1 and MLP (y,a) > 24`, map[string]interface{}{
			"x": 8,
			"y": 7,
			"a": 5,
		})
		require.NoError(t, err)
		require.True(t, res)
	})
	t.Run("using 5 rules", func(t *testing.T) {
		res, err := Evaluate(`x eq 8 and SUBTRACT (x,y) EQ 1 and SUBTRACT (x,y) EQ 1 and MLP (y,a) > 24 and MLP (y,a) eq 35`, map[string]interface{}{
			"x": 8,
			"y": 7,
			"a": 5,
		})
		require.NoError(t, err)
		require.True(t, res)
	})
}
