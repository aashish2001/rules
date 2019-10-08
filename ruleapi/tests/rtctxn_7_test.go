package tests

import (
	"context"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"

	"github.com/stretchr/testify/assert"
)

//add and delete in the same rtc
func Test_T7(t *testing.T) {

	rs, err := createRuleSession()
	assert.Nil(t, err)

	rule := ruleapi.NewRule("R7")
	err = rule.AddCondition("R7_c1", []string{"t1.none"}, trueCondition, nil)
	assert.Nil(t, err)
	rule.SetActionService(createActionServiceFromFunction(t, r7_action))
	rule.SetPriority(1)
	err = rs.AddRule(rule)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	txnCtx := txnCtx{t, 0}
	rs.RegisterRtcTransactionHandler(t7Handler, &txnCtx)
	err = rs.Start(nil)
	assert.Nil(t, err)

	i1, err := model.NewTupleWithKeyValues("t1", "t10")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), i1)
	assert.Nil(t, err)

	deleteRuleSession(t, rs)

}

func r7_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t1")].(model.MutableTuple)
	id, _ := t1.GetString("id")
	if id == "t10" {
		rs.Delete(ctx, t1)
	}
}

func t7Handler(ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {
	if done {
		return
	}

	txnCtx := handlerCtx.(*txnCtx)
	txnCtx.TxnCnt++
	t := txnCtx.Testing
	if txnCtx.TxnCnt == 1 {
		lA := len(rtxn.GetRtcAdded())
		if lA != 1 {
			t.Errorf("RtcAdded: Expected [%d], got [%d]\n", 1, lA)
			printTuples(t, "Added", rtxn.GetRtcAdded())
		}
		lM := len(rtxn.GetRtcModified())
		if lM != 0 {
			t.Errorf("RtcModified: Expected [%d], got [%d]\n", 0, lM)
			printModified(t, rtxn.GetRtcModified())
		}
		lD := len(rtxn.GetRtcDeleted())
		if lD != 1 {
			t.Errorf("RtcDeleted: Expected [%d], got [%d]\n", 1, lD)
			printTuples(t, "Deleted", rtxn.GetRtcDeleted())
		}
	}
}
