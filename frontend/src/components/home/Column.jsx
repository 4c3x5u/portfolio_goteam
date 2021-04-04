import React from 'react';
import PropTypes from 'prop-types';
import { Col } from 'react-bootstrap';

const Column = ({ name, header }) => (
  <Col className="Col" xs={3}>
    <div
      className={
        `Column ${name.charAt(0).toUpperCase() + name.slice(1)}Column`
      }
    >
      <div className="Header">
        <img src={header} alt={`${name} column header`} />
      </div>
      <div className="Body" />
    </div>
  </Col>
);

Column.propTypes = {
  name: PropTypes.string.isRequired,
  header: PropTypes.string.isRequired,
};

export default Column;
