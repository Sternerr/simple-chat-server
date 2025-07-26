# Requirements Specification
Version: v0.9.0

---

## 1. Overview
### 1.1 Purpose
This document describes the requiremnts for a terminal-based chat application designed for real time communication between multiple users.

### 1.2 Definitions
- **TUI**: A text based user interface rendered in the terminal.
- **CLI**: a command line interface.
- **Client**: A user program that connects to the server, send messages, and recieves messages from other clients.
- **Server**: Manages client connections and broadcast messages to all connected clients.
- **Session**: Active connection between a client and a server.
- **TCP**: A standard network communication protocol that enables reliable, ordered, and error-checked delivery of data between applications over a network.

### 1.2 Scope
The software provides a terminal-based platform for users to exchange text messages in real time. The system includes
- a TCP based chat server
- a client with two interface mode:
    - TUI.
    - CLI.
    - Real time messaging with support for multiple users.
    - Ephemeral (non-persistent) message storage.

---

## 2. Overview
### 2.1 Assumptions and Dependencies
- Users can aceess the terminal and run Go binaries, with optional TUI rendering.
- The network allows TCP communcation over a selected port.
- The TUI mode depends on a third-party Go library Bubbletea, which must be compatible with the user's environment.

### 2.2 Constraints
- Go 1.24.4 or newer is required.
- The application does not support encrypted communication in this beta version.
- Server must be started before clients connects

---

## 3 Functional Requirements

| ID | Requirements|
|---|---|
| F-001 | The SuD shall maintain active user sessions until disconnect |
| F-002 | The SuD shall notify users of conection status changes |
| F-003 | The SuD shall enable users to use a unique display name |
| F-004 | The SuD shall enable users to send and recieve text messaging in real time |
| F-005 | The SuD shall store a configurable number of recent messages in the ephemeral storage (default: 50) |
| F-006 | The SuD shall provide an optional user interface (cli, tui) |
| F-007 | The TUI shall display chat window, input field |
| F-008 | The TUI shall suport keyboard shortcuts |
| F-009 | The CLI shall display messages as plain text in the terminal |

