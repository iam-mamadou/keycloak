FROM centos7

ARG KEYCLOAK_VERSION=11.0.2
ARG KEYCLOAK_DIR=/opt/keycloak-${KEYCLOAK_VERSION}

ARG POSTGRES_DRIVER_VERSION=42.2.18
ARG POSTGRES_MODULE_DIR=${KEYCLOAK_DIR}/modules/system/layers/keycloak/org/postgresql/main

RUN mkdir -p ${POSTGRES_MODULE_DIR} 

RUN yum -y update && \
    yum -y install wget \
                   unzip \
                   java-1.8.0-openjdk \
                   less

# download keycloak
RUN wget https://downloads.jboss.org/keycloak/${KEYCLOAK_VERSION}/keycloak-${KEYCLOAK_VERSION}.zip && \
    unzip keycloak-${KEYCLOAK_VERSION}.zip -d /opt/ && \
    rm keycloak-${KEYCLOAK_VERSION}.zip 

# download postgres driver
RUN wget https://jdbc.postgresql.org/download/postgresql-${POSTGRES_DRIVER_VERSION}.jar \
    -O ${POSTGRES_MODULE_DIR}/postgresql.jar 

# add files
COPY ./entrypoint.sh /entrypoint.sh 
COPY ./configuration/standalone.xml ${KEYCLOAK_DIR}/standalone/configuration/standalone.xml
COPY ./modules/postgres/module.xml ${POSTGRES_MODULE_DIR}/module.xml
COPY ./setup_keycloak.sh /setup_keycloak.sh

# create symbolic link
RUN ln -s ${KEYCLOAK_DIR} /opt/keycloak

EXPOSE 8080

ENTRYPOINT [ "/entrypoint.sh" ]

