FROM python:3

ENV APP_DIR /app 
RUN mkdir ${APP_DIR}

WORKDIR ${APP_DIR}

COPY . .

RUN pip install -r requirements.txt