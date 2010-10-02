include $(GOROOT)/src/Make.inc

PKGDIR=$(GOROOT)/pkg/$(GOOS)_$(GOARCH)

TARG=gonewrong
CGOFILES=gonewrong.go
CGO_CFLAGS=-I. -I "$(GOROOT)/include"
CGO_LDFLAGS=-lzmq
GOFMT=$(GOROOT)/bin/gofmt -tabwidth=4 -spaces=true -tabindent=false -w 

include $(GOROOT)/src/Make.pkg

CLEANFILES+=clsrv $(PKGDIR)/$(TARG).a

again: clean install

format: 
	$(GOFMT) gonewrong.go
