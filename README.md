# GoTeam! – Kanban Boards
![GoTeam! App UI Screenshot](https://i.ibb.co/nCty58P/Screenshot-2021-04-29-at-19-20-34.png)

## Aim and Purpose
I have developed this application in an effort to demonstrate my competency in creating Python Web APIs using Django-REST.

Please note that I am aware I have gone about things the less idiomatic way by using function-based views. This is due to 
time restrictions imposed by not being actively employed and to custom logic being required in almost every piece that 
makes up the app.

The chances are, by the time you're reading this, I have already started working on refactoring the app to move over to 
class-based views by migrating the custom logic into new serializers that inherit the pre-existing ones.

## Running the App
### Requirements
1. [Python 3.5+](https://www.python.org/downloads/release/python-390/)
2. [PIP](https://pypi.org/project/pip/) (Comes with Python.)
3. [Node](https://nodejs.org/en/), 
4. [NPM](https://www.npmjs.com/get-npm) (Comes with Node.)
5. [PostgreSQL](https://www.postgresql.org/) – It must be running prior to following the *backend* instructions below.
6. In order to run either the frontend or the backend, you must provide some environment variables. *dotenv* is installed
on both projects, so you can just create a *.env* file in both the *frontend* and the backend folders, and declare the environment
variables within them. For a list of required environment variables, please see the *.env-sample* file in each folder.

### The Backend
> I prefer using [pipenv](https://pypi.org/project/pipenv/) for virtual environments and a *Pipfile* for dependencies.
However, I have included a *requirements.txt* file just so you can use the instructions below to quickly run the app.

1. Inside a terminal, navigate into the backend folder. 
    - `cd backend`
2. Create a new virtual environment. 
    - `python -m venv env` on Windows, or `python3 -m venv env` on Mac.
3. Activate the virtual environment.
    - `source env/bin/activate`
4. Install dependencies from the requirements file.
    - `pip install -r requirements.txt`
5. Run the app.
    - `python manage.py runserver`
    
### The Frontend
> I prefer to use [yarn](https://yarnpkg.com) as my JavaScript package manager, but you can use npm for convenience sake.
1. Keep the backend running, and create a new tab inside your terminal.
    - Usually **[CMD + T]** on Mac, or **[CTRL + T]** on Windows.
2. Navigate to the *frontend* folder.
    - `cd frontend` from the root, or `cd ../frontend` from the *backend* folder
3. Install dependencies.
    - `yarn install`, or `npm install`
4. Run the app.
    - `yarn start`, or `npm run start`
    
That's it! If all your environment variables checks out and you have followed the instructions, the frontend app must now be
running at *localhost:3000*, and the backend app at *localhost:8000*.
