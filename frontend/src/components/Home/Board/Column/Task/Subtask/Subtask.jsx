import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheckSquare, faSquare } from '@fortawesome/free-regular-svg-icons';

import AppContext from '../../../../../../AppContext';
import { patchSubtask } from '../../../../../../misc/api';

import './subtask.sass';

const Subtask = ({ id, title, done }) => {
  const { loadBoard } = useContext(AppContext);

  const check = (subtaskId) => (
    patchSubtask(subtaskId, {
      done: !done,
    }).then(() => (
      loadBoard()
    )).catch((err) => (
      console.error(err)
    ))
  );

  return (
    <li className="Subtask">
      <button
        className="CheckButton"
        onClick={() => check(id)}
        type="button"
      >
        {
          done
            ? <FontAwesomeIcon className="CheckBox" icon={faCheckSquare} />
            : <FontAwesomeIcon className="CheckBox" icon={faSquare} />
        }

        {title}
      </button>
    </li>
  );
};

Subtask.propTypes = {
  id: PropTypes.number.isRequired,
  title: PropTypes.string.isRequired,
  done: PropTypes.bool.isRequired,
};

export default Subtask;
