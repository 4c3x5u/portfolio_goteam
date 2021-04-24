/* eslint-disable react/jsx-props-no-spreading */
import React, { useContext, useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Droppable } from 'react-beautiful-dnd';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/free-solid-svg-icons';
import _ from 'lodash/core';

import Task from './Task/Task';
import columnOrder from './columnOrder';
import { capFirstLetterOf } from '../../../../misc/util';
import AppContext from '../../../../AppContext';

import './column.sass';
import window from '../../../../misc/window';

const Column = ({
  id, order, tasks, handleActivate,
}) => {
  const { user } = useContext(AppContext);
  const [name, setName] = useState('');

  useEffect(() => (
    order !== null && setName(columnOrder.parseInt(order))
  ), [order]);

  return (
    <Col className="ColumnWrapper" xs={3}>
      <div className={`Column ${name && capFirstLetterOf(name)}Column`}>
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
              {_.sortBy(tasks, (task) => task.order)
                .map((task) => (
                  <Task
                    id={task.id}
                    title={task.title}
                    description={task.description}
                    order={task.order}
                    assignedUser={task.user}
                    handleActivate={handleActivate}
                    subtasks={task.subtasks}
                  />
                ))}

              {provided.placeholder}

              {user.isAdmin && name === columnOrder.INBOX && (
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
    assignedUser: PropTypes.string.isRequired,
  }).isRequired,
  handleActivate: PropTypes.func,
};

Column.defaultProps = {
  handleActivate: () => console.log('Cannot create task here.'),
};

export default Column;
