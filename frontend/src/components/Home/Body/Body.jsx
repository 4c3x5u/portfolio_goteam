import React from 'react';
import PropTypes from 'prop-types';
import { Col, Row } from 'react-bootstrap';

import { capFirstLetterOf } from '../../../misc/util';

import './body.sass';

const Column = ({ name }) => (
  <Col className="Col" xs={3}>
    <div className={`Column ${capFirstLetterOf(name)}Column`}>
      <div className="Header">{name.toUpperCase()}</div>
      <div className="Body" />
    </div>
  </Col>
);

Column.propTypes = { name: PropTypes.string.isRequired };

const Body = () => (
  <Row id="Body">
    <Column name="inbox" />
    <Column name="ready" />
    <Column name="go" />
    <Column name="done" />
  </Row>
);

export default Body;
