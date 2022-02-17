package repository

import (
	// "log"
	// "os"
	// "github.com/joho/godotenv"
	"log"

	httperors "github.com/myrachanto/custom-http-error"
	"github.com/myrachanto/microservice/user/src/model"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//IndexRepo
var (
	IndexRepo indexRepo = indexRepo{}
	Operator            = map[string]string{"all": "all", "equal_to": "=", "not_equal_to": "<>", "less_than": "<",
		"greater_than": ">", "less_than_or_equal_to": "<=", "greater_than_ro_equal_to": ">=",
		"like": "like", "between": "between", "in": "in", "not_in": "not_in"}
)

//Layout ...
const (
	Layout   = "2006-01-02"
	layoutUS = "January 2, 2006"
)

type Db struct {
	DbType     string `mapstructure:"DbType"`
	DbName     string `mapstructure:"DbName"`
	DbUsername string `mapstructure:"DbUsername"`
	DbPassword string `mapstructure:"DbPassword"`
	DbPort     string `mapstructure:"DbPort"`
}

func LoaddbConfig() (db Db, err error) {
	viper.AddConfigPath("../../")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&db)
	return
}

///curtesy to gorm
type indexRepo struct{}

func (indexRepo indexRepo) InitDB() httperors.HttpErr {

	sdb, ers := LoaddbConfig()
	if ers != nil {
		return httperors.NewNotFoundError("Something went wrong with viper loading db config --db!")
	}
	GormDB, err := gorm.Open(mysql.New(mysql.Config{
		// DSN: sdb.DbUsername + ":" + sdb.DbPassword + "@tcp(user_database:3306)/" + sdb.DbName + "?charset=utf8&parseTime=True&loc=Local", // data source name
		DSN:                       sdb.DbUsername + ":" + sdb.DbPassword + "@tcp(127.0.0.1:3306)/" + sdb.DbName + "?charset=utf8&parseTime=True&loc=Local", // data source name
		DefaultStringSize:         256,                                                                                                                     // default size for string fields
		DisableDatetimePrecision:  true,                                                                                                                    // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                                                                    // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                                                                    // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                                                                   // auto configure based on currently MySQL version
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return httperors.NewNotFoundError("Something went wrong with viper loading db config --db!")
	}

	GormDB.AutoMigrate(&model.User{})
	GormDB.AutoMigrate(&model.Auth{})
	return nil
}
func (indexRepo indexRepo) Getconnected() (GormDB *gorm.DB, err httperors.HttpErr) {
	log.Println("Db Prep --------")
	sdb, ers := LoaddbConfig()
	if ers != nil {
		log.Println(ers)
		return nil, httperors.NewNotFoundError("Something went wrong with viper --db!")
	}
	log.Println("Db connection --------")
	GormDB, err1 := gorm.Open(mysql.New(mysql.Config{
		// DSN: sdb.DbUsername + ":" + sdb.DbPassword + "@tcp(user_database:3306)/" + sdb.DbName + "?charset=utf8&parseTime=True&loc=Local", // data source name
		DSN:                       sdb.DbUsername + ":" + sdb.DbPassword + "@tcp(127.0.0.1:3306)/" + sdb.DbName + "?charset=utf8&parseTime=True&loc=Local", // data source name
		DefaultStringSize:         256,                                                                                                                     // default size for string fields
		DisableDatetimePrecision:  true,                                                                                                                    // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                                                                    // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                                                                    // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                                                                   // auto configure based on currently MySQL version
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err1 != nil {
		return nil, httperors.NewNotFoundError("Something went wrong with viper --db!")
	}
	return GormDB, nil
}
func (indexRepo indexRepo) DbClose(GormDB *gorm.DB) {
	// defer GormDB.Close()
}

//Paginate the data from backend
func Paginate(page, pagesize int) func(GormDB *gorm.DB) *gorm.DB {
	return func(GormDB *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pagesize > 100:
			pagesize = 100
		case pagesize <= 0:
			pagesize = 10
		}

		offset := (page - 1) * pagesize
		return GormDB.Offset(offset).Limit(pagesize)
	}
}
