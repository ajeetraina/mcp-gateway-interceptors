FROM python:3.9-alpine
WORKDIR /app
RUN pip install flask
COPY interceptor.py .
EXPOSE 8080
CMD ["python", "interceptor.py"]
