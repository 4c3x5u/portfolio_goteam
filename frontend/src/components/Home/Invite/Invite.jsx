/* eslint-disable
jsx-a11y/no-static-element-interactions,
jsx-a11y/click-events-have-key-events,
jsx-a11y/mouse-events-have-key-events */

import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import {
  Form, Button, OverlayTrigger, Tooltip,
} from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faQuestion } from '@fortawesome/free-solid-svg-icons';

import FormGroup from '../../_shared/FormGroup/FormGroup';
import AppContext from '../../../AppContext';

import logo from './invite.svg';
import './invite.sass';

const Invite = ({ toggleOff }) => {
  const { team } = useContext(AppContext);

  const inviteLink = (
    `${process.env.REACT_APP_FRONTEND_URL}/register/${team.inviteCode}`
  );

  const handleSubmit = (e) => {
    e.preventDefault();

    const el = document.createElement('textarea');
    el.value = inviteLink;
    document.body.appendChild(el);
    el.select();
    document.execCommand('copy');
    document.body.removeChild(el);

    toggleOff();
  };

  return (
    <div className="Invite" onClick={toggleOff}>
      <Form
        className="Form"
        onSubmit={handleSubmit}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="HeaderWrapper">
          <img className="Header" alt="logo" src={logo} />
        </div>

        <FormGroup
          type="text"
          label="INVITE LINK"
          value={inviteLink}
          disabled
        />

        <OverlayTrigger
          placement="bottom"
          overlay={(
            <Tooltip id="HelpTooltip">
              This invite link will automatically be copied to your clipboard
              when you click &quot;GO!&quot;.
              <br />
              <br />
              Send it to your colleagues. Once they register through it, they
              will automatically be added to your team.
            </Tooltip>
          )}
        >
          <Button className="InviteHelp" type="button">
            <FontAwesomeIcon className="Icon" icon={faQuestion} />
          </Button>
        </OverlayTrigger>

        <div className="ButtonWrapper">
          <Button
            className="Button"
            type="submit"
            aria-label="submit"
          >
            GO!
          </Button>
        </div>
      </Form>
    </div>
  );
};

Invite.propTypes = {
  toggleOff: PropTypes.func.isRequired,
};

export default Invite;
