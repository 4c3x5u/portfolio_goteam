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

import 'bootstrap/dist/css/bootstrap.min.css';
import './app.sass';
import verifyToken from './misc/verifyToken';

const App = () => {
  const [currentUser, setCurrentUser] = useState({
    username: '',
    teamId: null,
    isAdmin: false,
    isAuthenticated: false,
  });

  useEffect(() => {
    verifyToken().then((res) => setCurrentUser({
      username: res.data.username,
      teamId: res.data.teamId,
      isAdmin: res.data.isAdmin,
      isAuthenticated: true,
    })).catch(() => setCurrentUser({
      username: '',
      teamId: null,
      isAdmin: false,
      isAuthenticated: false,
    }));
  }, []);

  return (
    <Router className="App">
      <Switch>
        <Route exact path="/">
          {currentUser.isAuthenticated
            ? <Home currentUser={currentUser} />
            : <Redirect to="/login" />}
        </Route>

        <Route exact path="/login">
          {!currentUser.isAuthenticated
            ? <Login setCurrentUser={setCurrentUser} />
            : <Redirect to="/" />}
        </Route>

        <Route exact path="/register">
          {!currentUser.isAuthenticated
            ? <Register setCurrentUser={setCurrentUser} />
            : <Redirect to="/" />}
        </Route>
      </Switch>
    </Router>
  );
};

export default App;
