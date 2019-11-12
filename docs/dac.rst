.. _database-as-code:

Database as Code
================

In the DDSL allows DDL (Data Definition Language) to be managed like any other
software code.

A RDS (Relational Database System) is a *runtime* environment for DDL.
Just as app source code would not be managed directly in its deployed, runtime
environment, so too DDL code should not be managed inside the RDS.

DDSL provides the Database DevOps tool to enable CICD (continuous integration,
continuous delivery) of database DDL code from source to RDS.

Migrations
----------

