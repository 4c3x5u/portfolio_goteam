import React from 'react';
import PropTypes from 'prop-types';
import { Col } from 'react-bootstrap';

const Column = ({ name }) => (
  <Col className="Col" xs={3}>
    <div
      className={
        `Column ${name.charAt(0).toUpperCase() + name.slice(1)}Column`
      }
    >
      <div className="Header">
        {name.toUpperCase()}
      </div>
      <div className="Body" />
    </div>
  </Col>
);

Column.propTypes = {
  name: PropTypes.string.isRequired,
};

export default Column;
