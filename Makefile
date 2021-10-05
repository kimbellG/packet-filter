
stand-up:
	./up.sh up

test:
	cd cmd/test && go test -tags integration

clean:
	./up.sh down

.PHONY: clean

