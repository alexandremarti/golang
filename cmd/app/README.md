
# Programa GO que cria novos produtos tanto via API REST quanto via fila kafka

No arquivo cmd/app/main.go está é iniciado tanto o servidor http quanto o consumer da fila kafka

Para rodar o app

go run cmd/app/main.go


Próximos passos
- Organizar para rodar o app no container


## Docker compose para criar o kafka no docker local, caso seja necessário
```
# docker-compose.yml
version: '3.3'
services:
  db:
    image: mysql:5.7
    restart: always
    environment:
      MYSQL_DATABASE: 'db'
      # So you don't have to use root, but you can if you like
      MYSQL_USER: 'user'
      # You can use whatever password you like
      MYSQL_PASSWORD: 'password'
      # Password for root access
      MYSQL_ROOT_PASSWORD: 'password'
    ports:
      # <Port exposed> : <MySQL Port running inside container>
      - '3306:3306'
      # Where our data will be persisted
    volumes:
      - my-db:/var/lib/mysql
# Names our volume
volumes:
  my-db:
 ```

## Cria a base de dados mysql
docker-compose exec db bash

mysql -uroot -ppassword

create database products;
exit;

## Acessa o banco de dados products e cria a tabela products
mysql -uroot -ppassword products
create table products (id varchar(255), name varchar(255), price float);
select * from products;


## Docker compose para criar o kafka no docker local, caso seja necessário
 
 ```
# docker-compose.yml
version: "3.7"
services:
  zookeeper:
    image: docker.io/bitnami/zookeeper:3.8
    ports:
      - "2181:2181"
    volumes:
      - "zookeeper-volume:/bitnami"
    environment:
      - ALLOW_ANONYMOUS_LOGIN=yes
  kafka:
    image: docker.io/bitnami/kafka:3.3
    ports:
      - "9093:9093"
    volumes:
      - "kafka-volume:/bitnami"
    environment:
      - KAFKA_BROKER_ID=1
      - KAFKA_CFG_ZOOKEEPER_CONNECT=zookeeper:2181
      - ALLOW_PLAINTEXT_LISTENER=yes
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CLIENT:PLAINTEXT,EXTERNAL:PLAINTEXT
      - KAFKA_CFG_LISTENERS=CLIENT://:9092,EXTERNAL://:9093
      - KAFKA_CFG_ADVERTISED_LISTENERS=CLIENT://kafka:9092,EXTERNAL://localhost:9093
      - KAFKA_CFG_INTER_BROKER_LISTENER_NAME=CLIENT
    depends_on:
      - zookeeper
volumes:
  kafka-volume:
  zookeeper-volume:
 ```

## Criar o tópico kafka
docker-compose exec kafka bash

kafka-console-producer.sh --bootstrap-server=localhost:9092 --topic=products

## aqui para enviar uma mensagem que resultará no cadastramento via a mensagem recebida
kafka-console-producer.sh --bootstrap-server=localhost:9092 --topic=products
###
colar esse json
{"name": "Meu produto 03 - via kafka","price": 120}