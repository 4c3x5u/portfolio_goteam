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
  getBoard,
} from './misc/api';

import 'bootstrap/dist/css/bootstrap.min.css';
import './app.sass';

const App = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [currentUser, setCurrentUser] = useState({
    username: '',
    token: '',
    teamId: null,
    isAdmin: false,
    isAuthenticated: false,
  });
  const [boards, setBoards] = useState([{ id: null, name: '' }]);
  const [activeBoard, setActiveBoard] = useState(activeBoardInit);
  const loadBoard = async (boardId) => {
    setIsLoading(true);
    try {
      const user = await verifyToken();
      setCurrentUser(user);

      // 1. Get boards, and set boards and active board ID
      const boardsRes = await getBoards(user.teamId);
      setBoards(boardsRes.data.boards);

      const activeBoardId = (
        boardId || activeBoard.id || boardsRes.data.boards[0].id
      );

      const boardRes = await getBoard(activeBoardId);

      console.log(`COLUMNS: ${JSON.stringify(boardRes.data.columns)}`);

      setActiveBoard({
        id: boardRes.data.id,
        columns: boardRes.data.columns,
      });

      // // 2. Get the columns
      // const columnsRes = await getColumns(activeBoardId);
      // const { columns } = columnsRes.data;
      //
      // // 3. Get the tasks per column
      // const rawTasksPerColumn = await Promise.all(
      //   columns.map(async (column) => {
      //     const tasksRes = await getTasks(column.id);
      //     return tasksRes.data.tasks;
      //   }),
      // );
      //
      // // 4. Get the subtasks and feed them into tasks
      // const tasksPerColumn = await Promise.all(
      //   rawTasksPerColumn.map(async (tasks) => (
      //     Promise.all(tasks.map(async (task) => {
      //       const subtasksRes = await getSubtasks(task.id);
      //       const { subtasks } = subtasksRes.data;
      //       return { ...task, subtasks };
      //     }))
      //   )),
      // );
      //
      // // 5. Set the active board
      // setActiveBoard({
      //   id: activeBoardId,
      //   columns: columns
      //     .sort((c1, c2) => c1.order - c2.order)
      //     .map((column) => (
      //       {
      //         id: column.id,
      //         order: columnOrder.parseInt(column.order),
      //         tasks: tasksPerColumn[column.order]
      //           .sort((t1, t2) => t1.order - t2.order),
      //       }
      //     )),
      // });
    } catch (err) {
      console.error(err);
    }
    setIsLoading(false);
  };

  useEffect(() => loadBoard(), []);

  return (
    <AppContext.Provider
      value={{
        currentUser,
        boards,
        activeBoard,
        loadBoard,
        isLoading,
        setIsLoading,
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
