Refactoring guide (What I did in no particular order):
- Created a Makefile to handle my commands.
- set up `sqlc` and `testify` for generating Golang code from SQL and testing my methods and functions respectively.
- set up golang migrate and created migration folder.
- redesigned the database schema using `dbdiagram.io`, exported to PostgreSQL and ran the migration file for it.
- created `util` package for utilities such as generating random numbers, mobile numbers, strings, names, emails and passwords (with a test package for password method).
- set up the tests for db connection, and `user.sql.go` which creates a new user account, gets a new user account by ID, updates the existing information of the user's account and deletes the user account.