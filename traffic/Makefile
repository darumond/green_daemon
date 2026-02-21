CXX := g++
CXXFLAGS := -std=c++20 -O3
CXXFLAGS += -Wall -Wextra -Wconversion -Wshadow
CXXFLAGS += -Wno-unused-parameter
INCLUDES := -I./include

SRCS := $(wildcard src/*.cc)
HEADERS := $(wildcard include/*.hh)

TARGET := main

all: $(TARGET)

$(TARGET): $(SRCS) $(HEADERS)
	$(CXX) $(CXXFLAGS) $(INCLUDES) $(SRCS) -o $(TARGET)

clean:
	rm -f $(TARGET)

.PHONY: all clean
