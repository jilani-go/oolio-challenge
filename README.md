# Glofox Fitness Class Booking System

A RESTful API service for managing fitness class bookings. This system allows studio owners to create classes and members to book those classes.

### Features

- Create fitness classes with capacity
- Allow members to Book classes
- In-memory data storage (for demonstration purposes)
- RESTful API interface
- OpenAPI specification for API documentation


### Prerequisites

- Go 1.24
- Git

### Setup

1. Clone the repository:
   ```
   git clone https://github.com/jilani-go/glofox.git
   cd glofox
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Build and run the application:
   ```
   make build
   make run
   ```

## API Documentation

The API is documented using OpenAPI Specification (OAS) which provides a standardized way to describe RESTful APIs. 

### Using the OpenAPI Specification

The project includes an `openapi.yaml` file in the root directory that contains the complete API specification. You can import the file directly in postman to play with API.

### Using cURL request
 - #### Create class sample request
      ```
      curl --location 'http://localhost:8080/api/classes' \
      --header 'Content-Type: application/json' \
      --header 'Accept: application/json' \
      --data '{
          "name": "Gym starter",
          "start_date": "2025-05-01",
          "end_date": "2025-05-30",
          "capacity": 1
      }'
      
      ```
 - #### Create class sample response
   ```
   {
     "id":"50d06c6f-d1cb-496e-beae-5c3d1c94a245",
      "message":"Class created successfully"
   }
   
   ```
- #### Create Booking sample request

   ```
   curl --location 'http://localhost:8080/api/bookings' \
   --header 'Content-Type: application/json' \
   --header 'Accept: application/json' \
   --data '{
       "class_id": "50d06c6f-d1cb-496e-beae-5c3d1c94a245",
       "member_name": "Jilani",
       "date": "2025-05-02"
   }'
   ```
- #### Create Booking sample response
   ```
   {
     "id":"78fd2762-8f3a-4ea1-acb2-25fb13dbca68",
     "message":"Booking created successfully"
    }

  
   ```









