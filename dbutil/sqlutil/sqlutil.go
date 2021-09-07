package sqlutil

import (
	"errors"
	"strings"

	"github.com/sssvip/goutil/logutil"
	"github.com/sssvip/goutil/strutil"
)

var ErrorCheckoutSQLCondition = errors.New("checkout conditions,let sql safe,you can sqlGen.ForceExecOnNoCondition() force exec current sql")

type SQLGen struct {
	printError   bool
	tableName    string
	queryColumns []string
	//辅助有序(输出的sql更加稳定)
	insertColumnKeys       []string
	insertColumnMap        map[string]interface{}
	updateColumnKeys       []string
	updateColumnMap        map[string]interface{}
	andConditionKeys       []string
	andConditionMap        map[string]interface{}
	orConditionKeys        []string
	orConditionMap         map[string]interface{}
	customCondition        string
	customConditionArgs    []interface{}
	orderByConditions      []string
	limit                  int // 常规简单limit
	limitEnd               int // limit 跳页 "limit 10,20"中的20
	forceExecOnNoCondition bool
}

func NewSQLGen(tableName string) *SQLGen {
	sqlGen := SQLGen{
		printError:      true,
		tableName:       tableName,
		insertColumnMap: make(map[string]interface{}),
		updateColumnMap: make(map[string]interface{}),
		andConditionMap: make(map[string]interface{}),
		orConditionMap:  make(map[string]interface{}),
		limit:           -1,
	}
	return &sqlGen
}

func (sqlGen *SQLGen) PrintError() bool {
	return sqlGen.printError
}

func (sqlGen *SQLGen) ClosePrintError() {
	sqlGen.printError = false
}

/*
	条件部分
*/
func (sqlGen *SQLGen) And(columnName string, condition interface{}) *SQLGen {
	sqlGen.andConditionKeys = append(sqlGen.andConditionKeys, columnName)
	sqlGen.andConditionMap[columnName] = condition
	return sqlGen
}

func (sqlGen *SQLGen) Or(columnName string, condition interface{}) *SQLGen {
	sqlGen.orConditionKeys = append(sqlGen.orConditionKeys, columnName)
	sqlGen.orConditionMap[columnName] = condition
	return sqlGen
}

func safeFormatValue(value interface{}) string {
	switch value.(type) {
	case int, int64:
		return strutil.Format("%d", value)
	case string:
		return strutil.Format(`'%s'`, value)
	default:
		return strutil.Format(`'%v'`, value)
	}
}

/*func safeFormatColumn(column string) string {
	return column
	//return strutil.Format("`%s`", column)
}*/

/*func safeFormatKV(columnName string, value interface{}) string {
	return strutil.Format(" %s=%s", safeFormatColumn(columnName), safeFormatValue(value))
}*/

func safeFormatKWithPlaceHolder(columnName string) string {
	return strutil.Format(" %s=?", columnName)
}

// condition一旦设置,追加到除了order by,limit关键字外的所有条件最后（需要用到条件时才生效）
// Deprecated: recommend use CustomConditionAndArgsAppend to use k,v sql .
func (sqlGen *SQLGen) CustomConditionAppend(condition string) *SQLGen {
	if condition != "" {
		sqlGen.customCondition += strutil.Format(" %s", condition)
	}
	return sqlGen
}

// condition一旦设置,追加到除了order by,limit关键字外的所有条件最后（需要用到条件时才生效）
func (sqlGen *SQLGen) CustomConditionAndArgsAppend(condition string, args ...interface{}) *SQLGen {
	if condition != "" {
		sqlGen.customCondition += strutil.Format(" %s", condition)
		sqlGen.customConditionArgs = append(sqlGen.customConditionArgs, args...)
	}
	return sqlGen
}

func (sqlGen *SQLGen) genConditions() (sqlStr string, args []interface{}) {
	sqlCondition := " where 1=1"
	//处理and逻辑
	var andConditions []string
	for _, key := range sqlGen.andConditionKeys {
		andConditions = append(andConditions, safeFormatKWithPlaceHolder(key))
		args = append(args, sqlGen.andConditionMap[key])
	}
	var andCondition string
	if len(andConditions) > 0 {
		andCondition = strutil.Format(" and%s", strings.Join(andConditions, " and"))
	}
	//处理or逻辑
	var orConditions []string
	for _, key := range sqlGen.orConditionKeys {
		orConditions = append(orConditions, safeFormatKWithPlaceHolder(key))
		args = append(args, sqlGen.orConditionMap[key])
	}
	var orCondition string
	if len(orConditions) > 0 {
		orCondition = strutil.Format(" and (%s)", strings.Join(orConditions, " or"))
	}
	sql := strutil.Format("%s%s%s", sqlCondition, andCondition, orCondition)
	if len(andCondition) < 1 && len(orConditions) < 1 && sqlGen.customCondition == "" {
		return "", args
	}
	if sqlGen.customCondition != "" {
		args = append(args, sqlGen.customConditionArgs...)
	}
	return strings.Replace(sql, "where 1=1 and", "where", 1), args
}

func (sqlGen *SQLGen) Count() (sqlStr string, args []interface{}, err error) {
	conditions, tArgs := sqlGen.genConditions()
	countCondition := "count(*)"
	/*if len(sqlGen.queryColumns) > 1 {
		return "", nil, errors.New(strutil.Format("count method just max allowe [1] arg, but accept [%d] args,-> %s", len(sqlGen.queryColumns), sqlGen.queryColumns))
	}*/
	if len(sqlGen.queryColumns) == 1 && (strings.HasPrefix(sqlGen.queryColumns[0], "count") || strings.HasPrefix(sqlGen.queryColumns[0], "COUNT")) {
		countCondition = sqlGen.queryColumns[0]
	}
	return strutil.Format("select %s from %s %s %s", countCondition, sqlGen.tableName, conditions, sqlGen.customCondition), tArgs, nil
}

func (sqlGen *SQLGen) ForceExecOnNoCondition() *SQLGen {
	sqlGen.forceExecOnNoCondition = true
	return sqlGen
}

/*
	修改部分
*/

func (sqlGen *SQLGen) Update() (sqlStr string, args []interface{}, err error) {
	var updateColumns []string
	for _, key := range sqlGen.updateColumnKeys {
		updateColumns = append(updateColumns, strutil.Format("%s", safeFormatKWithPlaceHolder(key)))
		args = append(args, sqlGen.updateColumnMap[key])
	}
	conditions, tArgs := sqlGen.genConditions()
	args = append(args, tArgs...)
	sqlStr = strutil.Format("update %s set%s%s%s%s%s", sqlGen.tableName, strings.Join(updateColumns, ","), conditions, sqlGen.getOrderConditions(), sqlGen.limitCondition(), sqlGen.customCondition)
	if !sqlGen.forceExecOnNoCondition && conditions == "" {
		logutil.Error.Println("sqlStr:" + sqlStr)
		return "", nil, ErrorCheckoutSQLCondition
	}
	return
}

func (sqlGen *SQLGen) UpdateColumn(columnName string, value interface{}) *SQLGen {
	sqlGen.updateColumnKeys = append(sqlGen.updateColumnKeys, columnName)
	sqlGen.updateColumnMap[columnName] = value
	return sqlGen
}

func (sqlGen *SQLGen) Update2Insert() *SQLGen {
	for _, key := range sqlGen.updateColumnKeys {
		sqlGen.InsertColumn(key, sqlGen.updateColumnMap[key])
	}
	return sqlGen
}

func (sqlGen *SQLGen) Insert2Update() *SQLGen {
	for _, key := range sqlGen.insertColumnKeys {
		sqlGen.UpdateColumn(key, sqlGen.insertColumnMap[key])
	}
	return sqlGen
}

/*
	新增部分
*/

func (sqlGen *SQLGen) InsertColumn(columnName string, value interface{}) *SQLGen {
	sqlGen.insertColumnKeys = append(sqlGen.insertColumnKeys, columnName)
	sqlGen.insertColumnMap[columnName] = value
	return sqlGen
}

func (sqlGen *SQLGen) Insert() (sqlStr string, args []interface{}, err error) {
	var columns []string
	var placeHolders []string
	for _, key := range sqlGen.insertColumnKeys {
		columns = append(columns, key)
		args = append(args, sqlGen.insertColumnMap[key])
		placeHolders = append(placeHolders, "?")
	}
	sqlStr = strutil.Format("insert into %s (%s) values (%s)", sqlGen.tableName, strings.Join(columns, ","), strings.Join(placeHolders, ","))
	return
}

/*
	删除部分
*/

func (sqlGen *SQLGen) Delete() (sqlStr string, args []interface{}, err error) {
	conditions, tArgs := sqlGen.genConditions()
	sqlStr = strutil.Format("delete from %s%s%s%s%s", sqlGen.tableName, conditions, sqlGen.getOrderConditions(), sqlGen.limitCondition(), sqlGen.customCondition)
	if !sqlGen.forceExecOnNoCondition && conditions == "" {
		logutil.Error.Println("sqlStr:" + sqlStr)
		return "", nil, ErrorCheckoutSQLCondition
	}
	return sqlStr, tArgs, nil
}

/*
	查询部分
*/

func (sqlGen *SQLGen) getOrderConditions() string {
	if len(sqlGen.orderByConditions) < 1 {
		return ""
	}
	return strutil.Format(" order by %s", strings.Join(sqlGen.orderByConditions, ","))
}

func (sqlGen *SQLGen) OrderByDesc(column string) *SQLGen {
	sqlGen.orderByConditions = append(sqlGen.orderByConditions, strutil.Format("%s desc", column))
	return sqlGen
}

func (sqlGen *SQLGen) OrderByAsc(column string) *SQLGen {
	sqlGen.orderByConditions = append(sqlGen.orderByConditions, strutil.Format("%s asc", column))
	return sqlGen
}

func (sqlGen *SQLGen) limitCondition() string {
	if sqlGen.limit > 0 && sqlGen.limitEnd > 0 {
		return strutil.Format(" limit %d,%d", sqlGen.limit, sqlGen.limitEnd)
	}
	if sqlGen.limit > 0 {
		return strutil.Format(" limit %d", sqlGen.limit)
	}
	return ""
}

func (sqlGen *SQLGen) Limit(limit int) *SQLGen {
	sqlGen.limit = limit
	return sqlGen
}

func (sqlGen *SQLGen) LimitPage(limit, limitEnd int) *SQLGen {
	sqlGen.limit = limit
	sqlGen.limitEnd = limitEnd
	return sqlGen
}

func (sqlGen *SQLGen) QueryColumnsCount() int {
	return len(sqlGen.queryColumns)
}

func COALESCE(columnName string, defaultValue interface{}) string {
	return strutil.Format("COALESCE(%s,%s)", columnName, safeFormatValue(defaultValue))
}

func (sqlGen *SQLGen) QueryColumns(columns ...string) *SQLGen {
	for _, c := range columns {
		sqlGen.queryColumns = append(sqlGen.queryColumns, c)
	}
	return sqlGen
}

func (sqlGen *SQLGen) getQueryColumns() string {
	/*var tColumns []string
	for _, c := range sqlGen.queryColumns {
		if strings.HasPrefix(c, "distinct ") || strings.HasPrefix(c, "DISTINCT ") {
			tColumns = append(tColumns, strutil.Format("DISTINCT COALESCE(%s, '')", safeFormatColumn(c[len("distinct "):])))
			continue
		}
		tColumns = append(tColumns, strutil.Format("COALESCE(%s, '')", safeFormatColumn(c)))
	}*/
	return strings.Join(sqlGen.queryColumns, ",")
}

func (sqlGen *SQLGen) Query() (sqlStr string, args []interface{}, err error) {
	var errStr string
	if len(sqlGen.queryColumns) < 1 {
		errStr = "must assign query columns by QueryColumns()"
		logutil.Error.Println(errStr)
		return "", args, errors.New(errStr)
	}
	conditions, tArgs := sqlGen.genConditions()
	return strutil.Format("select %s from %s%s%s%s%s", sqlGen.getQueryColumns(), sqlGen.tableName, conditions, sqlGen.customCondition, sqlGen.getOrderConditions(), sqlGen.limitCondition()), tArgs, nil
}
