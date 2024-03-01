FROM debian:stable-slim

# copy source destination.
# don't need the additional functionality of add
COPY dockertest /bin/dockertest
ENV PORT="8080"

CMD ["/bin/dockertest"]
