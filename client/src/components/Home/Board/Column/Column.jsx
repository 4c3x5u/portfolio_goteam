import React, { useContext, useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Droppable } from 'react-beautiful-dnd';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/free-solid-svg-icons';
import _ from 'lodash/core';

import Task from './Task/Task';
import columnOrder from './columnOrder';
import capFirstLetterOf from '../../../../misc/util';
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
      <div
        className={
          `Column ${name && capFirstLetterOf(name.replace('!', ''))}Column`
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
              {_.sortBy(tasks, (task) => task.order).map((task) => (
                <Task
                  id={task.id}
                  key={task.id}
                  title={task.title}
                  description={task.description}
                  order={task.order}
                  assignedUser={task.user}
                  columnId={id}
                  subtasks={task.subtasks}
                  handleActivate={handleActivate}
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
  id: PropTypes.number.isRequired,
  order: PropTypes.number.isRequired,
  tasks: PropTypes.arrayOf(
    PropTypes.exact({
      id: PropTypes.number.isRequired,
      title: PropTypes.string.isRequired,
      description: PropTypes.string.isRequired,
      order: PropTypes.number.isRequired,
      user: PropTypes.string.isRequired,
      subtasks: PropTypes.arrayOf(
        PropTypes.exact({
          id: PropTypes.number.isRequired,
          title: PropTypes.string.isRequired,
          order: PropTypes.number.isRequired,
          done: PropTypes.bool.isRequired,
        }),
      ),
    }),
  ).isRequired,
  handleActivate: PropTypes.func,
};

Column.defaultProps = {
  handleActivate: () => { },
};

export default Column;
