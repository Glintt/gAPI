package migrations

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"

	"github.com/Glintt/gAPI/api/utils"
)

const oracleMigrationFolder = "./migrations/oracle"

const (
	MIGRATION_EXISTS        = `select * from gapi_migrations where rownum = 1`
	MIGRATION_ALREADY_RUN   = `select COUNT(*) from gapi_migrations where id = :id`
	CREATE_MIGRATIONS_TABLE = `create table gapi_migrations (
		id VARCHAR2(255),
		created_at DATE DEFAULT (sysdate)
	)`
	ADD_RUN_MIGRATION = `insert into gapi_migrations(id, created_at) values (:id, DEFAULT)`
)

const RUN_MIGRATION_ENV_VAR = "RUN_MIGRATIONS"

func MigrateOracle(connectionString string) {
	if os.Getenv(RUN_MIGRATION_ENV_VAR) != "true" {
		return
	}

	db, err := sql.Open("goracle", connectionString)

	files, err := ioutil.ReadDir(oracleMigrationFolder)
	if err != nil {
		log.Fatal(err)
	}

	utils.LogMessage("==== RUNNING MIGRATIONS ====", utils.InfoLogType)

	_, err = db.Exec(CREATE_MIGRATIONS_TABLE)
	if err != nil {
		_, err = db.Exec(MIGRATION_EXISTS)
		if err != nil {
			utils.LogMessage("error creating migration table: "+err.Error(), utils.InfoLogType)
			return
		}
	}

	for _, f := range files {
		migrationID := f.Name()
		rows, err := db.Query(MIGRATION_ALREADY_RUN, migrationID)

		count := checkCount(rows)
		if count > 0 {
			continue
		}

		dat, _ := ioutil.ReadFile(oracleMigrationFolder + "/" + f.Name())

		utils.LogMessage("    -----> "+migrationID, utils.InfoLogType)
		_, err = db.Exec(string(dat))
		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = db.Exec(ADD_RUN_MIGRATION, migrationID)
	}

	utils.LogMessage("==== MIGRATIONS FINISHED RUNNING ====", utils.InfoLogType)

	defer db.Close()
}

func checkCount(rows *sql.Rows) (count int) {
	c := 0
	for rows.Next() {
		rows.Scan(&c)
	}
	return c
}
