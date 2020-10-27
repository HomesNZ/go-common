/*
	We use this package for migrations: https://github.com/DavidHuie/gomigrate

	Here's what you need to know:
	- Migrations go in the directory defined by the "migrationsDir" constant below.
	- They should be named like so: {{ id }}_{{ name }}_{{ "up" or "down" }}.sql
		For example: 1_add_users_table_up.sql / 1_add_users_table_down.sql
	- Migrations will be run on start up from main.go
	- To rollback use the `rollback` and `steps` flags

	AUTHORISE_KEY is a configurable variable, check .env files for more info.
*/

package migrator
