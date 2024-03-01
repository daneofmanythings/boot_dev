# define the linux image
FROM debian:stable-slim

# copy the code
COPY main.py main.py

# copy the data dependencies
COPY books/ books/

# update the container
RUN apt update
RUN apt upgrade -y

# install build tooling
RUN apt install -y build-essential zlib1g-dev libncurses5-dev libgdbm-dev libnss3-dev libssl-dev libreadline-dev libffi-dev libsqlite3-dev wget libbz2-dev

# get python interpreter
RUN wget https://www.python.org/ftp/python/3.10.8/Python-3.10.8.tgz
RUN tar -xf Python-3.10.*.tgz

# build the python interpreter
RUN cd Python-3.10.8 && ./configure --enable-optimizations && make && make altinstall

# finally run it
CMD ["python3.10", "main.py"]
