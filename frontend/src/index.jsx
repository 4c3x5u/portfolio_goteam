import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';

document.addEventListener('contextmenu', (e) => e.preventDefault());

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById('root'),
);
