/* eslint-disable
jsx-a11y/no-static-element-interactions,
jsx-a11y/click-events-have-key-events */

import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Form, Button } from 'react-bootstrap';

import FormGroup from '../../_shared/FormGroup/FormGroup';
import { inputType } from '../../../misc/inputType';

import logo from './createboard.svg';
import './createboard.sass';

const CreateBoard = ({ toggleOff }) => {
  const [title, setTitle] = useState('');

  const handleSubmit = () => console.log('TODO: IMPLEMENT');

  return (
    <div className="CreateBoard" onClick={toggleOff}>
      <Form
        className="Form"
        onSubmit={handleSubmit}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="HeaderWrapper">
          <img className="Header" alt="logo" src={logo} />
        </div>

        <FormGroup
          type={inputType.TEXT}
          label="title"
          value={title}
          setValue={setTitle}
        />

        <div className="ButtonWrapper">
          <Button className="Button" type="submit" aria-label="submit">
            GO!
          </Button>
        </div>
      </Form>
    </div>
  );
};

CreateBoard.propTypes = {
  toggleOff: PropTypes.func.isRequired,
};

export default CreateBoard;
