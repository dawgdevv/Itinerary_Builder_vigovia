# How to Run the Vigovia Itinerary API

This guide will walk you through the steps to set up and run the Vigovia Itinerary API on your local machine.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.18 or higher** - [Download Go](https://golang.org/dl/)
- **Postman** or **Insomnia** - API testing client
  - [Download Postman](https://www.postman.com/downloads/)
  - [Download Insomnia](https://insomnia.rest/download)

## Step 1: Navigate to the Project Directory

Open your terminal and navigate to the project folder:

```bash
cd vigovia-task
```

## Step 2: Install Dependencies

The project uses Go modules for dependency management. Install all required dependencies:

```bash
go mod download
```

This will download:

- Gin web framework
- gofpdf for PDF generation
- bcrypt for password hashing
- And other required packages

## Step 3: Run the Application

You have two options to run the application:

### Option A: Run Directly (Recommended for Development)

This method compiles and runs the application in one step:

```bash
go run main.go
```

You should see the following output:

```
Starting Itinerary Builder API on http://localhost:8080
```

The server is now running and ready to accept requests!

### Option B: Build and Execute (Recommended for Production)

This method creates a standalone executable binary:

1. **Build the application:**

   ```bash
   go build -o vigovia-task
   ```

   This creates an executable file named `vigovia-task` in the current directory.

2. **Run the executable:**

   ```bash
   ./vigovia-task
   ```

   You should see:

   ```
   Starting Itinerary Builder API on http://localhost:8080
   ```

**Benefits of this approach:**

- Faster startup (no compilation needed)
- Can run the binary on any machine without Go installed
- Better for deployment scenarios

## Step 4: Import the Postman/Insomnia Collection

To easily test all API endpoints, import the provided collection:

### For Postman:

1. **Open Postman**
2. Click **"Import"** button (top-left corner)
3. Choose **"Upload Files"** or drag and drop
4. Select the file: `Vigovia_API_Postman_Collection.json`
5. Click **"Import"**

### For Insomnia:

1. **Open Insomnia**
2. Click **"Create"** → **"Import From"** → **"File"**
3. Select the file: `Vigovia_API_Postman_Collection.json`
4. Click **"Scan"** and then **"Import"**

## Step 5: Configure Environment (Important!)

### For Postman:

1. Click the **"Environments"** dropdown (top-right, next to the eye icon)
2. Select **"Create new environment"** or click the **"+"** icon
3. Name it (e.g., "Vigovia Local")
4. Add the following variable (optional but recommended):
   - Variable: `baseUrl`
   - Initial Value: `http://localhost:8080`
   - Current Value: `http://localhost:8080`
5. Click **"Save"**
6. **Select this environment** from the dropdown

**Why this is important:** The collection uses cookies for authentication. Having a proper environment ensures cookies are stored and sent correctly with each request.

### For Insomnia:

1. Click on **"No Environment"** dropdown at the top
2. Select **"Manage Environments"**
3. Create a new environment or use the default
4. Add:
   ```json
   {
     "baseUrl": "http://localhost:8080"
   }
   ```
5. Click **"Done"**

## Step 6: Test the API

Once the server is running and the collection is imported:

1. **Test Health Check:**

   - Request: `GET /`
   - Should return: `{"message": "Itinerary Builder API is running!"}`

2. **Sign Up a User:**

   - Request: `POST /api/auth/signup`
   - Creates a new user and returns an authentication token

3. **Create an Itinerary:**

   - Request: `POST /api/itineraries`
   - The auth token is automatically sent via cookies

4. **Export to PDF:**
   - Request: `GET /api/itineraries/:id/export-pdf`
   - PDF file will be saved in the `/output` folder

## Step 7: Watch the Loom Video Tutorial

For a complete visual walkthrough of:

- Starting the server
- Testing endpoints in Postman/Insomnia
- Creating itineraries
- Generating PDFs
- Troubleshooting common issues

**Watch the Loom video included with this project.**

The video demonstrates:
✅ How to run the application effortlessly  
✅ How to use the Postman collection  
✅ How to use the Insomnia collection  
✅ Complete API workflow from signup to PDF export  
✅ Common pitfalls and how to avoid them

## Troubleshooting

### Server won't start:

- **Port 8080 already in use:** Kill the process using port 8080 or change the port in `main.go`
  ```bash
  # Find process on port 8080
  lsof -i :8080
  # Kill it (replace PID with actual process ID)
  kill -9 <PID>
  ```

### Dependencies not found:

- Run `go mod tidy` to clean up and download dependencies
- Ensure you're in the `vigovia-task` directory

### Authentication not working:

- Make sure you've created and selected an environment in Postman
- Check that cookies are enabled in your API client
- Look for the `auth_token` cookie after login/signup

### PDF not generating:

- Check that the `/output` folder exists in the project directory
- Verify file permissions allow writing to the folder

## Additional Information

- **API Documentation:** See `API_DOCUMENTATION.md` for complete endpoint details
- **Sample Data:** Check the `/examples` folder for request/response samples
- **Port:** Default is `8080` - can be changed in `main.go`
- **Storage:** Currently uses in-memory storage (data is lost on restart)

## Stopping the Server

To stop the running server:

- Press `Ctrl + C` in the terminal where the server is running

---

**Need Help?** Watch the Loom video tutorial or refer to `API_DOCUMENTATION.md` for detailed API usage.
