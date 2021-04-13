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
import UserContext from './UserContext';

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

  useEffect(() => {
    verifyToken().then((user) => {
      setCurrentUser(user);
      axios.get(
        `${process.env.REACT_APP_BACKEND_URL}/boards/?team_id=${user.teamId}`,
        { headers: { 'auth-user': user.username, 'auth-token': user.token } },
      ).then((res) => setBoards(
        res.data.boards.map((board) => ({
          id: board.id,
          name: board.name,
        })),
      )).catch((err) => (
        console.log(`LIST BOARDS ERROR: ${err.data}`)
      ));
    }).catch((err) => (setErrors(err)));

    console.log(`BOARDS: ${boards}`);
  }, []);

  return (
    <UserContext.Provider
      value={{
        currentUser,
        setCurrentUser,
        boards,
        setBoards,
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
    </UserContext.Provider>
  );
};

export default App;
