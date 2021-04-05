import React from 'react';
import PropTypes from 'prop-types';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/free-solid-svg-icons';

import './controlbutton.sass';

const ControlButton = ({
  name, isToggled, setIsToggled, icon,
}) => (
  <Col xs={4} className="ControlButtonWrapper">
    <button
      className="Button"
      onClick={setIsToggled}
      aria-label={`${name} button`}
      type="button"
    >
      <FontAwesomeIcon icon={icon} />
      {name.toUpperCase()}
    </button>
    {isToggled && (
      <div className="TeamControls">
        <button type="button">
          Some Member
        </button>
        <button type="button">
          Some Other Member
        </button>
        <button type="button">
          <FontAwesomeIcon icon={faPlusCircle} />
        </button>
      </div>
    )}

  </Col>
);

ControlButton.propTypes = {
  name: PropTypes.string.isRequired,
  isToggled: PropTypes.bool.isRequired,
  setIsToggled: PropTypes.func.isRequired,
  icon: PropTypes.string.isRequired,
};

export default ControlButton;
