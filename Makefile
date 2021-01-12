es-cli: *.go cmd/*.go go.mod go.sum
	go build

clean:
	rm es-cli

install: es-cli
	cp es-cli ${HOME}/bin/Darwin

completion: install
	es-cli completion bash > ${HOME}/.profile.d/es-cli.sh
