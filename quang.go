package quang

type Quang struct {
	evaluator evaluator_t
}

// Init the whole language. `query` is the expression provided by
// the user, for example: `size gt 0` which will be evaluated later.
func Init(query string) (*Quang, error) {
	l := createLexer(query)

	if err := l.lex(); err != nil {
		return nil, err
	}

	p := createParser(l.tokens)

	expr, err := p.parseExpression()

	if err != nil {
		return nil, err
	}

	evaluator := createEvaluator(expr)

	return &Quang{
		evaluator: evaluator,
	}, nil
}

// Build the set of available atoms.
// [You only need to do this once]
// Atoms works just like enums. You can say that an atom ":get" is 0
// So, everytime the user types ":get" in the query, i'll be substituted by 0,
// That way you make it easier to the user by specify a variable "kind", instead of typing 0, he types :get
func (q *Quang) SetupAtoms(atoms map[string]AtomType) *Quang {
	for key, value := range atoms {
		q.SetupAtom(key, value)
	}

	return q
}

// Build the set of available atoms.
// [You only need to do this once]
// Atoms works just like enums. You can say that an atom ":get" is 0
// So, everytime the user types ":get" in the query, i'll be substituted by 0,
// That way you make it easier to the user by specify a variable "kind", instead of typing 0, he types :get
func (q *Quang) SetupAtom(name string, value AtomType) *Quang {
	q.evaluator.setAtomValue(name, value)

	return q
}

// For each evaluation, you can provide different variable values.
// If, for example you want to do a query over a bunch of logs the user
// will provide the query, for example filtering by a specific user agent pattern
// then, for each log row, you can update the "agent" variable value to the current log row user agent
// so the, when the language lazy evaluate the "agent" variable the query will be applied to the current
// log row successfully
func (q *Quang) AddStringVar(name, value string) *Quang {
	q.evaluator.addStringVar(name, value)

	return q
}

// For each evaluation, you can provide different variable values.
// If, for example you want to do a query over a bunch of logs the user
// will provide the query, for example filtering by a specific user agent pattern
// then, for each log row, you can update the "agent" variable value to the current log row user agent
// so the, when the language lazy evaluate the "agent" variable the query will be applied to the current
// log row successfully
func (q *Quang) AddIntegerVar(name string, value IntegerType) *Quang {
	q.evaluator.addIntegerVar(name, value)

	return q
}

// For each evaluation, you can provide different variable values.
// If, for example you want to do a query over a bunch of logs the user
// will provide the query, for example filtering by a specific user agent pattern
// then, for each log row, you can update the "agent" variable value to the current log row user agent
// so the, when the language lazy evaluate the "agent" variable the query will be applied to the current
// log row successfully
func (q *Quang) AddFloatVar(name string, value FloatType) *Quang {
	q.evaluator.addFloatVar(name, value)

	return q
}

// For each evaluation, you can provide different variable values.
// If, for example you want to do a query over a bunch of logs the user
// will provide the query, for example filtering by a specific user agent pattern
// then, for each log row, you can update the "agent" variable value to the current log row user agent
// so the, when the language lazy evaluate the "agent" variable the query will be applied to the current
// log row successfully
func (q *Quang) AddBoolVar(name string, value bool) *Quang {
	q.evaluator.addBoolVar(name, value)

	return q
}

// For each evaluation, you can provide different variable values.
// If, for example you want to do a query over a bunch of logs the user
// will provide the query, for example filtering by a specific user agent pattern
// then, for each log row, you can update the "agent" variable value to the current log row user agent
// so the, when the language lazy evaluate the "agent" variable the query will be applied to the current
// log row successfully
func (q *Quang) AddAtomVar(name string, value AtomType) *Quang {
	q.evaluator.addAtomVar(name, value)

	return q
}

func (q Quang) Eval() (bool, error) {
	return q.evaluator.eval()
}
