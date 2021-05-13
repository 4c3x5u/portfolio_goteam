# GoTeam! – Kanban Boards
![GoTeam! App UI Screenshot](https://i.ibb.co/nCty58P/Screenshot-2021-04-29-at-19-20-34.png)

## The Purpose
The purpose of this application is to demonstrate **my skills and knowledge in creating Python web APIs using Django-REST**, as well as **the exploration of some React libraries that I haven't previously used**, such as *react-beautiful-dnd* for drag and drop controls, *toastify* for displaying server-side errors to the user, and several others.

The result is a Django-REST/React.js app that simulates a kanban boards experience for small teams by creating a team for the user upon registration, and then allowing them to invite other users to their team by sharing an invite link with them, which the other users click on and register. After this process, the original user (a.k.a. admin) controls what boards the invited users (a.k.a. members) can access, and what task they can move across the board and mark its subtasks done therein.

## The State of the App
GoTeam! is now largely functional and is written in idiomatic Python and JavaScript code for the most part. Though there are some issues that I might consider taking on soon — visit the [issues](https://github.com/alicandev/portfolio_goteam/issues) section to review them.

## Running the App
### Requirements
1. [Python 3.5+](https://www.python.org/downloads/release/python-390/)
2. [PIP](https://pypi.org/project/pip/)
4. [Node.js](https://nodejs.org/en/)
5. [Yarn](https://yarnpkg.com/getting-started/install) or [NPM](https://www.npmjs.com/get-npm)
6. [PostgreSQL 13](https://www.postgresql.org/) – It must be running prior to following the *backend* instructions below.
7. In order to run either the frontend or the backend, you must provide some environment variables. *dotenv* is installed on both projects, so you can just create a *.env* file in both the *frontend* and the backend folders, and declare the environment variables within them. For a list of required environment variables, please see the *.env-sample* file in each folder.

### The Backend
> I prefer using [Poetry](https://python-poetry.org/) for managing virtual environments and dependencies. However, I included a *requirements.txt* file inside the *backend* folder, so you can just use the instructions below to simply and quickly run the app.

1. Inside a terminal, navigate into the backend folder. 
    - `cd backend`
2. Create a new virtual environment. 
    - `python3 -m venv env` on Mac, or `python -m venv env` on Windows.
3. Activate the virtual environment.
    - `source env/bin/activate`
4. Install dependencies from the requirements file.
    - `pip install -r requirements.txt`
5. Run the app.
    - `python manage.py runserver`
    
### The Frontend
1. Keep the backend running, and create a new tab inside your terminal.
    - Usually **[CMD + T]** on Mac, or **[CTRL + T]** on Windows.
2. Navigate to the *frontend* folder.
    - `cd frontend` from the root, or `cd ../frontend` from the *backend* folder
3. Install dependencies.
    - `yarn install` or `npm install`
4. Run the app.
    - `yarn start` or `npm start`
    
That's it! If all your environment variables check out and you have followed the instructions, the frontend app must now be running at *http://localhost:3000*, and the backend app at *http://localhost:8000*.

## Running the Tests
There are 200+ tests inside the backend project, which you can run by executing `python manage.py test main.tests` from the *backend* directory while the virtual environment is active.
