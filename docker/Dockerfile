FROM alpine

RUN apk update && apk --no-cache add curl && apk add --no-cache bash

RUN mkdir -p /oidc-mock

ADD .data .data
ADD oidcmock /oidc-mock/oidcmock

EXPOSE 6999

ENTRYPOINT ["/oidc-mock/oidcmock"]
