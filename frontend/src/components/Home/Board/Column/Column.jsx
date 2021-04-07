import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/free-solid-svg-icons';

import Task from './Task/Task';
import { columnOrder } from './columnOrder';
import { capFirstLetterOf } from '../../../../misc/utils';

import './column.sass';

const Column = ({ id, name, toggleCreateTask }) => {
  const [tasks, setTasks] = useState([]);

  // TODO: API call based on column id here
  console.log(`Column ID: ${id}`);
  useEffect(() => (
    setTasks([
      {
        id: 0,
        title: 'Some task',
        description: 'Some task description',
        order: 0,
      },
    ])
  ), []);

  return (
    <Col className="Col" xs={3}>
      <div className={`Column ${capFirstLetterOf(name)}Column`}>
        <div className="Header">{name.toUpperCase()}</div>

        <div className="Body">
          {
            tasks
              .sort((task) => task.order)
              .map((task) => (
                <Task
                  id={task.id}
                  title={task.title}
                  description={task.description}
                />
              ))
          }

          {name === columnOrder.INBOX && (
            <button
              className="CreateButton"
              onClick={toggleCreateTask}
              type="button"
            >
              <FontAwesomeIcon className="Icon" icon={faPlusCircle} />
            </button>
          )}
        </div>
      </div>
    </Col>
  );
};

Column.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  toggleCreateTask: PropTypes.func,
};

Column.defaultProps = {
  toggleCreateTask: () => console.log('Cannot create task here.'),
};

export default Column;
