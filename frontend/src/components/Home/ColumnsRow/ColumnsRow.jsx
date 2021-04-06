import React from 'react';
import PropTypes from 'prop-types';
import { Row } from 'react-bootstrap';

import Column from './Column/Column';
import { columnNames } from './Column/columnNames';

import './columnsrow.sass';

const ColumnsRow = ({ toggleCreateTask }) => (
  <Row id="ColumnsRow">
    <Column name={columnNames.INBOX} toggleCreateTask={toggleCreateTask} />
    <Column name={columnNames.READY} />
    <Column name={columnNames.GO} />
    <Column name={columnNames.DONE} />
  </Row>
);

ColumnsRow.propTypes = { toggleCreateTask: PropTypes.func.isRequired };

export default ColumnsRow;
