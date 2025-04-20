package quang_test

import (
	"testing"

	"github.com/marcos-venicius/quang"
	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
	q, err := quang.Init("size gt 0 and method eq :get and status eq 200")

	atoms := map[string]quang.AtomType{
		":get": 0,
	}

	q.SetupAtoms(atoms)

	assert.Nil(t, err)

	type test_case_t struct {
		size   quang.IntegerType
		status quang.IntegerType
		method quang.AtomType
		result bool
	}

	tests := []test_case_t{
		{
			size:   0,
			method: atoms[":get"],
			status: 200,
			result: false,
		},
		{
			size:   1,
			method: atoms[":get"],
			status: 200,
			result: true,
		},
		{
			size:   1,
			method: 1,
			status: 200,
			result: false,
		},
		{
			size:   1,
			method: atoms[":get"],
			status: 204,
			result: false,
		},
	}

	for _, test := range tests {
		q.AddAtomVar("method", test.method).
			AddIntegerVar("size", test.size).
			AddIntegerVar("status", test.status)

		r, err := q.Eval()

		assert.Nil(t, err)
		assert.Equal(t, test.result, r)
	}
}
