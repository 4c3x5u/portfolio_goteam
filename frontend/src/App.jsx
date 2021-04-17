import React, { useState, useEffect } from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
} from 'react-router-dom';

import AppContext from './AppContext';
import UserAPI from './api/UserAPI';
import BoardsAPI from './api/BoardsAPI';
import Home from './components/Home/Home';
import Login from './components/Login/Login';
import Register from './components/Register/Register';
import activeBoardInit from './misc/activeBoardInit';

import 'bootstrap/dist/css/bootstrap.min.css';
import './app.sass';

const App = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [user, setUser] = useState({
    username: '',
    teamId: null,
    isAdmin: false,
    isAuthenticated: false,
  });
  const [boards, setBoards] = useState([{ id: null, name: '' }]);
  const [activeBoard, setActiveBoard] = useState(activeBoardInit);

  const loadBoard = async (boardId) => {
    setIsLoading(true);
    try {
      const userResponse = await UserAPI.verifyToken();
      delete userResponse.data.msg;
      setUser({ ...userResponse.data, isAuthenticated: true });

      const teamBoards = await BoardsAPI.get(null, userResponse.data.teamId);
      setBoards(teamBoards.data);

      const nestedBoard = await BoardsAPI.get((
        teamBoards.data.length === 1 && teamBoards.data[0].id
      ) || boardId || activeBoard.id);

      setActiveBoard(nestedBoard.data);
    } catch (err) {
      setUser({ ...user, isAuthenticated: false });
      sessionStorage.removeItem('username');
      sessionStorage.removeItem('auth-token');
      console.error(err);
    }
    setIsLoading(false);
  };

  useEffect(() => loadBoard(), []);

  return (
    <AppContext.Provider
      value={{
        user,
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
            {user.isAuthenticated
              ? <Home />
              : <Redirect to="/login" />}
          </Route>

          <Route exact path="/login">
            {!user.isAuthenticated
              ? <Login />
              : <Redirect to="/" />}
          </Route>

          <Route exact path="/register">
            {!user.isAuthenticated
              ? <Register />
              : <Redirect to="/" />}
          </Route>
        </Switch>
      </Router>
    </AppContext.Provider>
  );
};

export default App;
