#include <iostream>

#include "../include/server.hh"

int main(int argc, char **argv) {
  if (std::string(argv[1]) == "server") {
    std::cout << "running server\n";
    serverMain();
  } else {
    std::cout << "running client\n";
    clientMain();
  }
}
