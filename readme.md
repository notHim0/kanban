Kanban API
This project is a simple Kanban API built with Go.
Getting Started
Follow these steps to set up and run the project locally.

1. Clone the Repository
   Clone the project from GitHub and navigate into the directory.
   git clone [https://github.com/notHim0/kanban](https://github.com/notHim0/kanban)
   cd kanban

2. Install Dependencies
   Install all the required Go packages.
   go get [github.com/gorilla/mux](https://github.com/gorilla/mux) \
   [github.com/justinas/alice](https://github.com/justinas/alice) \
   [github.com/joho/godotenv](https://github.com/joho/godotenv) \
   [github.com/lib/pq](https://github.com/lib/pq) \
   golang.org/x/crypto \
   [github.com/xeipuuv/gojsonpointer](https://github.com/xeipuuv/gojsonpointer) \
   [github.com/xeipuuv/gojsonreference](https://github.com/xeipuuv/gojsonreference) \
   [github.com/xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema)

3. Setup Environment Variables
   Create a file named .env in the root directory and add your database URL and a JWT secret key.
   DATABASE_URL=your_postgres_database_url
   JWT_SECRET=your_jwt_secret_key

4. Run the Application
   Start the application by running the main Go file.
   go run cmd/kanban/main.go

5. API Endpoints
   The API provides the following core functionalities:
   User Management: Registration and login.
   Project Management: Creating, updating, and deleting projects.
   Data Retrieval: Fetching specific project data.
