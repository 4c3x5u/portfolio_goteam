import React, { useState, useEffect } from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
} from 'react-router-dom';

import AppContext from './AppContext';
import AuthAPI from './api/AuthAPI';
import UsersAPI from './api/UsersAPI';
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
  const [members, setMembers] = useState([{ username: '', isActive: false }]);
  const [boards, setBoards] = useState([{ id: null, name: '' }]);
  const [activeBoard, setActiveBoard] = useState(activeBoardInit);

  const loadBoard = async (boardId) => {
    setIsLoading(true);
    try {
      // 1. set the current user
      const userResponse = await AuthAPI.verifyToken();
      delete userResponse.data.msg;
      setUser({ ...userResponse.data, isAuthenticated: true });

      // 2. set boards lists
      const teamBoards = await BoardsAPI.get(null, userResponse.data.teamId);
      setBoards(teamBoards.data);

      // 3. set the active board
      const nestedBoard = await BoardsAPI.get((
        teamBoards.data.length === 1 && teamBoards.data[0].id
      ) || boardId || activeBoard.id);
      setActiveBoard(nestedBoard.data);

      // 4. set members list
      const teamMembers = await UsersAPI.get(
        userResponse.data.teamId,
        nestedBoard.data.id,
      );
      setMembers(teamMembers.data);
    } catch (err) {
      console.error(err);
    }
    setIsLoading(false);
  };

  useEffect(() => loadBoard(), []);

  return (
    <AppContext.Provider
      value={{
        user,
        members,
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
