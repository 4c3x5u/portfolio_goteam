/* eslint-disable
jsx-a11y/no-static-element-interactions,
jsx-a11y/click-events-have-key-events */

import React from 'react';
import PropTypes from 'prop-types';
import { Button } from 'react-bootstrap';

import './help.sass';
import logo from './help.svg';

const Help = ({ toggleOff }) => (
  <div className="Help" onClick={toggleOff}>
    <div className="Body" onClick={(e) => e.stopPropagation()}>
      <div className="HeaderWrapper">
        <img className="Header" alt="logo" src={logo} />
      </div>

      <h1>Summary</h1>
      <p>
        In both the control menus and the inbox column, you can click on the
        plus icon to add a new item.
        <br />
        You can also right click on any item to view additional controls for it.
        However, unless you are an admin,
        <br />
        you will be prohibited from accessing all CRUD controls except checking
        subtasks and dragging tasks.
      </p>

      <h1>Team Controls</h1>
      <ol>
        <li>
          Click the TEAM button located on the control bar to toggle the
          team controls menu.
        </li>
        <li>
          The team members that are added to the current board and therefore
          can view all its tasks will
          <br />
          have a yellow arrow to the left of their name, and the team admin
          will have an admin icon.
        </li>
        <li>
          Admins can click on a member to include or exclude them from
          the current board.
        </li>
        <li>
          Admins can view and click on the plus icon at the bottom of the menu
          to invite new members to their team.
        </li>
      </ol>

      <h1>Boards Controls</h1>
      <ol>
        <li>
          Click the BOARDS button located on the control bar to toggle
          the Boards controls menu.
        </li>
        <li>
          The currently board will have a yellow arrow to the left of its name.
        </li>
        <li>Click a board to load it on page.</li>
        <li>
          Admins can click the plus icon at the bottom of the list to create
          a new board.
        </li>
      </ol>

      <h1>Task Controls</h1>
      <ol>
        <li>
          Admins can view and click on the plus icon inside the INBOX column to
          create a new task.
        </li>
        <li>
          Drag and drop tasks from one column to another in order to move
          them around. Members
          <br />
          can only move the the tasks that are assigned to them.
        </li>
        <li>
          Admins can right click on tasks to edit or delete them, or to assign
          them to a team member.
        </li>
      </ol>
      <Button
        className="Button"
        type="button"
        aria-label="cancel"
        onClick={toggleOff}
      >
        OK
      </Button>
    </div>
  </div>
);

Help.propTypes = { toggleOff: PropTypes.func.isRequired };

export default Help;
