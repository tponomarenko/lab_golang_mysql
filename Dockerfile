FROM docker.io/ubuntu:latest AS first

RUN apt-get update --yes
RUN apt-get install golang ca-certificates git-core --yes

ADD ./ /phonebook
WORKDIR /phonebook
ENV CGO_ENABLED=0
RUN go build -o /phonebook/app ./

FROM docker.io/alpine:3.17

COPY --from=first /phonebook/app /phonebook

ENTRYPOINT /phonebook
