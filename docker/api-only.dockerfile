#include "backend.docker"

FROM alpine:3

#include "user.docker"

WORKDIR /var/lib/api-rback

COPY --from=build --chown=api-rback:api-rback /src/bin/rbns .

RUN chmod +x rbns \
    && mv rbns /usr/local/bin/

USER api-rback

CMD [ "rbns" ]