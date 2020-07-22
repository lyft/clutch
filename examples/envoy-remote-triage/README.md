# Envoy Remote Triage Example

This is an example of running the Envoy Remote Triage. Provided are:
- `envoy-config.yaml`: A simple Envoy configuration.
- `clutch-config.yaml`: A Clutch configuration with remote triage enabled.
- `docker-compose.yml`: Container orchestration for this example, bringing up a Clutch and an Envoy.

## Usage
1. `docker-compose up` from a shell in this folder.
1. Visit the Clutch UI at `localhost:8080` in the browser.
1. Click `Envoy - Remote Triage` or select it from the menu.
1. Enter the `envoy` in the text box and select any other options.
1. Submit the form to view remote triage details.
