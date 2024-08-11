# price-crawler-databucket

This is a data bucket for the price crawler project. It contains the data that is used by the price crawler project.

## Endpoints

- **GET** `/datapoints` - Get all the data points filtered by the query parameters.
    - Query Parameters:
        - `name` - The name of the data point.
        - `start` - The start date of the data point in RFC3339 format.
        - `end` - The end date of the data point in RFC3339 format.

- **POST** `/datapoints` - Create a new data point.
    - Request Body:
        - `name` - The name of the data point.
        - `value` - The value of the data point.
        - `timestamp` - The timestamp of the data point in RFC3339 format.

- **GET** `/datapoints/name` - Get all the unique data point names.
