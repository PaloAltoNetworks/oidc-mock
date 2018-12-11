FROM ubuntu

RUN apt-get upgrade -y && apt-get update -y

RUN mkdir -p /oidc-mock

ADD .data .data
ADD oidcmock /oidc-mock/oidcmock

EXPOSE 6999

ENTRYPOINT ["/oidc-mock/oidcmock"]
