# SSHAutomator

The SSHAutomator is a Go program that provides an automation tool for executing scripts or commands on one or more remote servers over SSH.

## Features

- Execute scripts or commands on remote servers using SSH
- Support for multiple servers with parallel execution
- Accept command-line arguments for server details (host, port, username, password/key file)
- Read script and server details from files
- Handle SSH connection and command execution

## Usage

The program accepts the following command-line arguments:

```
shautomator -script script.shautomator -host HOST -port PORT -username USERNAME [-password PASSWORD | -keyfile KEYFILE]
shautomator -script script.shautomator -servers serverlist.txt
```

Where:
- `script.shautomator`: Path to the file containing the commands to be executed
- `HOST`: The SSH host to connect to
- `PORT`: The SSH port (default is 22)
- `USERNAME`: The SSH username
- `PASSWORD`: The SSH password (optional if using a key file)
- `KEYFILE`: The path to the private key file (optional)
- `serverlist.txt`: The file containing the server details (one server per line in the format `HOST,USERNAME,PASSWORD,KEYFILE`)

## Example

To execute a script on a single server:

```
shautomator -script script.shautomator -host example.com -port 22 -username myuser -password mypassword
```

To execute a script on multiple servers from a servers file:

```
shautomator -script script.shautomator -servers serverlist.txt
```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please create an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
