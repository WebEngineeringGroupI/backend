# Contributing to the repository

This file contains some guidelines to contribute to the project.

## Architecture

The architecture follows the principles
of [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
and [Domain Driven Design](https://domaindrivendesign.org/ddd-domain-driven-design/).

## Testing

All the project is developed with BDD using [Ginkgo](https://onsi.github.io/ginkgo/)
and [Gomega](https://onsi.github.io/gomega/). You can install the `ginkgo` executable by executing:

```
go install github.com/onsi/ginkgo/ginkgo@v1
```

If you just want to execute some tests in a package, move to this package and execute `ginkgo` directly, e.g:

```
cd pkg/domain/url
ginkgo -r
```

To execute all the unit tests you can use:

```
make test-unit
```

To execute all the integration tests, you can use:

```
make test-integration
```

These tests will first, launch a local postgres DB in Docker and run the migrations before executing the integration
tests for the DB. The database will be killed after the tests are run.

## Database management

The database is created using migration scripts. The main purpose is to ensure that database changes, being made and
committed to source control, can be properly applied without breaking up database integrity, making configuration
changes manually, or to prevent data loss.

For example, renaming a column will be treated by source control system as dropping a column with an old name, and
creating a new column with a different name. This will result in data loss, when that specific change is applied against
a database. With a migration script assigned to this change, a column will be renamed in a way the data won’t be lost.

Whereas a build script creates a database, a migration script, or ‘change’ script, alters a database. It is called a
migration script because it changes all or part of a database from one version to another. It ‘migrates’ it between
versions. This alteration can be as simple as adding or removing a column to a table, or a complex refactoring task such
as splitting tables or changing column properties in a way that could affect the data it stores.

For every likely migration path between database versions, we need to store in version control the migration scripts
that describe precisely those steps required to perform the change and, if necessary, moving data around and
transforming it in the process.

In this repository, all migration scripts will be created and applied
by [golang-migrate](https://github.com/golang-migrate/migrate).

To install it, run:

```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Creating migrations

Create some migrations using migrate CLI. Here is an example:

```
migrate create -ext sql -dir database/migrate -seq create_users_table
```

Once you create your files, you should fill them.

**IMPORTANT:** In a project developed by more than one person there is a chance of migrations inconsistency - e.g. two
developers can create conflicting migrations, and the developer that created his migration later gets it merged to the
repository first. Developers and Teams should keep an eye on such cases (especially during code review).
[Here](https://github.com/golang-migrate/migrate/issues/179#issuecomment-475821264) is the issue summary if you would
like to read more.

Consider making your migrations idempotent - we can run the same sql code twice in a row with the same result. This
makes our migrations more robust. On the other hand, it causes slightly less control over database schema - e.g. let's
say you forgot to drop the table in down migration. You run down migration - the table is still there. When you run up
migration again - `CREATE TABLE` would return an error, helping you find an issue in down migration,
while `CREATE TABLE IF NOT EXISTS` would not. Use those conditions wisely.

In case you would like to run several commands/queries in one migration, you should wrap them in a transaction. This way
if one of commands fails, our database will remain unchanged.

### Run database

An example database can be run for development purposes with the following command:

```
make run-db
```

This will execute a Postgres database running in Docker with the following properties:

| Host        | Port   | Username   | Password   | Database   |
| ----------- | ------ | ---------- | ---------- | ---------- |
| `localhost` | `5432` | `postgres` | `postgres` | `postgres` |

### Run migrations

With the DB running, run your migrations through the CLI and check if they applied expected changes. To apply the
database migrations, you can execute:

```
make migrate-db
```

This will apply all the `up` SQL scripts to the database.
