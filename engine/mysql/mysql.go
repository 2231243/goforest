package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	User   string
	Pass   string
	Host   string
	DBName string
	Port   int
	Drive  string
	Debug  bool
}
type Option struct {
	Config
	*gorm.DB
	sqlDB    *sql.DB
	settings []Setting
}

const (
	defaultDriver          = "mysql"
	defaultPort            = 3306
	defaultMaxIdleConn     = 10
	defaultMaxOpenConn     = 100
	defaultConnMaxLifetime = time.Hour
	defaultCharset         = "utf8mb4"
	defaultLoc             = "Asia/Shanghai"
	cDSNFormat             = "%s%s=%s&"
)

func (m *Option) set(s ...Setting) {
	m.settings = append(m.settings, s...)
}

func (m *Option) realDSN() string {
	format := "%s:%s@tcp(%s:%d)/%s?%s"
	return strings.TrimRight(fmt.Sprintf(format, m.User, m.Pass, m.Host, m.Port, m.DBName, concatDSN(m.settings)), "?")
}

func DefaultOption() gorm.Option {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,          // 禁用彩色打印
		},
	)
	return &gorm.Config{
		Logger:      newLogger.LogMode(logger.Info),
		PrepareStmt: true, //预编译
	}
}

func (o *Option) connect() error {
	var err error

	var driver gorm.Dialector
	switch o.Drive {
	case defaultDriver:
		driver = mysql.Open(o.realDSN())
	case "SQLite":
		driver = sqlite.Open(o.DBName)
	default:
		return errors.New("driver not implemented")
	}
	if o.DB, err = gorm.Open(driver, DefaultOption()); err != nil {
		return err
	}
	if o.sqlDB, err = o.DB.DB(); err != nil {
		return err
	}

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	o.sqlDB.SetMaxIdleConns(defaultMaxIdleConn)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	o.sqlDB.SetMaxOpenConns(defaultMaxOpenConn)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	o.sqlDB.SetConnMaxLifetime(defaultConnMaxLifetime)
	return nil
}

func (o *Option) Close() {
	_ = o.sqlDB.Close()
}

func DefaultDB(c Config, opt ...Setting) (*Option, error) {
	opt = append(opt, SetCharset(defaultCharset))
	opt = append(opt, SetLocal(url.QueryEscape(defaultLoc)))

	return NewDB(c, opt...)

}

func NewDB(c Config, opt ...Setting) (*Option, error) {
	o := &Option{
		Config: c,
	}
	o.Config.Port = defaultPort
	o.Config.Drive = defaultDriver
	o.set(opt...)
	if err := o.connect(); err != nil {
		return nil, err
	}
	if !c.Debug {
		o.DB.Logger = logger.Default.LogMode(logger.Silent)
	}
	return o, nil
}
