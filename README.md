curl -v 127.0.0.1:10100/v1/cruds/cars/
curl -s -X POST -H "Content-Type: application/json" --data '{ "year":2020 }' 127.0.01:10100/v1/cruds/cars/create
curl -v -X DELETE 127.0.0.1:10100/v1/cruds/cars/"fc579c0e-8b11-46fa-a4f7-d0fcb7d01573"/delete

