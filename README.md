# Gym API

The Gym API is a project designed to help users manage and track their workouts. With this API, users can create, update, and view their exercise routines efficiently.

## Setup Instructions

Follow these steps to set up the project on your local machine:

1. **Clone the Project:**
   ```bash
   git clone git@github.com:davidspader/gym-api.git
   cd gym-api
   ```

2. **Configure Environment Variables:**
   - Copy the example environment file:
     ```bash
     cp .env-example .env
     ```
   - Edit the `.env` file to configure the environment variables:
     ```dotenv
     POSTGRES_USER=your_user
     POSTGRES_PASSWORD=your_password
     POSTGRES_DB=your_database
     DB_URL=postgresql://your_user:your_password@localhost/your_database
     API_PORT=your_preferred_port
     SECRET_KEY=your_base64_encoded_string
     ```

3. **Install Go Modules:**
   ```bash
   go mod tidy
   ```

4. **Start Docker Containers:**
   ```bash
   docker compose up -d
   ```

5. **Run the SQL Script:**
   - Import the database schema from the SQL script located at `sql/script.sql` into your PostgreSQL database.

6. **Generate Test Data (Optional):**
   - If you want to generate test data for the API, run the SQL script located at `sql/data.sql`. Note that the default password for the test users is `123456`.

Your system is now configured and ready to use!