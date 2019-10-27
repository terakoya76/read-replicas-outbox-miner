FROM node:12

RUN mkdir /kinesis-data-streams
WORKDIR /kinesis-data-streams

# Avoid `user "nobody" does not have permission to access`
# see: https://github.com/nodejs/node-gyp/issues/1236
RUN npm -g install kinesalite --unsafe

EXPOSE 4567
CMD ["kinesalite", "--path", "/kinesis-data-streams/data"]
