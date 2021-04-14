import React, { useState, useEffect } from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
} from 'react-router-dom';

import Home from './components/Home/Home';
import Login from './components/Login/Login';
import Register from './components/Register/Register';
import AppContext from './AppContext';
import activeBoardInit from './misc/activeBoardInit';
import {
  verifyToken,
  getBoards,
  getColumns,
  getTasks,
  getSubtasks,
} from './misc/apiCalls';

import 'bootstrap/dist/css/bootstrap.min.css';
import './app.sass';

const App = () => {
  const [currentUser, setCurrentUser] = useState({
    username: '',
    token: '',
    teamId: null,
    isAdmin: false,
    isAuthenticated: false,
  });
  const [activeBoard, setActiveBoard] = useState(activeBoardInit);

  useEffect(async () => {
    try {
      const user = await verifyToken();
      setCurrentUser(user);

      // 1. Get the active board ID
      const boardsRes = await getBoards(user.teamId);
      let tempActiveBoard = { id: boardsRes.data.boards[0].id };

      // 2. Get the columns
      const columnsRes = await getColumns(tempActiveBoard.id);
      const columns = columnsRes.data.columns.sort((column) => column.order);
      tempActiveBoard = { ...tempActiveBoard, columns };

      // 3. Get the tasks
      const rawTasksPerColumn = await Promise.all(
        tempActiveBoard.columns.map(async (column) => {
          const tasksRes = await getTasks(column.id);
          return tasksRes.data.tasks.sort((task) => task.order);
        }),
      );
      tempActiveBoard = {
        ...tempActiveBoard,
        columns: columns.map((column, i) => (
          { id: column.id, tasks: rawTasksPerColumn[i] }
        )),
      };

      // 4. Get the subtasks and feed them into tasks
      const tasksPerColumn = await Promise.all(
        tempActiveBoard.columns.map(async (column) => (
          Promise.all(column.tasks.map(async (task) => {
            const subtasksRes = await getSubtasks(task.id);
            const { subtasks } = subtasksRes.data;
            return { ...task, subtasks };
          }))
        )),
      );
      tempActiveBoard = {
        ...tempActiveBoard,
        columns: columns.map((column, i) => (
          { id: column.id, tasks: tasksPerColumn[i] }
        )),
      };

      // 5. Set the active board state as the accumulated result
      setActiveBoard(tempActiveBoard);
    } catch (err) {
      console.log(err);
    }
  }, []);

  return (
    <AppContext.Provider
      value={{
        currentUser, setCurrentUser, activeBoard, setActiveBoard,
      }}
    >
      <Router className="App">
        <Switch>
          <Route exact path="/">
            {currentUser.isAuthenticated
              ? <Home />
              : <Redirect to="/login" />}
          </Route>

          <Route exact path="/login">
            {!currentUser.isAuthenticated
              ? <Login />
              : <Redirect to="/" />}
          </Route>

          <Route exact path="/register">
            {!currentUser.isAuthenticated
              ? <Register />
              : <Redirect to="/" />}
          </Route>
        </Switch>
      </Router>
    </AppContext.Provider>
  );
};

export default App;
