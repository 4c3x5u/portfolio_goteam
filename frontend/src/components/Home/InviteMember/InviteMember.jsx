/* eslint-disable jsx-a11y/mouse-events-have-key-events */
import React, { useState } from 'react';
import PropTypes from 'prop-types';
import {
  Form, Button, OverlayTrigger, Tooltip,
} from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faQuestion } from '@fortawesome/free-solid-svg-icons';

import FormGroup from '../../_shared/FormGroup';

import logo from '../../../assets/invite.svg';
import './invitemember.sass';

const InviteMember = ({ toggleOff }) => {
  const [inviteLink] = useState('TODO: IMPLEMENT');

  const handleSubmit = () => console.log('TODO: IMPLEMENT');

  return (
    <button className="InviteMember" onClick={toggleOff} type="button">
      <Form
        className="Form"
        onSubmit={handleSubmit}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="HeaderWrapper">
          <img className="Header" alt="logo" src={logo} />
        </div>

        <FormGroup
          label="INVITE LINK"
          value={inviteLink}
          setValue={() => console.log('NOT ALLOWED')}
        />

        <OverlayTrigger
          placement="bottom"
          overlay={(
            <Tooltip id="InviteHelp">
              To invite members to your team, copy and send them this invite
              code.
              <br />
              <br />
              Once they register through it, they will be added to your team.
            </Tooltip>
          )}
        >
          <Button className="Help" type="button">
            <FontAwesomeIcon className="Icon" icon={faQuestion} />
          </Button>
        </OverlayTrigger>

        <div className="ButtonWrapper">
          <Button className="Button" type="submit" aria-label="submit">
            GO!
          </Button>
        </div>
      </Form>
    </button>
  );
};

InviteMember.propTypes = {
  toggleOff: PropTypes.func.isRequired,
};

export default InviteMember;
