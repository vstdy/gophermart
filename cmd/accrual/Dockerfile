FROM alpine:3.15

RUN apk add --no-cache libc6-compat

WORKDIR /accrual
COPY ./accrual_linux_amd64 .
RUN chmod +x accrual_linux_amd64

CMD ["./accrual_linux_amd64"]
