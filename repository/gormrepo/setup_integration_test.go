package gormrepo_test

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/storage"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var dbName string

func TestMain(m *testing.M) {
	err := config.Load()
	if err != nil {
		fmt.Printf("Config error: %s\n", err.Error())
		os.Exit(1)
	}

	err = initLogging()
	if err != nil {
		fmt.Printf("Logging error: %s\n", err.Error())
		os.Exit(1)
	}

	conn, err := prepareDB()
	if err != nil {
		fmt.Printf("Prepare db error: %s", err.Error())
		os.Exit(1)
	}

	retCode := m.Run()
	db, err := conn.DB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.Close()
	os.Exit(retCode)
}

func initLogging() error {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	log := logrus.StandardLogger()
	level, err := logrus.ParseLevel("DEBUG")
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)

	return err
}

func cleanDB(t *testing.T, db *gorm.DB) {
	defer func(t *testing.T) {
		sqlDB, err := db.DB()
		require.NoError(t, err)
		// Close
		sqlDB.Close()
	}(t)
	defer func(t *testing.T) {
		err := storage.TruncateNonRefTables(db)
		require.NoError(t, err)
	}(t)
}

func prepareDB() (dbConn *gorm.DB, err error) {
	dbName = "t_" + RandomString(10)
	err = storage.CreatePostgresDb(dbName)
	if err != nil {
		return
	}

	dbConn = storage.PostgresDbConn(&dbName)
	sqlDb, err := dbConn.DB()
	if err != nil {
		return
	}

	err = storage.MigratePostgresDb(sqlDb, nil, false, -1)
	if err != nil {
		return
	}

	return
}

func RandomString(n int) string {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes, err := RandomBytes(n)
	if err != nil {
		return ""
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}

func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
