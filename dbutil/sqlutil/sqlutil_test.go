package sqlutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var checkData = "insert into tb_test (source,tiny_url,created_on) values (?,?,?),[test test test];update tb_test set source=? where tiny_url=?,[test value test value];select tiny_url,created_on,origin_url,source from tb_test where tiny_url=? and ( created_on=? or source=?) order by created_on asc,tiny_url desc limit 10,[test value or value 2];select count(*) from tb_test  where tiny_url=? ,[test value];delete from tb_test where tiny_url=?,[test value];"

func TestExample(t *testing.T) {
	//fmt.Println(exampleBase(false))
	a := exampleBase(false) == checkData
	assert.True(t, a)
}

func TestCOALESCE(t *testing.T) {
	assert.Equal(t, "COALESCE(mini_program_count,'0')", COALESCE("mini_program_count", "0"))
}

func TestExample2(t *testing.T) {
	result := exampleBase(true) != ""
	assert.True(t, result)
}

func TestSQLGen_Update(t *testing.T) {
	//没有设置条件的情况
	sqlGen := NewSQLGen("t").UpdateColumn("name", "test")
	_, _, e := sqlGen.Update()
	assert.Equal(t, e, ErrorCheckoutSQLCondition)
	//断言强制执行
	_, _, e = sqlGen.ForceExecOnNoCondition().Update()
	assert.Equal(t, nil, e)
}

func TestSQLGen_CustomConditionAndArgsAppend(t *testing.T) {
	updateStr, args, _ := NewSQLGen("t").UpdateColumn("name", "test").And("id", 10).OrderByDesc("name").Limit(10).Update()
	//测试args
	mockArgs := []interface{}{"test",10}
	for i, arg := range args {
		assert.Equal(t, mockArgs[i], arg)
	}
	assert.Equal(t, 2, len(args))
	assert.Equal(t, "update t set name=? where id=? order by name desc limit 10", updateStr)
	queryStr, _, _ := NewSQLGen("t").QueryColumns("name", "test").And("age", 10).CustomConditionAndArgsAppend("and name=?", "test").Query()
	assert.Equal(t, "select name,test from t where age=? and name=?", queryStr)
	deleteStr, _, _ := NewSQLGen("t").And("age", 10).And("id>", 10).OrderByAsc("money").Limit(10).Delete()
	assert.Equal(t, "delete from t where age=? and id>=? order by money asc limit 10", deleteStr)
	//insertStr, _, _ := NewSQLGen("t").InsertColumn("age", 10).InsertColumn("age2", 10).CustomConditionAndArgsAppend("and name=?", "test").Insert()
	//fmt.Println(insertStr)
	//assert.Equal(t, "update `t` set  name=?  where 1=1  and name=?", queryStr)
}
