TARGET	= gl
SRCS	= gl.c ast.c
GENSRCS	= golexer.yy.c goparser.tab.c
GENHDRS	= goparser.tab.h
OBJS	:= $(SRCS:%.c=%.o) $(GENSRCS:%.c=%.o)
CFLAGS  = -Wall -g
LIBS    = -ljit -ldl -lm -lpthread

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) -o $(TARGET) $(OBJS) $(LIBS)

golexer.yy.c golexer.yy.h: golexer.l
	flex golexer.l

goparser.tab.c goparser.tab.h: goparser.y
	bison -b goparser goparser.y

clean:
	rm -f $(TARGET) $(OBJS) $(GENSRCS) $(GENHDRS)

golexer.yy.o: golexer.yy.c goparser.tab.h
