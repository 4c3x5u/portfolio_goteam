# goteam/db
This directory includes SQL scripts that are executed against the GoTeam! 
database. The `init.sql` script contains the full initial schema of the 
database.

Should a change to the database schema be necessary, a file with the name 
`migration-<yyyy>-<mm>-<dd>-<hh>-<mm>.sql` must be created in this directory
containing the schema alteration code and executed against the database to 
keep track of changes to the schema.
