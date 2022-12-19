COMPILER = go
EXE = solat.exe
SRC = $(wildcard *.go)

$(EXE):
	$(COMPILER) build -o $(EXE) $(SRC)

run:$(EXE)
	@./$< ||:

clean:
	rm -rf $(EXE)