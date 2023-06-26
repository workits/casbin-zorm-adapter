package zormadapter

import (
	"gitee.com/chunanyong/zorm"
	"github.com/casbin/casbin/v2"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"testing"
)

var dbDao, _ = zorm.NewDBDao(&zorm.DataSourceConfig{
	DSN:                   "root:password@tcp(127.0.0.1:3306)/casbin?charset=utf8&parseTime=true&loc=Local",
	DriverName:            "mysql",
	Dialect:               "mysql",
	SlowSQLMillis:         0,
	MaxOpenConns:          0,
	MaxIdleConns:          0,
	ConnMaxLifetimeSecond: 0,
	DefaultTxOptions:      nil,
})

func testGetPolicy(t *testing.T, e *casbin.Enforcer, res [][]string) {
	t.Helper()
	myRes := e.GetPolicy()
	log.Print("Policy: ", myRes)

	m := make(map[string]bool, len(res))
	for _, value := range res {
		key := strings.Join(value, ",")
		m[key] = true
	}

	for _, value := range myRes {
		key := strings.Join(value, ",")
		if !m[key] {
			t.Error("Policy: ", myRes, ", supposed to be ", res)
			break
		}
	}
}

func initPolicy(t *testing.T) {
	var err error

	// Because the DB is empty at first,
	// so we need to load the policy from the file adapter (.CSV) first.
	e, err := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")
	if err != nil {
		panic(err)
	}

	a, err := NewAdapter(dbDao)
	if err != nil {
		panic(err)
	}
	// This is a trick to save the current policy to the DB.
	// We can't call e.SavePolicy() because the adapter in the enforcer is still the file adapter.
	// The current policy means the policy in the Casbin enforcer (aka in memory).
	err = a.SavePolicy(e.GetModel())
	if err != nil {
		panic(err)
	}

	// Clear the current policy.
	e.ClearPolicy()
	testGetPolicy(t, e, [][]string{})

	// Load the policy from DB.
	err = a.LoadPolicy(e.GetModel())
	if err != nil {
		panic(err)
	}
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
}

func testSaveLoad(t *testing.T) {
	// Initialize some policy in DB.
	initPolicy(t)
	// Note: you don't need to look at the above code
	// if you already have a working DB with policy inside.

	// Now the DB has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	a, err := NewAdapter(dbDao)
	if err != nil {
		t.Fatal(err)
	}
	e, err := casbin.NewEnforcer("examples/rbac_model.conf", a)
	if err != nil {
		t.Fatal(err)
	}
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
}

func testAutoSave(t *testing.T) {
	// Initialize some policy in DB.
	initPolicy(t)
	// Note: you don't need to look at the above code
	// if you already have a working DB with policy inside.
	var err error
	// Now the DB has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	a, err := NewAdapter(dbDao)
	if err != nil {
		t.Fatal(err)
	}
	e, err := casbin.NewEnforcer("examples/rbac_model.conf", a)
	if err != nil {
		t.Fatal(err)
	}

	// AutoSave is enabled by default.
	// Now we disable it.
	e.EnableAutoSave(false)

	logErr := func(action string) {
		if err != nil {
			t.Fatalf("test action[%s] failed, err: %v", action, err)
		}
	}

	// Because AutoSave is disabled, the policy change only affects the policy in Casbin enforcer,
	// it doesn't affect the policy in the storage.
	_, err = e.AddPolicy("alice", "data1", "write")
	logErr("AddPolicy")
	// Reload the policy from the storage to see the effect.
	err = e.LoadPolicy()
	logErr("LoadPolicy")
	// This is still the original policy.
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})

	// Now we enable the AutoSave.
	e.EnableAutoSave(true)

	// Because AutoSave is enabled, the policy change not only affects the policy in Casbin enforcer,
	// but also affects the policy in the storage.
	_, err = e.AddPolicy("alice", "data1", "write")
	logErr("AddPolicy2")
	// Reload the policy from the storage to see the effect.
	err = e.LoadPolicy()
	logErr("LoadPolicy2")
	// The policy has a new rule: {"alice", "data1", "write"}.
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}, {"alice", "data1", "write"}})

	// Remove the added rule.
	_, err = e.RemovePolicy("alice", "data1", "write")
	logErr("RemovePolicy")
	err = e.LoadPolicy()
	logErr("LoadPolicy3")
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})

	// Remove "data2_admin" related policy rules via a filter.
	// Two rules: {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"} are deleted.
	_, err = e.RemoveFilteredPolicy(0, "data2_admin")
	logErr("RemoveFilteredPolicy")
	err = e.LoadPolicy()
	logErr("LoadPolicy4")

	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}})

}

func TestAdapters(t *testing.T) {
	testSaveLoad(t)
	testAutoSave(t)
}
