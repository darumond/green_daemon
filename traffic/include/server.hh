#ifndef SERVER_HH
#define SERVER_HH

#include <arpa/inet.h>
#include <iostream>
#include <netinet/in.h>
#include <sys/socket.h>
#include <unistd.h>

static int serverMain() {
  // 1. Create TCP socket
  int server_fd = socket(AF_INET, SOCK_STREAM, 0);

  // 2. Bind to Port 8000
  sockaddr_in addr{AF_INET, htons(8000), {INADDR_ANY}};
  bind(server_fd, (sockaddr *)&addr, sizeof(addr));

  // 3. Listen and Accept incoming connection
  listen(server_fd, 1);
  std::cout << "Listening on port 8000..." << std::endl;
  int client_fd = accept(server_fd, nullptr, nullptr);

  // 4. Read traffic in a loop
  char buffer[128];
  while (read(client_fd, buffer, sizeof(buffer)) > 0) {
    std::cout << "Received 128 bytes: " << buffer << std::endl;
  }

  close(client_fd);
  close(server_fd);
  return 0;
}

int clientMain() {
  // 1. Create TCP socket
  int sock_fd = socket(AF_INET, SOCK_STREAM, 0);

  // 2. Connect to localhost:8000
  sockaddr_in addr{AF_INET, htons(8000)};
  addr.sin_addr.s_addr = inet_addr("127.0.0.1");
  connect(sock_fd, (sockaddr *)&addr, sizeof(addr));

  // 3. Exactly 128 bytes of payload (C-string padded with nulls up to 128)
  char msg[128] =
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod "
      "tempor incididunt ut labore et dolore magna aliqua. Ut ";

  // 4. Send traffic in a loop
  while (true) {
    write(sock_fd, msg, sizeof(msg));
    sleep(1); // Optional: prevents terminal flooding
  }

  close(sock_fd);
  return 0;
}

#endif /* SERVER_HH */
