import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Form } from 'react-bootstrap';

import FormGroup from '../../_shared/FormGroup';
import { windowEnum } from '../windowEnum';
import logo from '../../../assets/createboard.svg';

import './createboard.sass';

const CreateBoard = ({ deactivate }) => {
  const [title, setTitle] = useState('');

  const handleSubmit = () => console.log('TODO: IMPLEMENT');

  return (
    <button
      className="CreateBoard"
      onClick={deactivate(windowEnum.CREATE_BOARD)}
      type="button"
    >
      <Form
        className="Form"
        onSubmit={handleSubmit}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="HeaderWrapper">
          <img className="Header" alt="logo" src={logo} />
        </div>

        <FormGroup label="title" value={title} setValue={setTitle} />

        <div className="ButtonWrapper">
          <button className="Button" type="submit" aria-label="submit">
            GO!
          </button>
        </div>
      </Form>
    </button>
  );
};

CreateBoard.propTypes = {
  deactivate: PropTypes.func.isRequired,
};

export default CreateBoard;
