@echo off
chcp 1251

SET server=localhost
SET database=postgres
SET port=5432
SET username=postgres

"C:\Program Files\PostgreSQL\15\bin\psql.exe" -h %server% -U %username% -d %database% -p %port%