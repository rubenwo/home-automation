FROM python:3.7

RUN pip install fastapi uvicorn PyP100 pyyaml

EXPOSE 80

COPY . /

ENTRYPOINT  ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "80"]