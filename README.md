any-migrate
===========

A generic data migrator.

**Status**: Planning

Rationale
---------

For proper continuous integration (CI) or continuous delivery (CD), many teams
use a database schema migration tool. Such tool makes it possible to versionize
and codify database schema changes to make deployments repeatable, revertable,
testable and automated. The general approach to do this is by

 * versioning database (table) schemas by defining ordered schema changes/deltas; and
 * as-needed applying these changes in a specific order. Bringing the database
   to the latest schema version is generally an
   [idempotent](https://en.wikipedia.org/wiki/Idempotence) operation. To
   support idenpotency, some tools keep track of every change applied in a
   database of its own.

[Most](http://migrate4j.sourceforge.net/)
[schema](https://liquibase.jira.com/wiki/display/CONTRIB/Available+Extensions)
[migration](https://flywaydb.org) [tools](http://www.mybatis.org/migrations/)
[are](https://bitbucket.org/liamstask/goose) [tied](http://sqitch.org) to the
low-level concept of a database and very specific interfaces, such as SQL. This
often limits them to very basic data migrations such as creating/dropping a
table, adding/deleting fields or updating a value of some field. At its best
they support custom SQL.

Unfortunately, simple SQL executions aren't always enough when migrating
application state. A deploy might require a more complex data migration than an
SQL execution; maybe you want to run a script, or support migrating a NoSQL
database of some kind. The script could clean up your database, execute
[`pt-online-schema-change`](https://www.percona.com/doc/percona-toolkit/2.1/pt-online-schema-change.html)
(for MySQL), do complex migrations of users, populate complex data, or other
one-off tasks you'd like to execute. `any-migrate` is a tool aiming to support
custom data migrations.

`any-migrate` will have the following features:

 * Keeps track of which migrations has been executed on a system.
 * Fail early. Better safe than sorry.
 * Asks, doesn't assume. Explicitness is important. No magic.
 * Verifies every change as much as possible before actual execution.
 * Support dry-run mode.
 * Pluggable storage _repositories_ of which version have been applied.
   Examples are:
   * On disk (mostly used for testing - make sure to backup!).
   * RDBM (MySQL/PostgresQL).
 * Pluggable system for migration _migrators_. Examples are:
   * Executing a script.
   * Executing an SQL query on a server.
   * Executing a shell command.
 * Easy deployment. Written in [Go](https://golang.org). Application, and
   plugins, compiles to statically linked binaries. No JVM, Python, node.js or
   Ruby environment and third party libraries required.
 * A basic test framework of migrations and reverts.

Details
-------

_Migrations_ are strictly ordered and applied in ascending order. `any-migrate`
will fail if two migrations have gaps in the migration numbering, or two
migrations have the same number. A migration is also not allowed to run unless
all previous migrations are in `MIGRATION_SUCCEEDED` state.

Each migration is a [finite-state
machine](https://en.wikipedia.org/wiki/Finite-state_machine) with the following
states:

| State                       | Description                                                |
|-----------------------------|------------------------------------------------------------|
| `NEW`                       | Unapplied migration.                                       |
| `TESTING`                   | The migration is being tested if it looks sane.            |
| `TESTING_FAILED`            | `TESTING` failed.                                          |
| `MIGRATING`                 | The migration is being applied.                            |
| `MIGRATION_FAILED`          | The migration failed.                                      |
| `VERIFYING`                 | The migration is verifying that the migration was correct. |
| `VERIFICATION_FAILED`       | `VERIFYING` failed.                                        |
| `MIGRATION_SUCCEEDED`       | The migration succeeded.                                   |

There are additional states for migrations that support reverting migrations:

| State                        | Description                              |
|------------------------------|------------------------------------------|
| `REVERT_TESTING`             | The initial stage when running a revert. |
| `REVERT_TESTING_FAILED`      | `REVERT_TESTING` failed.                 |
| `REVERT_MIGRATING`           | Revert is being executed.                |
| `REVERT_MIGRATION_FAILED`    | `REVERT_MIGRATING` failed.               |
| `REVERT_VERIFYING`           | Verifying that a revert succeeded.       |
| `REVERT_VERIFICATION_FAILED` | The revert failed.                       |
| `REVERT_SUCCEEDED`           | The revert succeeded.                    |

Reverting a migration is useful for quick rollback of a failed. Also, each
state is also accompanied by a free-text `payload` containing additional
information.

The following state transitions exist:

| Name | Source states | Destination state |
|------|---------------|-------------------|
| migrate | `NEW` | `TESTING` |
| fail-testing | `TESTING` | `TESTING_FAILED` |
| execute-migration | `TESTING` | `MIGRATING` |
| fail-migration | `MIGRATING` | `MIGRATION_FAILED` |
| verify-migration | `MIGRATING` | `VERIFYING`
| fail-verification | `VERIFYING` | `VERIFICATION_FAILED` |
| succeed-migration | `VERIFYING` | `MIGRATION_SUCCEEDED` |
| retry-migration | `TESTING_FAILED`, `MIGRATION_FAILED`, `VERIFICATION_FAILED` | `TESTING` |
| revert-migration | `MIGRATION_SUCCEEDED` | `REVERT_TESTING` |
| fail-revert-testing | `REVERT_TESTING` | `REVERT_TESTING_FAILED` |
| execute-revert | `REVERT_TESTING` | `REVERT_MIGRATING` |
| fail-revert-execution | `REVERT_MIGRATING` | `REVERT_MIGRATION_FAILED` |
| verify-revert | `REVERT_MIGRATING` | `REVERT_VERIFYING` |
| fail-revert-verification | `REVERT_VERIFYING` | `REVERT_VERIFICATION_FAILED` |
| succeed-migraton | `REVERT_VERIFYING` | `REVERT_SUCCEEDED` |
| retry-revert | `REVERT_TESTING_FAILED`, `REVERT_MIGRATION_FAILED`, `REVERT_VERIFICATION_FAILED` | `REVERT_TESTING` |

Getting started
---------------

First, [install Go](https://golang.org/dl/) and then set the `GOPATH`
environment variable to a directory where you'd like to install the binaries.
You can then install the binaries needed for this example by issuing

```bash
$ go install github.com/any-migrate/any-migrate/...
```

`any-migrate` is the command line utility that you'll use to interact with
any-migrate. A _repository_ is a plugin which is used for keeping track of
which migrations have been applied. The _file_ repository plugin stores this
information in a plain file. A _migrator_ is a plugin which applies changes.
The _mysql_ plugin applies changes to a MySQL database.

Then create a configuration (YAML file):

```yaml
# my.config
repository:
  type: file
  path: ./my.repository
migrations:
  path: ./my.changes
sources:
  mysql:
    migrator: mysql-migrator
    url: /var/run/mysqld/mysqld.sock
    username: root
    database: test_db
```

and execute `mkdir my.changes` to create the directory where we will store our
migrations.

Now, let's create and run our first migrations! Create the following files:

```yaml
# my.changes/01create_database.migration
# `source` is optional if only a single source is defined in `my.config`.
source: mysql
action:
  type: create database
  # If `database` is specified for source, then it's optional here.
  database: test_db
```
,
```yaml
# my.changes/02create_table.migration
action:
  type: raw
  migration: |
    CREATE TABLE `users` (
      `id` varchar(255) NOT NULL,
      `username` varchar(255) DEFAULT NULL,
      PRIMARY KEY (`id`),
      UNIQUE KEY `username` (`username`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
  reversion: |
    DROP TABLE `users`;
```
and
```yaml
# my.changes/03add_column.migration
action:
  type: add column
  column:
    name: age
    specification: INT DEFAULT NULL
```
. Each version must be an incremented integer numbered from 1. Gaps in versions
are not allowed.

The migration setup can be tested by executing

```bash
$ $GOPATH/bin/any-migrate -f my.config test
```

and executed using

```bash
$ $GOPATH/bin/any-migrate -f my.config migrate
```

Roadmap
-------
`any-migrate` is currently in planning stage. Please follow the
[issues](https://www.github.com/JensRantil/any-migrate/issues) to keep track of
progress.

Find a missing feature? Please [submit
one](https://www.github.com/JensRantil/any-migrate/issues/new) to have us
discuss it!

Best practices
--------------
To create a great tool, there are some practise we try to adhere to:

 * Fail early. This includes running validation before making actual changes if
   possible.
 * Maximize trust. Trust in a migrator is uterly important. Upgrades should not
   break. Related to above;
 * Ask before assuming something. Explicitness is important. No magic.
 * _Always_ test migrations before applying to production. Preferably on an
   isolated staging environment. This will make you confident that you haven't
   made any errors.
