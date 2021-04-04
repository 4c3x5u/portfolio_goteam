import React from 'react';
import PropTypes from 'prop-types';
import { Col, Button } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/free-solid-svg-icons';

const ControlButton = ({
  name, isToggled, setIsToggled, icon,
}) => (
  <Col xs={4} className="ControlButtonWrapper">
    <Button
      className="Button"
      onClick={setIsToggled}
      aria-label={`${name} button`}
    >
      <FontAwesomeIcon icon={icon} />
      {name.toUpperCase()}
    </Button>
    {isToggled && (
      <div className="TeamControls">
        <ul>
          <li>Some Member</li>
          <li>Some Other Member</li>
          <Button>
            <FontAwesomeIcon icon={faPlusCircle} />
          </Button>
        </ul>
      </div>
    )}

  </Col>
);

ControlButton.propTypes = {
  name: PropTypes.string.isRequired,
  isToggled: PropTypes.bool.isRequired,
  setIsToggled: PropTypes.func.isRequired,
  // eslint-disable-next-line react/forbid-prop-types
  icon: PropTypes.object.isRequired,
};

export default ControlButton;
