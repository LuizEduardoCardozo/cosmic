FILES=$(shell ls *.md)
EXE=cosmic

build:
	go build -o ${EXE}

clean:
	rm -f *.svg *.gv	

run:
	./${EXE} ${FILES}

svg:
	dot -Tsvg -Kneato -O *.gv

gen:
	make clean
	make build
	make run
	make svg
