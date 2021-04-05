import React from 'react';
import { Row } from 'react-bootstrap';

import Column from './Column/Column';

import './body.sass';

const Body = () => (
  <Row id="Body">
    <Column name="inbox" />
    <Column name="ready" />
    <Column name="go" />
    <Column name="done" />
  </Row>
);

export default Body;
