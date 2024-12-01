# Enterprise Resource Planning (ERP) System

<!-- <p align="center">
  <img src="path/to/your/erp/image.png" alt="ERP System Screenshot" width="600"> 
</p> -->


## Team Name: NSU_ERP_Team

North South University, Department of Computer Science and Engineering, Software Engineering (CSE-327)


## Team Members:

1. Tanvir Ahmed Khan - 2131491642
2. Sharjil Khan - 2131861642
3. Tahsinul Haque Wrudra - 2131252642
4. Mahir Shahriar Tamim - 2131377642


## How to Use

This section details how to run and interact with the ERP system.  Remember to replace placeholders with your actual project details.

**Prerequisites:**

* **Go:** Go 1.18 or later.  [Download and Installation Instructions](https://go.dev/dl/)
* **PostgreSQL:** PostgreSQL database server. [Download and Installation Instructions](https://www.postgresql.org/download/)
* **Git:** Git version control system. [Download and Installation Instructions](https://git-scm.com/downloads)


**Database Setup:**

1. **Create the Database:**  Use the `psql` command-line tool to create a database named `erp`.  You might need to create a user with appropriate permissions first.  Example (replace with your username and password):
   ```bash
   createdb -U postgres -p 5432 erp
   ```
2. **Run Migrations:**  Execute the database migration scripts to set up the necessary tables and schema.  The location of these scripts will depend on your project structure.  Example (replace with the actual path):
   ```bash
   psql -U postgres -d erp -f ./migrations/init.sql 
   ```


## 1. Clone the Project

- Open your terminal and navigate to your desired project directory.
- Run the command: `git clone <your_git_repo_url>`  (Replace `<your_git_repo_url>` with your actual repository URL)

## 2. Navigate to the Project Directory

- `cd <your_project_directory>` (Replace `<your_project_directory>` with your project's directory name)

## 3. Set Environment Variables

- Create a `.env` file in your project's root directory and add the following (replace with your actual credentials):

```
DB_USER=postgres
DB_PASSWORD=<your_db_password>
DB_NAME=erp
DB_HOST=localhost
DB_PORT=5432
```

- **Important:**  Never commit your `.env` file to your Git repository!  It should be added to your `.gitignore`.


## 4. Run the Project

-  Run the following command to start the ERP system: `go run main.go` (or adjust this command to match your project's startup command).


## 5. Access the System

- Once the application has started, access the ERP system via your web browser at `http://localhost:8080` (or the appropriate address/port specified by your application).



## How to Develop

This section outlines the development workflow.

## 1. Create a New Branch

- `git checkout -b <your_feature_branch_name>`

## 2. Make Changes

- Develop your features.  Adhere to the coding style guide (if one exists).

## 3. Run Tests

- Before committing, run the test suite:  `make test ./...`


## 4. Commit and Push Changes

- `git add .`
- `git commit -m "Your descriptive commit message"`
- `git push origin <your_feature_branch_name>`

## 5. Create a Pull Request

- Create a pull request on GitHub

## 6. Review and Merge

- After code review and approval, merge the pull request.

## 7. Update Local Repository

- `git checkout main`
- `git pull origin main`


Remember to replace all placeholders (like `<your_git_repo_url>`, `<your_project_directory>`, database credentials, and ports).