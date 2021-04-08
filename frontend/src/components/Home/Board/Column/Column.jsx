import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/free-solid-svg-icons';

import Task from './Task/Task';
import { columnOrder } from './columnOrder';
import { capFirstLetterOf } from '../../../../misc/utils';

import './column.sass';
import window from '../../../../misc/window';

const Column = ({ id, name, handleActivate }) => {
  const [tasks, setTasks] = useState([]);

  // TODO: API call based on column id here
  console.log(`Column ID: ${id}`);
  useEffect(() => (
    name === 'inbox' && setTasks([
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
                  handleActivate={handleActivate}
                />
              ))
          }

          {name === columnOrder.INBOX && (
            <button
              className="CreateButton"
              onClick={handleActivate(window.CREATE_TASK)}
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
  handleActivate: PropTypes.func,
};

Column.defaultProps = {
  handleActivate: () => console.log('Cannot create task here.'),
};

export default Column;
