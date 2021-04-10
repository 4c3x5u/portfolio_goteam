import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from 'react-router-dom';

import Auth from './components/Auth/Auth';
import Home from './components/Home/Home';

import 'bootstrap/dist/css/bootstrap.min.css';
import './app.sass';

const App = () => (
  <Router className="App">
    <Switch>
      <Route exact path="/" component={Home} />
      <Route path="/" component={Auth} />
    </Switch>
  </Router>
);

export default App;
