import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import Register from './components/auth/Register';
import Login from './components/auth/Login';
import 'bootstrap/dist/css/bootstrap.min.css';
import './app.sass';

const App = () => (
  <Router className="App">
    <Switch>
      <Route exact path="/login" component={Login} />
      <Route path="/" component={Register} />
    </Switch>
  </Router>
);

export default App;
