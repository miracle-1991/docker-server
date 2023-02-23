FROM drogonframework/drogon

WORKDIR /home

COPY . .

RUN apt-get update -y \
    && apt-get install gdb -y \
    && apt-get install vim -y \
    && rm -rf build \
    && mkdir build \
    && cd build \
    && cmake .. \
    && make \
    && mv rtk_service .. 

CMD [ "/home/rtk_service" ]
