FROM ubuntu

WORKDIR /srv

COPY ./build /srv/
RUN apt-get update && apt-get install wget -y

RUN chmod +x sai-eth-interaction
CMD ./sai-eth-indexer start


EXPOSE 8018
EXPOSE 8080
EXPOSE 8019
