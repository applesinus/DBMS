Commands format: command, then args splitted by space
Example: "set year 2024 in pool.scheme.collection"

Commands be like:


```
createPool <name>
deletePool <name>

createScheme <name> in <pool name>
deleteScheme <name> in <pool name>

createCollection <name> in <pool name>.<scheme name>
deleteCollection <name> in <pool name>.<scheme name>

set <key> <value> in <pool name>.<scheme name>.<collection name>
update <key> <value> in <pool name>.<scheme name>.<collection name>
get <key> in <pool name>.<scheme name>.<collection name>
getRange <key minimum> <key maximum> in <pool name>.<scheme name>.<collection name>
delete <key> in <pool name>.<scheme name>.<collection name>
```
