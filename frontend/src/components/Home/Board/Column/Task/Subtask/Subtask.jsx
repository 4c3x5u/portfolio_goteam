import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheckSquare, faSquare } from '@fortawesome/free-regular-svg-icons';

import AppContext from '../../../../../../AppContext';
import SubtasksAPI from '../../../../../../api/SubtasksAPI';

import './subtask.sass';

const Subtask = ({
  id, title, done, assignedUser, taskId, columnId,
}) => {
  const {
    user, activeBoard, setActiveBoard, loadBoard, notify,
  } = useContext(AppContext);

  const check = (subtaskId) => {
    // Update client state to avoid load screen.
    setActiveBoard({
      ...activeBoard,
      columns: activeBoard.columns.map((column) => (
        column.id === columnId ? {
          ...column,
          tasks: column.tasks.map((task) => (
            task.id === taskId ? {
              ...task,
              subtasks: task.subtasks.map((subtask) => (
                subtask.id === subtaskId ? {
                  ...subtask,
                  done: !subtask.done,
                } : subtask
              )),
            } : task
          )),
        } : column
      )),
    });

    // Update subtask in database
    SubtasksAPI
      .patch(subtaskId, { done: !done })
      .catch((err) => {
        notify(
          'Unable to check subtask done.',
          `${err?.message || 'Server Error'}.`,
        );
      })
      .finally(loadBoard);
  };

  return (
    <li className="Subtask">
      <button
        className="CheckButton"
        onClick={() => check(id)}
        type="button"
        disabled={!user.isAdmin && user.username !== assignedUser}
      >
        {done
          ? <FontAwesomeIcon className="CheckBox" icon={faCheckSquare} />
          : <FontAwesomeIcon className="CheckBox" icon={faSquare} />}

        {title}
      </button>
    </li>
  );
};

Subtask.propTypes = {
  id: PropTypes.number.isRequired,
  title: PropTypes.string.isRequired,
  done: PropTypes.bool.isRequired,
  assignedUser: PropTypes.string.isRequired,
  columnId: PropTypes.number.isRequired,
  taskId: PropTypes.number.isRequired,
};

export default Subtask;
