FROM dcarbon/go-shared as builder

WORKDIR /dcarbon/projects
COPY . .

RUN go mod tidy && go build -buildvcs=false -o projects && \
    cp  projects /usr/bin && \
    echo "Build image successs...!"


FROM dcarbon/dimg:minimal
COPY --from=builder /usr/bin/projects /usr/bin/projects

CMD [ "projects" ]