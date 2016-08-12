FROM scratch
WORKDIR /app
COPY consulreaper /app/
EXPOSE 8080
ENTRYPOINT ["./consulreaper"]