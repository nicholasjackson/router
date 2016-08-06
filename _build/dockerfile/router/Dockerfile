FROM nicholasjackson/microservice-basebox

EXPOSE 8001

# Create directory for server files
RUN mkdir /router

# Add s6 config
ADD s6-etc /etc/s6
RUN chmod -R 755 /etc/s6; \
chmod -R 755 /etc/s6

# Add consul template
ADD config.ctmpl /router/config.ctmpl

# Add server files
ADD swagger_spec /swagger
ADD router /router/router

RUN chmod 755 /router/router

ENTRYPOINT ["/usr/bin/s6-svscan","/etc/s6"]
CMD []
