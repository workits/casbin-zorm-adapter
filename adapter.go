package zormadapter

import (
	"context"
	"errors"
	"gitee.com/chunanyong/zorm"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

// table name
var casbinRuleEntityTableName = "casbin_rule"

// casbin rule struct
type casbinRuleEntity struct {
	zorm.EntityStruct
	ID    int64  `column:"id"`
	Ptype string `column:"ptype"`
	V0    string `column:"v0"`
	V1    string `column:"v1"`
	V2    string `column:"v2"`
	V3    string `column:"v3"`
	V4    string `column:"v4"`
	V5    string `column:"v5"`
}

// GetTableName get the table name of casbin-rule
func (entity *casbinRuleEntity) GetTableName() string {
	return casbinRuleEntityTableName
}

type Adapter struct {
	ctx context.Context
}

// NewAdapter get a zorm-adapter instance
func NewAdapter(db *zorm.DBDao, tableName ...string) (*Adapter, error) {
	// custom table name
	switch len(tableName) {
	case 0:
	case 1:
		casbinRuleEntityTableName = tableName[0]
	default:
		return nil, errors.New("too many parameters")
	}

	// get the zorm-adapter instance with datasource
	if ctx, err := db.BindContextDBConnection(context.Background()); err != nil {
		return nil, err
	} else {
		return &Adapter{
			ctx: ctx,
		}, nil
	}
}

// LoadPolicy loads all policies rules from database
func (a *Adapter) LoadPolicy(model model.Model) error {
	// create a slice to store the query result
	lines := make([]casbinRuleEntity, 0)
	finder := zorm.NewSelectFinder(casbinRuleEntityTableName)
	// to query result
	if err := zorm.Query(a.ctx, finder, &lines, nil); err != nil {
		return err
	}

	// load policy one by one
	for _, line := range lines {
		var p = []string{line.Ptype, line.V0, line.V1, line.V2, line.V3, line.V4, line.V5}

		idx := len(p) - 1
		for p[idx] == "" {
			idx--
		}
		idx += 1
		p = p[:idx]
		err := persist.LoadPolicyArray(p, model)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Adapter) savePolicyLine(ptype string, rule []string) *casbinRuleEntity {
	line := new(casbinRuleEntity)

	line.Ptype = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

// SavePolicy saves all policies rules to database
func (a *Adapter) SavePolicy(model model.Model) error {
	// using transaction
	_, err := zorm.Transaction(a.ctx, func(ctx context.Context) (interface{}, error) {
		// truncate the table
		finder := zorm.NewDeleteFinder(casbinRuleEntityTableName)
		if _, err := zorm.UpdateFinder(a.ctx, finder); err != nil {
			return nil, err
		}

		// batch insert the rules
		var lines []zorm.IEntityStruct
		flushEvery := 1000
		for ptype, ast := range model["p"] {
			for _, rule := range ast.Policy {
				lines = append(lines, a.savePolicyLine(ptype, rule))
				if len(lines) > flushEvery {
					if _, err := zorm.InsertSlice(a.ctx, lines); err != nil {
						return nil, err
					}
					lines = nil
				}
			}
		}
		for ptype, ast := range model["g"] {
			for _, rule := range ast.Policy {
				lines = append(lines, a.savePolicyLine(ptype, rule))
				if len(lines) > flushEvery {
					if _, err := zorm.InsertSlice(a.ctx, lines); err != nil {
						return nil, err
					}
					lines = nil
				}
			}
		}
		if len(lines) > 0 {
			if _, err := zorm.InsertSlice(a.ctx, lines); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	return err
}

// AddPolicy adds a policy rule to database
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	line := a.savePolicyLine(ptype, rule)
	_, err := zorm.Insert(a.ctx, line)
	return err
}

// RemovePolicy removes a policy rule from database
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	finder := zorm.NewDeleteFinder(casbinRuleEntityTableName)
	finder.Append("where ptype = ?", ptype)
	if len(rule) > 0 {
		finder.Append("and v0 = ?", rule[0])
	}
	if len(rule) > 1 {
		finder.Append("and v1 = ?", rule[1])
	}
	if len(rule) > 2 {
		finder.Append("and v2 = ?", rule[2])
	}
	if len(rule) > 3 {
		finder.Append("and v3 = ?", rule[3])
	}
	if len(rule) > 4 {
		finder.Append("and v4 = ?", rule[4])
	}
	if len(rule) > 5 {
		finder.Append("and v5 = ?", rule[5])
	}
	_, err := zorm.UpdateFinder(a.ctx, finder)
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from database
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	finder := zorm.NewDeleteFinder(casbinRuleEntityTableName)
	finder.Append("where ptype = ?", ptype)
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) && len(fieldValues[0-fieldIndex]) > 0 {
		finder.Append("and v0 = ?", fieldValues[0-fieldIndex])
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) && len(fieldValues[1-fieldIndex]) > 0 {
		finder.Append("and v1 = ?", fieldValues[1-fieldIndex])
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) && len(fieldValues[2-fieldIndex]) > 0 {
		finder.Append("and v2 = ?", fieldValues[2-fieldIndex])
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) && len(fieldValues[3-fieldIndex]) > 0 {
		finder.Append("and v3 = ?", fieldValues[3-fieldIndex])
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) && len(fieldValues[4-fieldIndex]) > 0 {
		finder.Append("and v4 = ?", fieldValues[4-fieldIndex])
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) && len(fieldValues[5-fieldIndex]) > 0 {
		finder.Append("and v5 = ?", fieldValues[5-fieldIndex])
	}
	_, err := zorm.UpdateFinder(a.ctx, finder)
	return err
}
