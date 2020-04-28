# swoop
moves tables from one database to another

```
swoop run -f task.yaml
```

task.yaml
```
source:
  name: erp-database
  flavor: mysql
  connectionURL: "user:password@(host)/database?parseTime=true"
destination:
  name: enterprise-data-warehouse
  flavor: mssql
  connectionURL: "sqlserver://user:password@host?database=databse"
secretProvider:
  flavor: raw
truncate: true
writeLimit: 200
tables:
    - source: table1
      destination: myOtherTable
    - source: table2
      destination: woah_another_table
```

or it even supports vault

```
source:
  name: erp-database
  flavor: mysql
destination:
  name: enterprise-data-warehouse
  flavor: mssql
secretProvider:
  flavor: vault
  connectionURL: http://my_vault_server:8200
truncate: true
writeLimit: 200
tables:
    - source: table1
      destination: myOtherTable
    - source: table2
      destination: woah_another_table
```

## TODO

- [ ] add postgres support
- [ ] add snowflake support
