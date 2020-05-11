#!/usr/bin/env bash

set -a; source .env; set +a
set -e


CMD="docker-compose -f docker-compose.yml"

if [[ ${ENVIRONMENT} != "production" ]]; then
  CMD="${CMD} -f docker-compose.development.yml"
fi

CMD="${CMD} -f docker-compose.${JITSI42_CONSUMER_TYPE:-api}.yml"

${CMD} ${@}
