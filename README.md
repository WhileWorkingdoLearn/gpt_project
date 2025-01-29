Prompt to ChatGpt:

Hello, I would like to start a software development project as an exercise. Since I am currently without ideas and don't know what I could do, I need your help. I would like to develop a small to medium sized software project. The project should be as close as possible to real problems / tasks in software development. Since I want to focus on the bakcend, frontend elements can occur in the context of the fullstack, but should play a subordinate role. 
The content of the project should include databases, api calls, possibly multithreading or tasks, filewriting reading ect. 
The requirement should be that of an intermediate developer. Please play the customer. Sectors could be industry or the financial sector. As programming language I would like to use Golang. If possible, the requirements for the project should be backed up with a few test cases against which I can then develop. If you have any further questions, please feel free to ask.


Response from Gpt:

Here is a project idea for you:

Project: “Order management system for the industry”
Description:

You are developing a backend system that can be used by industrial companies to manage production orders. The system offers the following functions:

    Order creation and management:
        CRUD operations for orders.
        Storage in a relational database (e.g. PostgreSQL).

    Status tracking:
        Each order has statuses (e.g. “New”, “In progress”, “Completed”).
        Automatic status changes based on defined conditions.

    API integration:
        RESTful API for communication with external systems (e.g. ERP systems).
        Retrieval of production resources or data via external APIs.

    File processing:
        Import of order data from CSV files.
        Export of reports (e.g. monthly reports) as PDF or Excel.



    Task scheduling:
        Regular checking of jobs by background tasks (e.g. via go routines).
        Notifications/logs for critical events.

    Multithreading:
        Optimization of tasks such as status updates or data processing through goroutines.

Requirements (as “customer requests”):

    API endpoints:
        POST /orders: Create a new order.
        GET /orders: List all orders.
        PATCH /orders/{id}: Update the status of an order.
        DELETE /orders/{id}: Delete an order.

    Database model:
        Orders table with columns such as id, name, status, created_at, updated_at.

    CSV import:
        An endpoint /import is to read CSV files and create orders.

    Regular task:
    Every hour, expired tasks should be automatically marked as “Completed”.
Would you like detailed instructions on how to start the project? 

