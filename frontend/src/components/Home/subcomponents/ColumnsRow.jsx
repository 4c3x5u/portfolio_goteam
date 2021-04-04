import React from 'react';
import { Row } from 'react-bootstrap';

import Column from './Column';

import inboxHeader from '../../../assets/inboxHeader.svg';
import readyHeader from '../../../assets/readyHeader.svg';
import goHeader from '../../../assets/goHeader.svg';
import doneHeader from '../../../assets/doneHeader.svg';

const ColumnsRow = () => (
  <Row className="ColumnsRow">
    <Column name="inbox" header={inboxHeader} />
    <Column name="ready" header={readyHeader} />
    <Column name="go" header={goHeader} />
    <Column name="done" header={doneHeader} />
  </Row>
);

export default ColumnsRow;
