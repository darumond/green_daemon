#include <iostream>

#include "../include/server.hh"

int main(int argc, char **argv) {
  std::string mode = argv[1];
  if (mode == "server") {
    std::cout << "running server\n";
    serverMain();
  } else if (mode == "client") {
    std::cout << "running client\n";
    clientMain();
  } else {
    std::cout << "running thrash\n";
    clientMain();
  }
}
