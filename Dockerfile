# temp container 
FROM golang:alpine as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
# RUN go test
RUN go build -o auth-rest-api .
# Final build with minimal FS
FROM golang:alpine as finalBuild
WORKDIR /app
COPY --from=builder /app/auth-rest-api /app/auth-rest-api
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/resources /app/resources
COPY --from=builder /app/.env /app/.env
CMD ["/app/auth-rest-api"]
