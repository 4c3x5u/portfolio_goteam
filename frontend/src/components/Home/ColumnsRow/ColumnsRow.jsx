import React from 'react';
import { Row } from 'react-bootstrap';

import Column from './Column/Column';
import { columnNames } from './Column/columnNames';

import './columnsrow.sass';

const ColumnsRow = () => (
  <Row id="ColumnsRow">
    <Column name={columnNames.INBOX} />
    <Column name={columnNames.READY} />
    <Column name={columnNames.GO} />
    <Column name={columnNames.DONE} />
  </Row>
);

export default ColumnsRow;
