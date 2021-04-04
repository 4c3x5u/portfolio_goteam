import React from 'react';
import { Row } from 'react-bootstrap';

import Column from './Column';

const ColumnsRow = () => (
  <Row className="ColumnsRow">
    <Column name="inbox" />
    <Column name="ready" />
    <Column name="go" />
    <Column name="done" />
  </Row>
);

export default ColumnsRow;
