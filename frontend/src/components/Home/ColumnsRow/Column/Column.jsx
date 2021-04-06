import React from 'react';
import PropTypes from 'prop-types';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/free-solid-svg-icons';

import Task from './Task/Task';
import { columnNames } from './columnNames';
import { capFirstLetterOf } from '../../../../misc/utils';

import './column.sass';

const Column = ({ name, toggleCreateTask }) => (
  <Col className="Col" xs={3}>
    <div className={`Column ${capFirstLetterOf(name)}Column`}>
      <div className="Header">{name.toUpperCase()}</div>

      <div className="Body">
        <Task />

        {name === columnNames.INBOX && (
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

Column.propTypes = {
  name: PropTypes.string.isRequired,
  toggleCreateTask: PropTypes.func,
};

Column.defaultProps = {
  toggleCreateTask: () => (
    console.log('Cannot create task here.')
  ),
};

export default Column;
