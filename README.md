# price-crawler

This is a project aimed at crawling historical price data from various sources. The data is then stored in a data bucket for further processing.

# Architecture

The project uses a microservices architecture. The main components are:

- **parser** - This service is responsible for parsing the data from the source into structured data.
- **requester** - This service is responsible for byparsing any restrictions set by the source (like human verification, cookies) and fetching the data.
- **data-bucket** - This service is responsible is the interface to the data storage. It provides endpoints to store and retrieve data.
- **watcher** - This service provides endpoints to manually update a given aimed product's price data. It also automatically updates the data every set interval.
- 🚧 **notification** - This service provides endpoints to set alerts on a given product. It sends an notification when the price of the product falls below a certain threshold. 
- 🚧 **viewer** - This is a web-based interface to view the historical chart of the product's price data.