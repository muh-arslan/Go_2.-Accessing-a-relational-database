# Accessing a relational database

## Setup DB
    mysql -u root -p
    create database recordings;
    use recordings;
    source [/path/to]/create-tables.sql

## How to run
Set env variables:
    set DBUSER=username
    set DBPASS=password

start: 
    go run .