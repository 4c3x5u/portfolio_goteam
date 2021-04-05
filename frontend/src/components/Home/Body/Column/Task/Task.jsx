import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSquare, faCheckSquare } from '@fortawesome/free-regular-svg-icons';

import './task.sass';

const Task = () => (
  <div className="Task">
    <h1 className="Title">Do Something</h1>

    <p className="Description">
      Lorem ipsum dolor sit amet falan filan iste kanka daha ne olsun yani
      anlatabiliyo muyum? Peki, o zaman.
    </p>

    <ul className="Subtasks">
      <li className="Subtask">
        <button
          className="CheckButton"
          onClick={() => console.log('check/uncheck')}
          type="button"
        >
          <FontAwesomeIcon className="CheckBox" icon={faSquare} />
          Undone Subtask
        </button>
      </li>

      <li className="Subtask">
        <button
          className="CheckButton"
          onClick={() => console.log('check/uncheck')}
          type="button"
        >
          <FontAwesomeIcon className="CheckBox" icon={faCheckSquare} />
          Done task
        </button>
      </li>
    </ul>
  </div>
);

export default Task;
