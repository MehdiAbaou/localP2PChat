# localP2PChat

A lightweight, peer-to-peer (P2P) chat client for local networks, featuring a responsive Terminal User Interface (TUI). Built in Go with a focus on simplicity and zero-configuration networking.

---

## Overview

**localP2PChat** enables instant communication between devices on a shared local network without the need for a central server. It leverages UDP multicast for peer discovery and a reactive TUI for an efficient CLI-based messaging experience.

### Key Features

*   **Zero-Config Discovery**: Automatically finds and connects to peers on the local network using UDP multicast.
*   **Reactive TUI**: A polished terminal interface built with `bubbletea` and `lipgloss`, supporting stateful transitions between login and chat views.
*   **Real-Time Messaging**: Asynchronous message handling via Go routines for seamless background updates.
*   **Integrity Verification**: Uses MD5 hashing to verify peer identities and ensure message origin consistency.
*   **Peer Lifecycle Tracking**: Dynamically manages the active peer list as users join or leave the network.

---

## Getting Started

### Prerequisites

*   **Go**: Version 1.18 or higher.
*   **Dependencies**: The project utilizes the following libraries:
    *   `[github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)`
    *   `[github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)`

### Installation

1.  **Clone the repository then**:
    ```bash
    cd localP2PChat
    ```

2.  **Build or Run**:
    To build the binary:
    ```bash
    go build -o p2pchat
    ```
    To run directly:
    ```bash
    go run .
    ```

> [!TIP]
> Ensure your firewall allows UDP traffic on the local network to enable the automatic peer discovery feature.

---

## Technical Details

### Architecture
The client operates as both a server and a client (servent). It listens for incoming UDP packets on a dedicated goroutine while simultaneously handling user input via the TUI main loop.



### Roadmap & Status

| Feature | Status | Description |
| :--- | :--- | :--- |
| **UDP Multicast** | 🟡 In Progress | Automatic peer discovery via shared multicast addresses. |
| **IP Verification** | Completed | MD5 hashing for basic IP-based message integrity. |
| **Real-Time Logic** | Completed | Background listener for asynchronous chat updates. |
| **Lifecycle Mgmt** | 🟡 In Progress | Dynamic signaling for joining/leaving the network and a periodic "are-you-alive" signal. |
| **Stateful Layouts** | Completed | Switching between Login and Chat views. |
| **Buffer Control** | Completed | Scrolling message history (last 10 messages). |

---

## Usage

1.  **Launch** the application.
2.  **Enter your username** at the login prompt.
3.  **Start chatting!** Messages are broadcasted to all discovered peers on your local subnet.

> [!NOTE]
> This tool is designed for **local-only** communication. It does not route traffic over the internet or through external gateways.
```