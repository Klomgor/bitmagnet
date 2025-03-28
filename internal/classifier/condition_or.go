package classifier

const orName = "or"

type orCondition struct{}

func (orCondition) name() string {
	return orName
}

var orConditionSpec = payloadSingleKeyValue[[]any]{
	key: orName,
	valueSpec: payloadMustSucceed[[]any]{payloadList[any]{
		itemSpec: payloadGeneric[any]{
			jsonSchema: map[string]any{
				"$ref": "#/definitions/condition",
			},
		},
		description: "A condition that is satisfied if any of the conditions in a list are satisfied",
	}},
}

func (orCondition) compileCondition(ctx compilerContext) (condition, error) {
	rawConds, err := orConditionSpec.Unmarshal(ctx)
	if err != nil {
		return condition{}, err
	}

	conds := make([]condition, len(rawConds))

	for i, rawCond := range rawConds {
		cond, err := ctx.compileCondition(ctx.child(numericPathPart(i), rawCond))
		if err != nil {
			return condition{}, err
		}

		conds[i] = cond
	}

	return condition{func(ctx executionContext) (bool, error) {
		for _, c := range conds {
			if result, err := c.check(ctx); err != nil {
				return false, err
			} else if result {
				return true, nil
			}
		}
		return false, nil
	}}, nil
}

func (orCondition) JSONSchema() JSONSchema {
	return orConditionSpec.JSONSchema()
}
