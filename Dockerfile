FROM gcr.io/distroless/static-debian12

COPY artifacts/svc /svc

EXPOSE 8080

WORKDIR /

USER nonroot:nonroot

CMD ["./svc"]
