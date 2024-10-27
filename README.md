
# Jailhouse Roll Call

A simple Go application to retrieve and display information about jails from a TrueNAS server.  Makes an API call to fetch jail details, processes the returned data, and serves it in an HTML table on a local web server.

## Prerequisites

- **Go**: Ensure Go is installed on your machine.
- **TrueNAS API Access**: API URL and access key for the TrueNAS instance.
- **Configuration file**: `config.json` (explained below).

## Setup

1. **Clone the repository**:

   ```bash
   git clone git@github.com:jwschman/jrc
   cd jrc
   ```

2. **Create `config.json`**: In the root directory, create a `config.json` file with your TrueNAS API URL and API Key:

   ```json
   {
       "api_url": "http://<TrueNAS-IP>/api/v2.0/jail",
       "api_key": "<API-Key>"
   }
   ```

## Running the Application

1. **Start the Server**:

   ```bash
   go run main.go
   ```

2. **Access the Application**: Open a browser and go to [http://localhost:3000](http://localhost:3000).

## File Structure

- **main.go**: Main application file containing the primary functionality.
- **config.json**: Configuration file with API URL and API key.
- **templates/index.html**: HTML template for rendering jail data.
- **static/stylesheets/main.css**: CSS file for styling the HTML template.

## Acknowledgements

Parts of `index.html` and all of `main.css` were written with assistance from ChatGPT.
