/* eslint-disable react/jsx-props-no-spreading */
import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Droppable } from 'react-beautiful-dnd';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/free-solid-svg-icons';

import Task from './Task/Task';
import columnOrder from './columnOrder';
import { capFirstLetterOf } from '../../../../misc/util';

import './column.sass';
import window from '../../../../misc/window';

const Column = ({
  id, order, tasks, handleActivate,
}) => {
  const [name, setName] = useState('');

  useEffect(() => (
    order !== null && setName(columnOrder.parseInt(order))
  ), [order]);

  return (
    <Col className="Col" xs={3}>
      {console.log(`NAME: ${name}`)}
      <div
        className={
          `Column ${name && name !== '' && capFirstLetterOf(name)}Column`
        }
      >
        <div className="Header">
          {name && name.toUpperCase()}
        </div>

        <Droppable droppableId={`${id}`}>
          {(provided) => (
            <div
              className="Body"
              {...provided.droppableProps}
              ref={provided.innerRef}
            >
              {tasks
                .sort((task) => task.order)
                .map((task) => (
                  task && (
                    <Task
                      id={task.id}
                      title={task.title}
                      description={task.description}
                      order={task.order}
                      handleActivate={handleActivate}
                      subtasks={task.subtasks}
                    />
                  )
                ))}

              {provided.placeholder}

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
          )}
        </Droppable>
      </div>
    </Col>
  );
};

Column.propTypes = {
  id: PropTypes.string.isRequired,
  order: PropTypes.number.isRequired,
  tasks: PropTypes.arrayOf({
    id: PropTypes.number.isRequired,
    title: PropTypes.string.isRequired,
    description: PropTypes.string.isRequired,
    order: PropTypes.number.isRequired,
  }).isRequired,
  handleActivate: PropTypes.func,
};

Column.defaultProps = {
  handleActivate: () => console.log('Cannot create task here.'),
};

export default Column;
