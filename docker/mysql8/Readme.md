dump mysql
```
docker exec mysql_8 sh -c 'exec /usr/bin/mysqldump -u root --password=123456 nbp > /tmp/dump/nbp.sql'
```

resort(Haven't try may have bug)
```
docker exec -i mysql_8 sh -c 'exec mysql -u root -p123456 nbp < nbp.sql'
```