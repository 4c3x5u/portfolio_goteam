import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheckSquare, faSquare } from '@fortawesome/free-regular-svg-icons';

import './subtask.sass';

const Subtask = ({ title }) => {
  const [checked, setChecked] = useState(false);

  const icon = checked ? (
    <FontAwesomeIcon className="CheckBox" icon={faCheckSquare} />
  ) : (
    <FontAwesomeIcon className="CheckBox" icon={faSquare} />
  );

  return (
    <li className="Subtask">
      <button
        className="CheckButton"
        onClick={() => setChecked(!checked)}
        type="button"
      >
        {icon}
        {title}
      </button>
    </li>
  );
};

Subtask.propTypes = {
  title: PropTypes.string.isRequired,
};

export default Subtask;
