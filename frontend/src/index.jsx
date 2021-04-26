import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';

document.addEventListener('contextmenu', (e) => e.preventDefault());

ReactDOM.render(
  <App />,
  document.getElementById('root'),
);
