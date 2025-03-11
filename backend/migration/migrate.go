package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"strings"

	"backend/util"

	"github.com/pressly/goose"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

var (
	flags  = flag.NewFlagSet("goose", flag.ExitOnError)
	dir    = flags.String("dir", "./sql", "directory with migration files")
	config *viper.Viper
)

func init() {
	config = viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("../env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}
}

func main() {
	dbString, err := util.NewDBStringFromConfig(config)
	if err != nil {
		log.Fatalf("fail to create db connection string: %v", err)
	}

	flags.Parse(os.Args[1:])

	driver := "postgres"
	command := "up"
	args := flags.Args()
	if len(args) > 0 {
		command = args[0]
	}

	if err := goose.SetDialect(driver); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(driver, dbString)
	if err != nil {
		log.Fatalf("-dbstring=%q: %v\n", dbString, err)
	}

	if err := goose.Run(command, db, *dir); err != nil {
		log.Fatalf("goose run: %v", err)
	}
}
