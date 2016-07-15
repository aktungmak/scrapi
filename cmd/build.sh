go build \
    -ldflags "-s -w -X 'github.com/aktungmak/scrapi.BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')'" \
    -o scrapi.exe \
    main.go
