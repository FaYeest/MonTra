# syntax=docker/dockerfile:1
FROM python:3.12-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1

RUN groupadd -r montra && useradd -r -g montra -d /app montra

WORKDIR /app

COPY pyproject.toml ./
COPY app ./app
COPY training ./training
COPY evaluation ./evaluation
COPY artifacts ./artifacts

RUN pip install --upgrade pip && \
    pip install . && \
    chown -R montra:montra /app

USER montra
EXPOSE 8000

CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
