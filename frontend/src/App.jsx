import React, { useState } from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from 'react-router-dom';

import Home from './components/Home/Home';
import Login from './components/Login/Login';
import Register from './components/Register/Register';

import 'bootstrap/dist/css/bootstrap.min.css';
import './app.sass';

const App = () => {
  const [currentUser, setCurrentUser] = useState({
    username: '',
    teamId: '',
    isAdmin: '',
  });

  return (
    <Router className="App">
      <Switch>
        <Route exact path="/">
          <Home currentUser={currentUser} />
        </Route>
        <Route exact path="/login">
          <Login setCurrentUser={setCurrentUser} />
        </Route>
        <Route exact path="/register">
          <Register setCurrentUser={setCurrentUser} />
        </Route>
      </Switch>
    </Router>
  );
};

export default App;
