import React from 'react';
import PropTypes from 'prop-types';

import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faQuestionCircle } from '@fortawesome/free-solid-svg-icons';

import './helptoggler.sass';

const HelpToggler = ({ toggle }) => (
  <Col id="HelpToggler" xs={4}>
    <button
      className="Toggler"
      onClick={toggle}
      aria-label="help modal toggler"
      type="button"
    >
      <FontAwesomeIcon icon={faQuestionCircle} />
      HELP
    </button>
  </Col>
);

HelpToggler.propTypes = {
  toggle: PropTypes.func.isRequired,
};

export default HelpToggler;
