# Brute Force Login Simulation

This Go program simulates a brute-force attack on a login API by sending multiple POST requests with randomly generated credentials. The program is designed to handle a high number of requests with configurable concurrency.

⚠️ **Disclaimer:** This code is for educational purposes only. Unauthorized use of brute-force attacks is illegal and unethical. Always ensure you have explicit permission to perform any such actions.

## Features

- Randomly generates usernames and passwords for each request.
- Configurable number of total requests (`numRequests`).
- Configurable concurrency (`concurrency`) to handle multiple requests in parallel.
- Customizable target URL and request headers.
- Monitors and controls the total number of requests sent.

## Requirements

- Go 1.16 or later
