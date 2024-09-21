FROM postgres:14

# Time ZoneAc
ENV TZ Asia/Tokyo

#DDLとなるファイルをコピー
#dockerのpostgresだと/docker-entrypoint-initdb.d配下のモノを実行してくれるらしい
#https://qiita.com/taka_22/items/5991e28ae94035710e0f
COPY /sql/*.sql /docker-entrypoint-initdb.d

#COPY ../*.sh /