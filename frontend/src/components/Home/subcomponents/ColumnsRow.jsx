import React from 'react';
import PropTypes from 'prop-types';
import { Col, Row } from 'react-bootstrap';

import { capitalizeFirstLetter } from '../../../misc/util';

import './columnsrow.sass';

const Column = ({ name }) => (
  <Col className="Col" xs={3}>
    <div className={`Column ${capitalizeFirstLetter(name)}Column`}>
      <div className="Header">{name.toUpperCase()}</div>
      <div className="Body" />
    </div>
  </Col>
);
Column.propTypes = { name: PropTypes.string.isRequired };

const ColumnsRow = () => (
  <Row id="ColumnsRow">
    <Column name="inbox" />
    <Column name="ready" />
    <Column name="go" />
    <Column name="done" />
  </Row>
);

export default ColumnsRow;
