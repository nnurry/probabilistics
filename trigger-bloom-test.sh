/usr/local/go/bin/go clean -cache

/usr/local/go/bin/go test -v -timeout 140s -run ^TestClassicBloomBasic$ github.com/nnurry/probabilistics/v2/test > TestClassicBloomBasic.txt
/usr/local/go/bin/go test -v -timeout 140s -run ^TestCountingBloomBasic$ github.com/nnurry/probabilistics/v2/test > TestCountingBloomBasic.txt

