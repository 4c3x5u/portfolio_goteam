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
} from './misc/api';

import 'bootstrap/dist/css/bootstrap.min.css';
import './app.sass';
import columnOrder from './components/Home/Board/Column/columnOrder';

const App = () => {
  const [currentUser, setCurrentUser] = useState({
    username: '',
    token: '',
    teamId: null,
    isAdmin: false,
    isAuthenticated: false,
  });
  const [activeBoard, setActiveBoard] = useState(activeBoardInit);
  const loadActiveBoard = async () => {
    try {
      const user = await verifyToken();
      setCurrentUser(user);

      // 1. Get the active board ID
      const boardsRes = await getBoards(user.teamId);
      const activeBoardId = boardsRes.data.boards[0].id;

      // 2. Get the columns
      const columnsRes = await getColumns(activeBoardId);
      const columns = columnsRes.data.columns.sort((column) => column.order);

      // 3. Get the tasks per column
      const rawTasksPerColumn = await Promise.all(
        columns.map(async (column) => {
          const tasksRes = await getTasks(column.id);
          return tasksRes.data.tasks.sort((task) => task.order);
        }),
      );

      // 4. Get the subtasks and feed them into tasks
      const tasksPerColumn = await Promise.all(
        rawTasksPerColumn.map(async (tasks) => (
          Promise.all(tasks.map(async (task) => {
            const subtasksRes = await getSubtasks(task.id);
            const { subtasks } = subtasksRes.data;
            return { ...task, subtasks };
          }))
        )),
      );

      // 5. Set the active board
      setActiveBoard({
        id: activeBoardId,
        columns: columns.sort((column) => column.order).map((column, i) => (
          {
            id: column.id,
            order: columnOrder.parseInt(column.order),
            tasks: tasksPerColumn[i],
          }
        )),
      });
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => loadActiveBoard(), []);

  return (
    <AppContext.Provider
      value={{
        currentUser,
        setCurrentUser,
        activeBoard,
        setActiveBoard,
        loadActiveBoard,
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
