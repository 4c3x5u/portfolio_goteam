import React from 'react';

import Subtask from './Subtask/Subtask';

import './task.sass';

const Task = () => (
  <div className="Task">
    <h1 className="Title">Do Something</h1>

    <p className="Description">
      Lorem ipsum dolor sit amet falan filan iste kanka daha ne olsun yani
      anlatabiliyo muyum? Peki, o zaman.
    </p>

    <ul className="Subtasks">
      <Subtask title="Subtask #1" />
      <Subtask title="Subtask #2" />
    </ul>
  </div>
);

export default Task;
