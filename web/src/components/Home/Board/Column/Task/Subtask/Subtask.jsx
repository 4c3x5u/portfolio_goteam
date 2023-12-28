import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheckSquare, faSquare } from '@fortawesome/free-regular-svg-icons';

import AppContext from '../../../../../../AppContext';

import './subtask.sass';

const Subtask = ({
  title, done, assignee, toggleDone,
}) => {
  const {
    user
  } = useContext(AppContext);

  return (
    <li className="Subtask">
      <button
        className="CheckButton"
        onClick={toggleDone}
        type="button"
        disabled={!user.isAdmin && assignee && user.username !== assignee}
      >
        {done
          ? <FontAwesomeIcon className="CheckBox" icon={faCheckSquare} />
          : <FontAwesomeIcon className="CheckBox" icon={faSquare} />}

        {title}
      </button>
    </li >
  );
};

Subtask.propTypes = {
  title: PropTypes.string.isRequired,
  done: PropTypes.bool.isRequired,
  assignee: PropTypes.string,
};

export default Subtask;
