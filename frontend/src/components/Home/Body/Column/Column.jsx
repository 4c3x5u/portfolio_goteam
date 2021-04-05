import React from 'react';
import PropTypes from 'prop-types';
import { Col } from 'react-bootstrap';

import Task from './Task/Task';
import { capFirstLetterOf } from '../../../../misc/util';

import './column.sass';

const Column = ({ name }) => (
  <Col className="Col" xs={3}>
    <div className={`Column ${capFirstLetterOf(name)}Column`}>
      <div className="Header">{name.toUpperCase()}</div>
      <div className="Body">
        <Task />
      </div>
    </div>
  </Col>
);

Column.propTypes = { name: PropTypes.string.isRequired };

export default Column;
