SASSC   ?= sassc
GO      ?= go

GO_SRC   = server/$(wildcard *.go)

override SASSCFLAGS += -t compressed

all: CrashDragon upload_syms minidump-stackwalk/stackwalker

CrashDragon: $(GO_SRC) assets/stylesheets/app.css
	$(GO) build -o bin/CrashDragon $(GO_SRC)

upload_syms: upload_syms/main.go
	$(GO) build -o bin/upload_syms upload_syms/main.go

assets/stylesheets/app.css:
	$(SASSC) $(SASSCFLAGS) $(@D)/app.scss > $@.tmp && mv $@.tmp $@

minidump-stackwalk/stackwalker:
	cd minidump-stackwalk && make

clean:
	rm -rf bin/
	rm -f assets/stylesheets/app.css
	cd minidump-stackwalk && make distclean

.PHONY: clean all
