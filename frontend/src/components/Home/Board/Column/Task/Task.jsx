import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

import Subtask from './Subtask/Subtask';

import './task.sass';

const Task = ({ id, title, description }) => {
  const [subtasks, setSubtasks] = useState([]);

  // TODO: API call based on task ID
  console.log(`Task ID: ${id}`);
  useEffect(() => (
    setSubtasks([
      {
        id: 0,
        title: 'Some subtask',
        order: 0,
        done: false,
      },
      {
        id: 1,
        title: 'Some subtask',
        order: 1,
        done: false,
      },
    ])
  ), []);

  return (
    <div className="Task">
      <h1 className="Title">
        {title}
      </h1>

      <p className="Description">
        {description}
      </p>

      {subtasks.length > 0 && (
        <ul className="Subtasks">
          {subtasks.sort((subtask) => subtask.order).map((subtask) => (
            <Subtask
              id={subtask.id}
              title={subtask.title}
              done={subtask.done}
            />
          ))}
        </ul>
      )}
    </div>
  );
};

Task.propTypes = {
  id: PropTypes.number.isRequired,
  title: PropTypes.string.isRequired,
  description: PropTypes.string.isRequired,
};

export default Task;
