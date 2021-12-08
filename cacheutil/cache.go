package cacheutil

import (
	"os"

	"github.com/sssvip/goutil/dbutil"
	"github.com/sssvip/goutil/dbutil/sqlutil"
	"github.com/sssvip/goutil/logutil"
	"github.com/sssvip/goutil/strutil"
)

const (
	defaultDBDir               = "./sqliteDB"
	defaultDBFileName          = "acq_key_value_cache"
	defaultDBTableName         = "acq_key_value_cache"
	defaultCrateTableNotExists = "create table if not exists "
	defaultCrateTableSQL       = `
	(
		id         integer primary key autoincrement,
		key        TEXT unique,
		value      TEXT,
		gmt_create timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime')),
		gmt_modify timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime'))
	);`
)

type CacheConfig struct {
	DBDir       string
	DBFileName  string
	DBTableName string
}

type CacheDB struct {
	Config CacheConfig
	dbFile string
	db     *dbutil.DBWrapper
}

func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		DBDir:       defaultDBDir,
		DBFileName:  defaultDBFileName,
		DBTableName: defaultDBTableName,
	}
}

func NewCache(cfg CacheConfig) CacheUtil {
	cache := &CacheDB{Config: cfg}
	cache.Init()
	return cache
}

func (c *CacheDB) Init() {
	c.dbFile = c.Config.DBDir + "/" + c.Config.DBFileName
	if e := c.dbFileCreate(); e != nil {
		return
	}
	c.db = dbutil.NewSqliteDB(c.dbFile, "", "")
	if e := c.dbTableCreate(); e != nil {
		return
	}
}

func (d *CacheDB) dbFileCreate() error {
	if isExist(d.dbFile) {
		return nil
	}
	// 创建 DB 文件目录
	err := os.MkdirAll(d.Config.DBDir, os.ModePerm)
	if err != nil {
		logutil.Error.Println("creat db dir failed, err: ", err)
		return err
	}
	f, e := os.OpenFile(d.dbFile, os.O_CREATE|os.O_RDWR, os.ModePerm)
	// 判断是否出错
	if e != nil {
		logutil.Error.Println("creat db file failed, err: ", err)
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	return nil
}

func isExist(fileAddr string) bool {
	// 读取文件信息，判断文件是否存在
	_, err := os.Stat(fileAddr)
	if err != nil {
		if os.IsExist(err) { // 根据错误类型进行判断
			return true
		}
		return false
	}
	return true
}

func (c *CacheDB) dbTableCreate() error {
	_, err := c.db.Exec(strutil.Format("%s %s %s", defaultCrateTableNotExists, c.Config.DBTableName, defaultCrateTableSQL))
	if err != nil {
		logutil.Error.Println("db sql exec failed, err:", err)
		return err
	}
	return nil
}

func (c *CacheDB) Set(key, value string) (err error) {
	if c.Has(key) {
		_, err = c.db.UpdateTableBySQLGen(sqlutil.NewSQLGen(c.Config.DBTableName).
			UpdateColumn("value", value).
			And("key", key))
	} else {
		_, err = c.db.InsertTableBySQLGen(sqlutil.NewSQLGen(c.Config.DBTableName).
			InsertColumn("value", value).
			InsertColumn("key", key))
	}
	return err
}

func (c *CacheDB) Has(key string) bool {
	minCount := 1
	cnt, err := c.db.CountBySQLStr(strutil.Format("select count(*) from %s where key='%s'", c.Config.DBTableName, key))
	if err != nil {
		logutil.Error.Println(err)
		return false
	}
	if cnt < minCount {
		return false
	}
	return true
}

func (c *CacheDB) Get(key string, defaultValue ...string) string {
	row, err := c.db.GetRowBySQLGen(sqlutil.NewSQLGen(c.Config.DBTableName).
		QueryColumns("value").
		And("key", key))
	if err != nil {
		logutil.Error.Println(err)
	}
	if len(row) < 1 {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return row[0]
}

func (c *CacheDB) GetDetail(key string) *KeyData {
	var keyData KeyData
	row, err := c.db.GetRowBySQLStr(strutil.Format("select key,value,gmt_create,gmt_modify from %s "+
		"where key='%s'", c.Config.DBTableName, key))
	if err != nil {
		logutil.Error.Println(err)
		return nil
	}
	if len(row) < 2 {
		return nil
	}
	keyData.Key = row[0]
	keyData.Value = row[1]
	keyData.GmtCreate = row[2]
	keyData.GmtModify = row[3]
	return &keyData
}

func (c *CacheDB) GetOriginDB() *dbutil.DBWrapper {
	return c.db
}

func (c *CacheDB) GetCacheConfig() CacheConfig {
	return c.Config
}

type KeyData struct {
	Key       string
	Value     string
	GmtCreate string
	GmtModify string
}

type CacheUtil interface {
	Init()
	Set(key, value string) error
	Has(key string) bool
	Get(key string, defaultValue ...string) string
	GetDetail(key string) *KeyData
	GetOriginDB() *dbutil.DBWrapper
	GetCacheConfig() CacheConfig
}
