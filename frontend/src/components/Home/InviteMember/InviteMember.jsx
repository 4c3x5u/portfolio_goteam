import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Form, Button } from 'react-bootstrap';

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
