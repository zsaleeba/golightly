# Makefile for building Lua
# See ../doc/readme.html for installation and customization instructions.

# == CHANGE THE SETTINGS BELOW TO SUIT YOUR ENVIRONMENT =======================

# Your platform. See PLATS for possible values.
PLAT= none

CC= gcc
CFLAGS= -O2 -Wall -DLUA_COMPAT_ALL $(SYSCFLAGS) $(MYCFLAGS)
LDFLAGS= $(SYSLDFLAGS) $(MYLDFLAGS)
LIBS= -lm $(SYSLIBS) $(MYLIBS)

AR= ar rcu
RANLIB= ranlib
RM= rm -f

SYSCFLAGS=
SYSLDFLAGS=
SYSLIBS=

MYCFLAGS=
MYLDFLAGS=
MYLIBS=
MYOBJS=

# == END OF USER SETTINGS -- NO NEED TO CHANGE ANYTHING BELOW THIS LINE =======

PLATS= aix ansi bsd freebsd generic linux macosx mingw posix solaris

GLI_A=	libgolightly.a
GLI_CORE_O=lapi.o lcode.o lctype.o ldebug.o ldo.o ldump.o lfunc.o lgc.o llex.o \
	lmem.o lobject.o lopcodes.o lparser.o lstate.o lstring.o ltable.o \
	ltm.o lundump.o lvm.o lzio.o
GLI_LIB_O=lauxlib.o lbaselib.o lbitlib.o lcorolib.o ldblib.o liolib.o \
	lmathlib.o loslib.o lstrlib.o ltablib.o loadlib.o linit.o
GLI_BASE_O=$(GLI_CORE_O) $(GLI_LIB_O) $(MYOBJS)

GLI_T=	gli
GLI_O=	gli.o

GLC_T=	glc
GLC_O=	glc.o

ALL_O= $(GLI_BASE_O) $(GLI_O) $(GLC_O)
ALL_T= $(GLI_A) $(GLI_T) $(GLC_T)
ALL_A= $(GLI_A)

# Targets start here.
default: $(PLAT)

all:	$(ALL_T)

o:	$(ALL_O)

a:	$(ALL_A)

$(GLI_A): $(GLI_BASE_O)
	$(AR) $@ $(GLI_BASE_O)
	$(RANLIB) $@

$(GLI_T): $(GLI_O) $(GLI_A)
	$(CC) -o $@ $(LDFLAGS) $(GLI_O) $(GLI_A) $(LIBS)

$(GLC_T): $(GLC_O) $(GLI_A)
	$(CC) -o $@ $(LDFLAGS) $(GLC_O) $(GLI_A) $(LIBS)

clean:
	$(RM) $(ALL_T) $(ALL_O)

depend:
	@$(CC) $(CFLAGS) -MM l*.c

echo:
	@echo "PLAT= $(PLAT)"
	@echo "CC= $(CC)"
	@echo "CFLAGS= $(CFLAGS)"
	@echo "LDFLAGS= $(SYSLDFLAGS)"
	@echo "LIBS= $(LIBS)"
	@echo "AR= $(AR)"
	@echo "RANLIB= $(RANLIB)"
	@echo "RM= $(RM)"

# Convenience targets for popular platforms
ALL= all

none:
	@echo "Please do 'make PLATFORM' where PLATFORM is one of these:"
	@echo "   $(PLATS)"

aix:
	$(MAKE) $(ALL) CC="xlc" CFLAGS="-O2 -DUSE_POSIX -DUSE_DLOPEN" SYSLIBS="-ldl" SYSLDFLAGS="-brtl -bexpall"

ansi:
	$(MAKE) $(ALL) SYSCFLAGS="-DLUA_ANSI"

bsd:
	$(MAKE) $(ALL) SYSCFLAGS="-DUSE_POSIX -DUSE_DLOPEN" SYSLIBS="-Wl,-E"

freebsd:
	$(MAKE) $(ALL) SYSCFLAGS="-DUSE_LINUX" SYSLIBS="-Wl,-E -lreadline"

generic: $(ALL)

linux:
	$(MAKE) $(ALL) SYSCFLAGS="-DUSE_LINUX" SYSLIBS="-Wl,-E -ldl -lreadline"

macosx:
	$(MAKE) $(ALL) SYSCFLAGS="-DUSE_MACOSX" SYSLIBS="-lreadline"

mingw:
	$(MAKE) "LUA_A=lua52.dll" "LUA_T=lua.exe" \
	"AR=$(CC) -shared -o" "RANLIB=strip --strip-unneeded" \
	"SYSCFLAGS=-DLUA_BUILD_AS_DLL" "SYSLIBS=" "SYSLDFLAGS=-s" lua.exe
	$(MAKE) "LUAC_T=luac.exe" luac.exe

posix:
	$(MAKE) $(ALL) SYSCFLAGS="-DUSE_POSIX"

solaris:
	$(MAKE) $(ALL) SYSCFLAGS="-DUSE_POSIX -DUSE_DLOPEN" SYSLIBS="-ldl"

# list targets that do not create files (but not all makes understand .PHONY)
.PHONY: all $(PLATS) default o a clean depend echo none

# DO NOT DELETE

lapi.o: lapi.c golightly.h gliconf.h lapi.h llimits.h lstate.h lobject.h ltm.h \
 lzio.h lmem.h ldebug.h ldo.h lfunc.h lgc.h lstring.h ltable.h lundump.h \
 lvm.h
lauxlib.o: lauxlib.c golightly.h gliconf.h lauxlib.h
lbaselib.o: lbaselib.c golightly.h gliconf.h lauxlib.h lualib.h
lbitlib.o: lbitlib.c golightly.h gliconf.h lauxlib.h lualib.h
lcode.o: lcode.c golightly.h gliconf.h lcode.h llex.h lobject.h llimits.h \
 lzio.h lmem.h lopcodes.h lparser.h ldebug.h lstate.h ltm.h ldo.h lgc.h \
 lstring.h ltable.h lvm.h
lcorolib.o: lcorolib.c golightly.h gliconf.h lauxlib.h lualib.h
lctype.o: lctype.c lctype.h golightly.h gliconf.h llimits.h
ldblib.o: ldblib.c golightly.h gliconf.h lauxlib.h lualib.h
ldebug.o: ldebug.c golightly.h gliconf.h lapi.h llimits.h lstate.h lobject.h \
 ltm.h lzio.h lmem.h lcode.h llex.h lopcodes.h lparser.h ldebug.h ldo.h \
 lfunc.h lstring.h lgc.h ltable.h lvm.h
ldo.o: ldo.c golightly.h gliconf.h lapi.h llimits.h lstate.h lobject.h ltm.h \
 lzio.h lmem.h ldebug.h ldo.h lfunc.h lgc.h lopcodes.h lparser.h \
 lstring.h ltable.h lundump.h lvm.h
ldump.o: ldump.c golightly.h gliconf.h lobject.h llimits.h lstate.h ltm.h \
 lzio.h lmem.h lundump.h
lfunc.o: lfunc.c golightly.h gliconf.h lfunc.h lobject.h llimits.h lgc.h \
 lstate.h ltm.h lzio.h lmem.h
lgc.o: lgc.c golightly.h gliconf.h ldebug.h lstate.h lobject.h llimits.h ltm.h \
 lzio.h lmem.h ldo.h lfunc.h lgc.h lstring.h ltable.h
linit.o: linit.c golightly.h gliconf.h lualib.h lauxlib.h
liolib.o: liolib.c golightly.h gliconf.h lauxlib.h lualib.h
llex.o: llex.c golightly.h gliconf.h lctype.h llimits.h ldo.h lobject.h \
 lstate.h ltm.h lzio.h lmem.h llex.h lparser.h lstring.h lgc.h ltable.h
lmathlib.o: lmathlib.c golightly.h gliconf.h lauxlib.h lualib.h
lmem.o: lmem.c golightly.h gliconf.h ldebug.h lstate.h lobject.h llimits.h \
 ltm.h lzio.h lmem.h ldo.h lgc.h
loadlib.o: loadlib.c golightly.h gliconf.h lauxlib.h lualib.h
lobject.o: lobject.c golightly.h gliconf.h lctype.h llimits.h ldebug.h lstate.h \
 lobject.h ltm.h lzio.h lmem.h ldo.h lstring.h lgc.h lvm.h
lopcodes.o: lopcodes.c lopcodes.h llimits.h golightly.h gliconf.h
loslib.o: loslib.c golightly.h gliconf.h lauxlib.h lualib.h
lparser.o: lparser.c golightly.h gliconf.h lcode.h llex.h lobject.h llimits.h \
 lzio.h lmem.h lopcodes.h lparser.h ldebug.h lstate.h ltm.h ldo.h lfunc.h \
 lstring.h lgc.h ltable.h
lstate.o: lstate.c golightly.h gliconf.h lapi.h llimits.h lstate.h lobject.h \
 ltm.h lzio.h lmem.h ldebug.h ldo.h lfunc.h lgc.h llex.h lstring.h \
 ltable.h
lstring.o: lstring.c golightly.h gliconf.h lmem.h llimits.h lobject.h lstate.h \
 ltm.h lzio.h lstring.h lgc.h
lstrlib.o: lstrlib.c golightly.h gliconf.h lauxlib.h lualib.h
ltable.o: ltable.c golightly.h gliconf.h ldebug.h lstate.h lobject.h llimits.h \
 ltm.h lzio.h lmem.h ldo.h lgc.h lstring.h ltable.h lvm.h
ltablib.o: ltablib.c golightly.h gliconf.h lauxlib.h lualib.h
ltm.o: ltm.c golightly.h gliconf.h lobject.h llimits.h lstate.h ltm.h lzio.h \
 lmem.h lstring.h lgc.h ltable.h
lua.o: lua.c golightly.h gliconf.h lauxlib.h lualib.h
luac.o: luac.c golightly.h gliconf.h lauxlib.h lobject.h llimits.h lstate.h \
 ltm.h lzio.h lmem.h lundump.h ldebug.h lopcodes.h
lundump.o: lundump.c golightly.h gliconf.h ldebug.h lstate.h lobject.h \
 llimits.h ltm.h lzio.h lmem.h ldo.h lfunc.h lstring.h lgc.h lundump.h
lvm.o: lvm.c golightly.h gliconf.h ldebug.h lstate.h lobject.h llimits.h ltm.h \
 lzio.h lmem.h ldo.h lfunc.h lgc.h lopcodes.h lstring.h ltable.h lvm.h
lzio.o: lzio.c golightly.h gliconf.h llimits.h lmem.h lstate.h lobject.h ltm.h \
 lzio.h

