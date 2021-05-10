FROM python:3.7

COPY ./req.txt /
RUN pip install -r req.txt

EXPOSE 80

COPY . /

ENTRYPOINT  ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "80"]