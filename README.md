# rdpctl

`rdpctl` is a command-line interface (CLI) tool written in Go for securely managing your RDP connection profiles and launching FreeRDP (`xfreerdp`) sessions.

It provides an encrypted vault to store your RDP host, username, domain, and optionally, passwords, all protected by a master password.

## Usage

To build and run `rdpctl`:

1.  Navigate to the `rdpctl` directory:
    ```bash
    cd rdpctl
    ```
2.  Build the application:
    ```bash
    go build -trimpath -ldflags "-s -w"
    ```
3.  Run the application:
    ```bash
    ./rdpctl
    ```

### First-Time Setup

Upon your first run, `rdpctl` will detect that no vault exists and guide you through creating a new master password and an empty vault. Follow the on-screen prompts to set up your vault and add your first RDP connection.
