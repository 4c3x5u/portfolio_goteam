import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheckSquare, faSquare } from '@fortawesome/free-regular-svg-icons';

import './subtask.sass';

const Subtask = ({ id, title, done }) => {
  // TODO: Use subtask's 'DONE' field rather than a state here
  const [checked, setChecked] = useState(false);

  // TODO Use subtask ID to handle set done/undone
  console.log(`Subtask ID: ${id}`);

  return (
    <li className="Subtask">
      <button
        className="CheckButton"
        onClick={() => setChecked(!checked)}
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
