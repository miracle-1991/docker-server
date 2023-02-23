FROM pytorch/pytorch

WORKDIR /home

COPY . .

RUN apt-get update -y && \
    apt-get install -y libglib2.0-0 libsm6 libxext6 libxrender-dev libgl1-mesa-dev && \
    pip3 install flask && \
    pip3 install flask_cors && \
    pip3 install geopy && \
    pip3 install easyocr && \
    pip3 install keplergl && \
    pip3 install boto3 && \
    pip3 install readerwriterlock 

CMD [ "/opt/conda/bin/python" "/home/server.py" ]