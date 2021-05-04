# GoTeam! – Kanban Boards
![GoTeam! App UI Screenshot](https://i.ibb.co/nCty58P/Screenshot-2021-04-29-at-19-20-34.png)

## The Purpose
The purpose of my having developed this application is to demonstrate **my competency in creating Python Web APIs using Django-REST**,as well as **the exploration of some React libraries** that I haven't previously used, such as *react-beautiful-dnd* for drag and drop controls, *toastify* for displaying server-side errors to the user, and several others.

The result is a Django-REST/React.js app that simulates a kanban boards experience for small teams by creating a team for the user upon registration, and then allowing them to invite other users to their team by sharing an invite link with them, which the other users click on and register. After this process, the original user (a.k.a. admin) controls what boards the invited users (a.k.a. members) can access, and what task they move across the board and mark as done therein.

> Please note that I am aware I have gone about things the less idiomatic way by using function-based views. This is due to time restrictions imposed by not being actively employed and to custom logic being required in almost every piece that makes up the app. The chances are, by the time you're reading this, I have already started working on refactoring the app to move over to class-based views by migrating the custom logic into new serializers that inherit the pre-existing ones.

## Running the App
### Requirements
1. [Python 3.5+](https://www.python.org/downloads/release/python-390/)
2. [PIP](https://pypi.org/project/pip/)
4. [Node.js](https://nodejs.org/en/)
5. [Yarn](https://yarnpkg.com/getting-started/install) or [NPM](https://www.npmjs.com/get-npm)
6. [PostgreSQL 13](https://www.postgresql.org/) – It must be running prior to following the *backend* instructions below.
7. In order to run either the frontend or the backend, you must provide some environment variables. *dotenv* is installed on both projects, so you can just create a *.env* file in both the *frontend* and the backend folders, and declare the environment variables within them. For a list of required environment variables, please see the *.env-sample* file in each folder.

### The Backend
> I prefer using [pipenv](https://pypi.org/project/pipenv/) for virtual environments and a *Pipfile* for dependencies. However, I have included a *requirements.txt* file inside the *backend* folder, so you can just use the instructions below to simply and quickly run the app.

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
    
That's it! If all your environment variables checks out and you have followed the instructions, the frontend app must now be running at *http://localhost:3000*, and the backend app at *http://localhost:8000*.

## Running the Tests
There are 200+ tests inside the backend project, which you can run by executing `python manage.py test main.tests` from the *backend* directory while the virtual environment is active.
