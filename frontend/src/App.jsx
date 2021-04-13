import React, { useState, useEffect } from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
} from 'react-router-dom';
import axios from 'axios';

import Home from './components/Home/Home';
import Login from './components/Login/Login';
import Register from './components/Register/Register';
import AppContext from './AppContext';

import 'bootstrap/dist/css/bootstrap.min.css';
import './app.sass';
import verifyToken from './misc/verifyToken';

const App = () => {
  const [, setErrors] = useState([]);
  const [currentUser, setCurrentUser] = useState({
    username: '',
    token: '',
    teamId: null,
    isAdmin: false,
    isAuthenticated: false,
  });
  const [boards, setBoards] = useState([{ id: null, name: '' }]);
  const [activeBoardId, setActiveBoardId] = useState(null);

  useEffect(() => {
    verifyToken().then((user) => {
      setCurrentUser(user);
      axios.get(
        `${process.env.REACT_APP_BACKEND_URL}/boards/?team_id=${user.teamId}`,
        { headers: { 'auth-user': user.username, 'auth-token': user.token } },
      ).then((res) => {
        setBoards(
          res.data.boards.map((board) => ({
            id: board.id,
            name: board.name,
          })),
        );
        setActiveBoardId(res.data.boards[0].id);
      }).catch((err) => (
        // TODO: handle
        console.log(`LIST BOARDS ERROR: ${err.data}`)
      ));
    }).catch((err) => (setErrors(err)));
  }, []);

  return (
    <AppContext.Provider
      value={{
        currentUser,
        setCurrentUser,
        boards,
        setBoards,
        activeBoardId,
        setActiveBoardId,
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
