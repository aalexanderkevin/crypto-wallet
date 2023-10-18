package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/repository/gormrepo"

	"github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	gorm_logrus "github.com/onrik/gorm-logrus"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetPostgresDb() *gorm.DB {
	dbName := config.Instance().Postgres.Database
	return PostgresDbConn(&dbName)
}

func PostgresDbConn(dbName *string) *gorm.DB {
	dbURL := getPostgresUrl(dbName)

	logrusLogger := gorm_logrus.New()
	logrusLogger.LogMode(logger.Silent)
	logrusLogger.Debug = false
	if config.Instance().Postgres.Debug {
		logrusLogger.Debug = true
		logrusLogger.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("error: %v for %v", err.Error(), dbURL))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("error: %v for %v", err.Error(), dbURL))
	}

	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.Instance().Postgres.MaxConnLifeTime))
	sqlDB.SetMaxOpenConns(config.Instance().Postgres.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(config.Instance().Postgres.MaxIdleConnections)

	return db
}

func CreatePostgresDb(dbName string) error {
	dbConn := PostgresDbConn(nil)

	return dbConn.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName)).Error
}

func MigratePostgresDb(db *sql.DB, migrationFolder *string, rollback bool, versionToForce int) error {
	dbConfig := config.Instance().Postgres

	var validMigrationFolder = dbConfig.Migration.Path
	if migrationFolder != nil && *migrationFolder != "" {
		validMigrationFolder = *migrationFolder
	}

	if validMigrationFolder == "" {
		return fmt.Errorf("empty migration folder")
	}
	logrus.Infof("Migration folder: %s", validMigrationFolder)

	driver, err := migrate_postgres.WithInstance(db, &migrate_postgres.Config{})
	if err != nil {
		logrus.WithError(err).Warning("Error when instantiating driver")
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+validMigrationFolder,
		dbConfig.Client,
		driver)
	if err != nil {
		logrus.WithError(err).Warning("Error when instantiating migrate")
		return err
	}
	if rollback {
		logrus.Info("About to Rolling back 1 step")
		err = m.Steps(-1)
	} else if versionToForce != -1 {
		logrus.Info(fmt.Sprintf("About to force version %d", versionToForce))
		err = m.Force(versionToForce)
	} else {
		logrus.Info("About to run migration")
		err = m.Up()
	}
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}

func CloseDB(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	sqlDB, _ := db.DB()
	err := sqlDB.Close()
	if err != nil {
		return err
	}
	return nil
}

func TruncateNonRefTables(db *gorm.DB) error {
	models := []interface{}{
		gormrepo.User{},
	}
	for _, v := range models {
		err := db.Statement.Parse(v)
		if err != nil {
			return err
		}

		tableName := db.Statement.Schema.Table
		err = db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func getPostgresUrl(dbName *string) string {
	dbConfig := config.Instance().Postgres

	dbNameTmp := "postgres"
	if dbName != nil {
		dbNameTmp = *dbName
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=disable", dbConfig.Host,
		dbConfig.Username, dbConfig.Password, dbNameTmp, dbConfig.Port)
}
