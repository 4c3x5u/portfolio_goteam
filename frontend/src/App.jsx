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
import UserContext from './UserContext';

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

  useEffect(() => { verifyToken(setCurrentUser); }, []);

  return (
    <UserContext.Provider value={{ currentUser, setCurrentUser }}>
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
