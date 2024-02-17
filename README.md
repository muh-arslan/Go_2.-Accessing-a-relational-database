# Accessing a relational database

## Setup go.work
    go work init
    go work use ./data-access
    go work use ./web-service-gin

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